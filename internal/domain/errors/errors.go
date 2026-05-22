package errors

import "fmt"

var (
	ErrInvalidIP        = fmt.Errorf("invalid IP address")
	ErrInvalidPort      = fmt.Errorf("invalid port number")
	ErrPermissionDenied = fmt.Errorf("permission denied: raw sockets require root privileges")
	ErrScanTimeout      = fmt.Errorf("scan timed out")
	ErrNoHostsFound     = fmt.Errorf("no hosts found in the target range")
	ErrContextCanceled  = fmt.Errorf("scan was canceled")
)

// ScanError wraps errors that occur during scanning.
type ScanError struct {
	Op   string
	Err  error
	Host string
	Port int
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("%s failed for %s:%d: %v", e.Op, e.Host, e.Port, e.Err)
}

func (e *ScanError) Unwrap() error {
	return e.Err
}
