package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/usecases"
)

// CLIScannerHandler handles CLI commands for the port scanner.
type CLIScannerHandler struct {
	scanUseCase     *usecases.ScanUseCase
	discoverUseCase *usecases.DiscoverUseCase
	reportUseCase   *usecases.ReportUseCase
	printer         *Printer
}

// NewCLIScannerHandler creates a new CLI handler.
func NewCLIScannerHandler(
	scanUC *usecases.ScanUseCase,
	discoverUC *usecases.DiscoverUseCase,
	reportUC *usecases.ReportUseCase,
	printer *Printer,
) *CLIScannerHandler {
	return &CLIScannerHandler{
		scanUseCase:     scanUC,
		discoverUseCase: discoverUC,
		reportUseCase:   reportUC,
		printer:         printer,
	}
}

// RunScan executes a port scan with the given parameters.
func (h *CLIScannerHandler) RunScan(targets []string, ports []int, scanType models.ScanType) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		h.printer.Warning("\nScan interrupted by user. Exiting...")
		cancel()
	}()

	h.printer.Info("Starting scan (type: %s)", scanType)

	output, err := h.scanUseCase.Execute(ctx, usecases.PortScannerInput{
		Targets:  targets,
		Ports:    ports,
		ScanType: scanType,
	})
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	h.printer.PrintResult(output.Result)
	return nil
}

// DiscoverHosts discovers hosts on a network CIDR.
func (h *CLIScannerHandler) DiscoverHosts(cidr string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	h.printer.Info("Discovering hosts in %s...", cidr)

	output, err := h.discoverUseCase.Execute(ctx, usecases.HostDiscoveryInput{
		CIDR: cidr,
	})
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	h.printer.PrintHosts(output.Hosts)
	return nil
}
