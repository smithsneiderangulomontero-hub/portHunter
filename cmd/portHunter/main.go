package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/smithsneiderangulomontero-hub/portHunter/config"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/infrastructure/network"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/infrastructure/output"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/interfaces/cli"
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/usecases"
	"github.com/smithsneiderangulomontero-hub/portHunter/pkg/utils"
)

func main() {
	cfg := config.Load()

	target := flag.String("t", cfg.DefaultTarget, "Target IP, hostname, or CIDR (e.g. 192.168.1.0/24)")
	ports := flag.String("p", cfg.DefaultPorts, "Ports to scan (e.g. 80,443 or 1-1000)")
	scanType := flag.String("s", string(cfg.ScanType), "Scan type: tcp_connect, syn, udp, ping")
	workers := flag.Int("w", cfg.Workers, "Number of concurrent workers")
	timeout := flag.Int("timeout", cfg.Timeout, "Timeout per port in milliseconds")
	_ = flag.String("o", cfg.OutputFormat, "Output format: table, json")
	discover := flag.Bool("discover", false, "Discover live hosts in the target network")
	flag.Parse()

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// --- Dependency Injection (wiring de Clean Architecture) ---
	scannerSvc := network.NewTCPScannerService()
	discoverySvc := network.NewPingDiscoveryService()
	jsonReporter := output.NewJSONReportService(true)

	scanUC := usecases.NewScanUseCase(scannerSvc)
	discoverUC := usecases.NewDiscoverUseCase(discoverySvc)
	reportUC := usecases.NewReportUseCase(jsonReporter)

	printer := cli.NewPrinter()
	handler := cli.NewCLIScannerHandler(scanUC, discoverUC, reportUC, printer)

	if *discover {
		if err := handler.DiscoverHosts(*target); err != nil {
			printer.Error("Discovery failed: %v", err)
			os.Exit(1)
		}
		return
	}

	portList, err := utils.ParsePorts(*ports)
	if err != nil {
		printer.Error("Invalid ports: %v", err)
		os.Exit(1)
	}

	targets, err := utils.ParseCIDR(*target)
	if err != nil {
		printer.Error("Invalid target: %v", err)
		os.Exit(1)
	}

	printer.Info("portHunter v1.0.0 - Ethical Port Scanner")
	printer.Info("Target: %s (%d hosts)", *target, len(targets))
	printer.Info("Ports: %d", len(portList))
	printer.Info("Workers: %d", *workers)
	printer.Info("Timeout: %dms", *timeout)

	if err := handler.RunScan(targets, portList, models.ScanType(*scanType)); err != nil {
		printer.Error("Scan failed: %v", err)
		os.Exit(1)
	}
}
