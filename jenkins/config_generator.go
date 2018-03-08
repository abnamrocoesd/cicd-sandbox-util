package jenkins

import (
	"../model"
	"../util"
	"encoding/json"
	"fmt"
	"log"
)

func KeycloakConfig(hostname string, port string, securityRealm string) {
	realm := securityRealm
	sslRequired := "external"
	resource := "jenkins"
	publicClient := true
	serverUrlTemplate := "http://XXX:YYY/auth"
	authServerUrl := util.ReplaceHostnameAndPort(serverUrlTemplate, hostname, port)
	config := model.JenkinsKeycloakConfig{
		Realm:         realm,
		SslRequired:   sslRequired,
		PublicClient:  publicClient,
		Resource:      resource,
		AuthServerURL: authServerUrl,
	}
	// Convert structs to JSON.
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
