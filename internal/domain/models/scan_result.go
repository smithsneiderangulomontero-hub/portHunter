package models

import "time"

// ScanResult holds the complete result of a scanning session.
type ScanResult struct {
	Target     string
	ScanType   ScanType
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Hosts      []Host
	OpenPorts  []Port
	TotalPorts int
	TotalHosts int
}

// ScanType defines the type of scan to perform.
type ScanType string

const (
	ScanTypeTCPConnect ScanType = "tcp_connect"
	ScanTypeSYN        ScanType = "syn"
	ScanTypeUDP        ScanType = "udp"
	ScanTypePing       ScanType = "ping"
)
