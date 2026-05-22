package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
)

// Printer handles formatted output to the terminal.
type Printer struct {
	writer *tabwriter.Writer
}

// NewPrinter creates a new Printer.
func NewPrinter() *Printer {
	return &Printer{
		writer: tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0),
	}
}

// Info prints an informational message.
func (p *Printer) Info(format string, args ...interface{}) {
	fmt.Printf("[*] "+format+"\n", args...)
}

// Warning prints a warning message.
func (p *Printer) Warning(format string, args ...interface{}) {
	fmt.Printf("[!] "+format+"\n", args...)
}

// Error prints an error message.
func (p *Printer) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[-] "+format+"\n", args...)
}

// Success prints a success message.
func (p *Printer) Success(format string, args ...interface{}) {
	fmt.Printf("[+] "+format+"\n", args...)
}

// PrintResult prints a scan result in a readable format.
func (p *Printer) PrintResult(result models.ScanResult) {
	fmt.Println()
	p.Info("Scan Results")
	p.Info("============")
	p.Info("Target: %s", result.Target)
	p.Info("Scan Type: %s", result.ScanType)
	p.Info("Duration: %v", result.Duration)
	p.Info("Total Hosts: %d", result.TotalHosts)
	p.Info("Open Ports: %d", len(result.OpenPorts))
	fmt.Println()

	if len(result.OpenPorts) == 0 {
		p.Warning("No open ports found.")
		return
	}

	fmt.Fprintln(p.writer, "PORT\tSTATE\tSERVICE\tVERSION")
	fmt.Fprintln(p.writer, "----\t-----\t-------\t-------")
	for _, port := range result.OpenPorts {
		fmt.Fprintf(p.writer, "%d/%s\t%s\t%s\t%s\n",
			port.Number, port.Protocol, port.State, port.Service, port.Version)
	}
	p.writer.Flush()
}

// PrintHosts prints discovered hosts in a table.
func (p *Printer) PrintHosts(hosts []models.Host) {
	fmt.Println()
	if len(hosts) == 0 {
		p.Warning("No hosts found.")
		return
	}

	p.Success("Found %d host(s):", len(hosts))
	fmt.Fprintln(p.writer, "IP\tMAC\tHOSTNAME\tSTATUS")
	fmt.Fprintln(p.writer, "--\t---\t--------\t------")
	for _, host := range hosts {
		fmt.Fprintf(p.writer, "%s\t%s\t%s\t%s\n",
			host.IP.String(), host.MAC, host.Hostname, host.Status)
	}
	p.writer.Flush()
}

// PrintProgress prints a progress indicator.
func (p *Printer) PrintProgress(current, total int, startTime time.Time) {
	elapsed := time.Since(startTime)
	percent := float64(current) / float64(total) * 100
	fmt.Printf("\r[*] Progress: %d/%d (%.1f%%) | Elapsed: %v",
		current, total, percent, elapsed.Round(time.Second))
}
