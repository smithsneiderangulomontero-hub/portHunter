package usecases

import (
	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
)

// ReportService defines how scan results are reported/output.
type ReportService interface {
	Generate(result models.ScanResult) (string, error)
}

// ReportUseCase handles generating reports from scan results.
type ReportUseCase struct {
	reporter ReportService
}

// NewReportUseCase creates a new ReportUseCase.
func NewReportUseCase(reporter ReportService) *ReportUseCase {
	return &ReportUseCase{reporter: reporter}
}

// GenerateReport generates a report from a scan result.
func (uc *ReportUseCase) GenerateReport(result models.ScanResult) (string, error) {
	return uc.reporter.Generate(result)
}
