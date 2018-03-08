package model

type JenkinsKeycloakConfig struct {
	Realm         string `json:"realm"`
	AuthServerURL string `json:"auth-server-url"`
	SslRequired   string `json:"ssl-required"`
	Resource      string `json:"resource"`
	PublicClient  bool   `json:"public-client"`
}
