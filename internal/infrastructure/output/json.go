package output

import (
	"encoding/json"

	"github.com/smithsneiderangulomontero-hub/portHunter/internal/domain/models"
)

// JSONReportService generates JSON reports from scan results.
type JSONReportService struct {
	Pretty bool
}

// NewJSONReportService creates a new JSON report service.
func NewJSONReportService(pretty bool) *JSONReportService {
	return &JSONReportService{Pretty: pretty}
}

// Generate produces a JSON string representation of the scan result.
func (rs *JSONReportService) Generate(result models.ScanResult) (string, error) {
	var data []byte
	var err error

	if rs.Pretty {
		data, err = json.MarshalIndent(result, "", "  ")
	} else {
		data, err = json.Marshal(result)
	}
	if err != nil {
		return "", err
	}

	return string(data), nil
}
