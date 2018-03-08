package sonarqube

import (
	"fmt"
	"strings"
	"net/url"
	"../model"
	"../util"
	"encoding/json"
	"log"
	"net/http"
)

func UpdateSonarQubeConfig(openidcConfig string, sonarQubeConfig model.SonarQubeConfig) {
	serverBaseUrl := fmt.Sprintf("http://%v:%v%v", sonarQubeConfig.ExternalHostname, sonarQubeConfig.ExternalPort, sonarQubeConfig.ContextRoot)
	UpdateSonarQubeSettings("sonar.core.serverBaseURL", serverBaseUrl, sonarQubeConfig)
	UpdateSonarQubeSettings("sonar.auth.oidc.clientId.secured", sonarQubeConfig.KeycloakClientId, sonarQubeConfig)
	UpdateSonarQubeSettings("sonar.auth.oidc.enabled", "true", sonarQubeConfig)
	UpdateSonarQubeSettings("sonar.auth.oidc.groupsSync", "true", sonarQubeConfig)
	UpdateSonarQubeSettings("sonar.auth.oidc.groupsSync.claimName", "groups", sonarQubeConfig)
	UpdateSonarQubeSettings("sonar.auth.oidc.providerConfiguration", openidcConfig, sonarQubeConfig)
}

func UpdateSonarQubeSettings(key string, value string, sonarQubeConfig model.SonarQubeConfig) {
	fmt.Printf(">> Updating SonarQube settings: key=%v, value=%v\n", key, value)
	rawUrl := util.ReplaceHostnameAndPort("http://XXX:YYY/ZZZ/api/settings/set", sonarQubeConfig.InternalHostname, sonarQubeConfig.InternalPort)
	rawUrl = strings.Replace(rawUrl, "/ZZZ", sonarQubeConfig.ContextRoot, 1)
	var apiUrl *url.URL
	apiUrl, err := url.Parse(rawUrl)
	if err != nil {
		panic("boom")
	}

	parameters := url.Values{}
	parameters.Add("key", key)
	parameters.Add("value", value)
	apiUrl.RawQuery = parameters.Encode()
	urlStr := apiUrl.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader("")) // <-- URL-encoded payload
	r.SetBasicAuth(sonarQubeConfig.SonarQubeUser, sonarQubeConfig.SonarQubePass)

	fmt.Printf(" > URL: %v\n", urlStr)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf(" < Failed to update setting: %v", err)
	} else {
		fmt.Printf(" < %v\n", resp.Status)
	}

}

func GenerateSonarKeycloakConfig(externalKeycloakHostname string, keycloakPort string, sonarQubeConfig model.SonarQubeConfig) string {
	// TODO: replace realm with flag input
	issuer := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	authorizationEndpoint := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/auth", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	tokenEndpoint := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/token", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	tokenIntrospectionEndpoint := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/token/introspect", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	userinfoEndpoint := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/userinfo", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	endSessionEndpoint := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/logout", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	jwksURI := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/certs", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	checkSessionIframe := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/protocol/openid-connect/login-status-iframe.html", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	grantTypesSupported := []string{"authorization_code", "implicit", "refresh_token", "password", "client_credentials"}
	responseTypesSupported := []string{"code", "none", "id_token", "token", "id_token token", "code id_token", "code token", "code id_token token"}
	subjectTypesSupported := []string{"public", "pairwise"}
	iDTokenSigningAlgValuesSupported := []string{"RS256"}
	userinfoSigningAlgValuesSupported := []string{"RS256"}
	requestObjectSigningAlgValuesSupported := []string{"none", "RS256"}
	responseModesSupported := []string{"query", "fragment", "form_post"}
	registrationEndpoint := util.ReplaceHostnamePortAndRealm("http://XXX:YYY/auth/realms/ZZZ/clients-registrations/openid-connect", externalKeycloakHostname, keycloakPort, sonarQubeConfig.SecurityRealm)
	tokenEndpointAuthMethodsSupported := []string{"private_key_jwt", "client_secret_basic", "client_secret_post"}
	tokenEndpointAuthSigningAlgValuesSupported := []string{"RS256"}
	claimsSupported := []string{"sub", "iss", "auth_time", "name", "given_name", "family_name", "preferred_username", "email"}
	claimTypesSupported := []string{"normal"}
	claimsParameterSupported := false
	scopesSupported := []string{"openid", "offline_access"}
	requestParameterSupported := true
	requestURIParameterSupported := true

	config := model.SonarKeycloakConfig{
		Issuer:                                     issuer,
		AuthorizationEndpoint:                      authorizationEndpoint,
		TokenEndpoint:                              tokenEndpoint,
		TokenIntrospectionEndpoint:                 tokenIntrospectionEndpoint,
		UserinfoEndpoint:                           userinfoEndpoint,
		EndSessionEndpoint:                         endSessionEndpoint,
		JwksURI:                                    jwksURI,
		CheckSessionIframe:                         checkSessionIframe,
		GrantTypesSupported:                        grantTypesSupported,
		ResponseTypesSupported:                     responseTypesSupported,
		SubjectTypesSupported:                      subjectTypesSupported,
		IDTokenSigningAlgValuesSupported:           iDTokenSigningAlgValuesSupported,
		UserinfoSigningAlgValuesSupported:          userinfoSigningAlgValuesSupported,
		RequestObjectSigningAlgValuesSupported:     requestObjectSigningAlgValuesSupported,
		ResponseModesSupported:                     responseModesSupported,
		RegistrationEndpoint:                       registrationEndpoint,
		TokenEndpointAuthMethodsSupported:          tokenEndpointAuthMethodsSupported,
		TokenEndpointAuthSigningAlgValuesSupported: tokenEndpointAuthSigningAlgValuesSupported,
		ClaimsSupported:                            claimsSupported,
		ClaimTypesSupported:                        claimTypesSupported,
		ClaimsParameterSupported:                   claimsParameterSupported,
		ScopesSupported:                            scopesSupported,
		RequestParameterSupported:                  requestParameterSupported,
		RequestURIParameterSupported:               requestURIParameterSupported,
	}

	// Convert structs to JSON.
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
	openidcConfig := fmt.Sprintf("%s", data)
	return openidcConfig
}