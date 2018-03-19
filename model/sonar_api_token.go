package model

import "encoding/json"

// ApiTokens is a security token from SonarQube to be used for calling the sonarqube api.
type ApiTokens struct {
	Login      string `json:"login"`
	UserTokens []struct {
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
	} `json:"userTokens"`
}

// Unmarshal is a implementation of ApiTokens interface's Unmarshal
func (s *ApiTokens) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &s)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return s, nil
}

// ApiToken is the single token structure you get back from SonarQube when creating a new API Token.
type ApiToken struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// Unmarshal is a implementation of ApiTokens interface's Unmarshal
func (s *ApiToken) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &s)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return s, nil
}
