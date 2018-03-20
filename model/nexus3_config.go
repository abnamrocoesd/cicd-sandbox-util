package model

// Nexus3Config basic configuration details required for interacting with a Nexus 3 server.
type Nexus3Config struct {
	ExternalHostname string
	ExternalPort     string
	InternalHostname string
	InternalPort     string
	ContextRoot      string
	User             string
	Pass             string
}
