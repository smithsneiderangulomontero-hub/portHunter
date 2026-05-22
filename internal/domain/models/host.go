package models

import "net"

type Host struct {
	IP       net.IP
	MAC      string
	Hostname string
	OS       string
	Status   HostStatus
	Ports    []Port
}

type HostStatus string

const (
	HostStatusAlive   HostStatus = "alive"
	HostStatusDead    HostStatus = "dead"
	HostStatusUnknown HostStatus = "unknown"
)
