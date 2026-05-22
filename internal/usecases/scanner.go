package usecases

import (
	"context"
	"time"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
)

// PortScannerInput defines the input for a port scanning operation.
type PortScannerInput struct {
	Targets  []string
	Ports    []int
	ScanType models.ScanType
	Timeout  time.Duration
	Workers  int
}

// PortScannerOutput defines the output from a port scanning operation.
type PortScannerOutput struct {
	Result models.ScanResult
}

// PortScannerService is the interface that the infrastructure must implement.
type PortScannerService interface {
	Scan(ctx context.Context, input PortScannerInput) (*PortScannerOutput, error)
}

// ScanUseCase orchestrates the port scanning business logic.
type ScanUseCase struct {
	scanner PortScannerService
}

// NewScanUseCase creates a new ScanUseCase.
func NewScanUseCase(scanner PortScannerService) *ScanUseCase {
	return &ScanUseCase{scanner: scanner}
}

// Execute runs the port scan use case.
func (uc *ScanUseCase) Execute(ctx context.Context, input PortScannerInput) (*PortScannerOutput, error) {
	if input.Timeout == 0 {
		input.Timeout = 3 * time.Second
	}
	if input.Workers <= 0 {
		input.Workers = 100
	}
	return uc.scanner.Scan(ctx, input)
}

// DefaultPorts returns a list of commonly scanned ports.
func DefaultPorts() []int {
	return []int{
		21, 22, 23, 25, 53, 80, 110, 111, 135, 139,
		143, 443, 445, 993, 995, 1433, 1521, 2049,
		3306, 3389, 5432, 5900, 5985, 5986, 6379,
		8080, 8443, 9000, 9090, 27017,
	}
}
