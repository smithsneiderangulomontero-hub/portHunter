package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/usecases"
)

// TCPScannerService implements PortScannerService using TCP Connect scanning.
type TCPScannerService struct{}

// NewTCPScannerService creates a new TCP scanner service.
func NewTCPScannerService() *TCPScannerService {
	return &TCPScannerService{}
}

// Scan performs a TCP Connect scan against the specified targets and ports.
func (s *TCPScannerService) Scan(ctx context.Context, input usecases.PortScannerInput) (*usecases.PortScannerOutput, error) {
	startTime := time.Now()
	result := models.ScanResult{
		ScanType:  input.ScanType,
		StartTime: startTime,
	}

	var mu sync.Mutex
	sem := semaphore.NewWeighted(int64(input.Workers))
	var wg sync.WaitGroup

	for _, target := range input.Targets {
		for _, port := range input.Ports {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			if err := sem.Acquire(ctx, 1); err != nil {
				return nil, err
			}
			wg.Add(1)

			go func(target string, port int) {
				defer sem.Release(1)
				defer wg.Done()

				host, open := s.scanPort(ctx, target, port, input.Timeout)
				if open != nil {
					mu.Lock()
					result.OpenPorts = append(result.OpenPorts, *open)
					result.Hosts = append(result.Hosts, *host)
					mu.Unlock()
				}
			}(target, port)
		}
	}

	wg.Wait()
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.TotalHosts = len(result.Hosts)
	result.TotalPorts = len(result.OpenPorts)

	return &usecases.PortScannerOutput{Result: result}, nil
}

func (s *TCPScannerService) scanPort(ctx context.Context, target string, port int, timeout time.Duration) (*models.Host, *models.Port) {
	addr := net.JoinHostPort(target, fmt.Sprintf("%d", port))
	var dialer net.Dialer
	dialer.Timeout = timeout

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, nil
	}
	conn.Close()

	names, _ := net.LookupAddr(target)
	hostname := ""
	if len(names) > 0 {
		hostname = names[0]
	}

	service := knownServices[port]

	host := &models.Host{
		IP:       net.ParseIP(target),
		Hostname: hostname,
		Status:   models.HostStatusAlive,
	}

	openPort := &models.Port{
		Number:   port,
		Protocol: "tcp",
		State:    models.PortStateOpen,
		Service:  service,
	}

	return host, openPort
}

var knownServices = map[int]string{
	21: "FTP", 22: "SSH", 23: "Telnet", 25: "SMTP", 53: "DNS",
	80: "HTTP", 110: "POP3", 111: "RPC", 135: "RPC", 139: "NetBIOS",
	143: "IMAP", 443: "HTTPS", 445: "SMB", 993: "IMAPS", 995: "POP3S",
	1433: "MSSQL", 1521: "Oracle", 2049: "NFS", 3306: "MySQL",
	3389: "RDP", 5432: "PostgreSQL", 5900: "VNC", 5985: "WinRM-HTTP",
	5986: "WinRM-HTTPS", 6379: "Redis", 8080: "HTTP-Proxy",
	8443: "HTTPS-Alt", 9000: "SonarQube", 9090: "WebLogic",
	27017: "MongoDB",
}
