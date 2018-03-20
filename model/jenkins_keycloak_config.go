package model

// JenkinsKeycloakConfig json structure for the Keycloak configuration for the Jenkins Keycloak plugin.
type JenkinsKeycloakConfig struct {
	Realm         string `json:"realm"`
	AuthServerURL string `json:"auth-server-url"`
	SslRequired   string `json:"ssl-required"`
	Resource      string `json:"resource"`
	PublicClient  bool   `json:"public-client"`
}
