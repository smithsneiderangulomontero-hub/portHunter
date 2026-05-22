package network

import (
	"context"
	"net"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/usecases"
)

// PingDiscoveryService implements HostDiscoveryService using TCP ping.
type PingDiscoveryService struct{}

// NewPingDiscoveryService creates a new PingDiscoveryService.
func NewPingDiscoveryService() *PingDiscoveryService {
	return &PingDiscoveryService{}
}

// Discover discovers live hosts in a CIDR range.
func (s *PingDiscoveryService) Discover(ctx context.Context, input usecases.HostDiscoveryInput) (*usecases.HostDiscoveryOutput, error) {
	ip, ipnet, err := net.ParseCIDR(input.CIDR)
	if err != nil {
		return nil, err
	}

	output := &usecases.HostDiscoveryOutput{}
	var mu sync.Mutex
	sem := semaphore.NewWeighted(int64(input.Workers))
	var wg sync.WaitGroup

	ip = ip.Mask(ipnet.Mask)
	for ; ipnet.Contains(ip); incrementIP(ip) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		ipCopy := net.ParseIP(ip.String())
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		wg.Add(1)

		go func(ip net.IP) {
			defer sem.Release(1)
			defer wg.Done()

			host := s.pingHost(ctx, ip, time.Duration(input.Timeout)*time.Millisecond)
			if host != nil {
				mu.Lock()
				output.Hosts = append(output.Hosts, *host)
				mu.Unlock()
			}
		}(ipCopy)
	}

	wg.Wait()
	return output, nil
}

func (s *PingDiscoveryService) pingHost(ctx context.Context, ip net.IP, timeout time.Duration) *models.Host {
	// TCP ping to port 80
	addr := net.JoinHostPort(ip.String(), "80")
	var dialer net.Dialer
	dialer.Timeout = timeout

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		// Try port 443
		addr = net.JoinHostPort(ip.String(), "443")
		conn, err = dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil
		}
	}
	conn.Close()

	names, _ := net.LookupAddr(ip.String())
	hostname := ""
	if len(names) > 0 {
		hostname = names[0]
	}

	return &models.Host{
		IP:       ip,
		Hostname: hostname,
		Status:   models.HostStatusAlive,
	}
}

func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
