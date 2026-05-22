package models

// Port represents a network port and its state.
type Port struct {
	Number   int
	Protocol string // "tcp" or "udp"
	State    PortState
	Service  string
	Version  string
}

// PortState represents the state of a port after scanning.
type PortState string

const (
	PortStateOpen     PortState = "open"
	PortStateClosed   PortState = "closed"
	PortStateFiltered PortState = "filtered"
)
