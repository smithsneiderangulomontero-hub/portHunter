package usecases

import (
	"context"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
)

// HostDiscoveryInput defines the input for host discovery.
type HostDiscoveryInput struct {
	CIDR    string // e.g. "192.168.1.0/24"
	Timeout int    // milliseconds
	Workers int
}

// HostDiscoveryOutput defines the output from host discovery.
type HostDiscoveryOutput struct {
	Hosts []models.Host
}

// HostDiscoveryService is the interface for discovering hosts on a network.
type HostDiscoveryService interface {
	Discover(ctx context.Context, input HostDiscoveryInput) (*HostDiscoveryOutput, error)
}

// DiscoverUseCase orchestrates the host discovery business logic.
type DiscoverUseCase struct {
	discoverer HostDiscoveryService
}

// NewDiscoverUseCase creates a new DiscoverUseCase.
func NewDiscoverUseCase(discoverer HostDiscoveryService) *DiscoverUseCase {
	return &DiscoverUseCase{discoverer: discoverer}
}

// Execute runs the host discovery use case.
func (uc *DiscoverUseCase) Execute(ctx context.Context, input HostDiscoveryInput) (*HostDiscoveryOutput, error) {
	if input.Timeout <= 0 {
		input.Timeout = 5000
	}
	if input.Workers <= 0 {
		input.Workers = 50
	}
	return uc.discoverer.Discover(ctx, input)
}
