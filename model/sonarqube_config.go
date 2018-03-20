package model

// SonarQubeConfig contains all the information we might need for configuring settings in SonarQube (6.3+) instance.
type SonarQubeConfig struct {
	ExternalHostname string
	ExternalPort     string
	SecurityRealm    string
	InternalHostname string
	InternalPort     string
	ContextRoot      string
	KeycloakClientId string
	SonarQubeUser    string
	SonarQubePass    string
	APITokenName     string
}
