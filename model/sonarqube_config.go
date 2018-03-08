package model

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
