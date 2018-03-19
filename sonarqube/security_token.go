package sonarqube

import (
	"../model"
	"../util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func RetrieveAPIToken(sonarQubeConfig model.SonarQubeConfig) model.ApiTokens {
	fmt.Printf(">> Retrieving SonarQube API Tokens\n")
	tokens := model.ApiTokens{}
	rawUrl := util.ReplaceHostnameAndPort("http://XXX:YYY/ZZZ/api/user_tokens/search", sonarQubeConfig.InternalHostname, sonarQubeConfig.InternalPort)
	rawUrl = strings.Replace(rawUrl, "/ZZZ", sonarQubeConfig.ContextRoot, 1)
	var apiUrl *url.URL
	apiUrl, err := url.Parse(rawUrl)
	if err != nil {
		panic("boom")
	}

	urlStr := apiUrl.String()
	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, strings.NewReader("")) // <-- URL-encoded payload
	r.SetBasicAuth(sonarQubeConfig.SonarQubeUser, sonarQubeConfig.SonarQubePass)

	fmt.Printf(" > URL: %v\n", urlStr)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf(" < Failed to update setting: %v", err)
	} else {
		fmt.Printf(" < %v\n", resp.Status)
		if resp.StatusCode != http.StatusOK {
			fmt.Printf(" > call failed: %v", resp.StatusCode)
			return tokens
		}
		defer resp.Body.Close()
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		} else {
			json.Unmarshal([]byte(body), &tokens)
			fmt.Printf("  > Found tokens: %v\n", tokens)
		}
	}
	return tokens
}

func GenerateAPIToken(sonarQubeConfig model.SonarQubeConfig) model.ApiToken {
	token := model.ApiToken{}

	fmt.Printf(">> Generating SonarQube API Token for:%v\n", sonarQubeConfig.APITokenName)
	rawUrl := util.ReplaceHostnameAndPort("http://XXX:YYY/ZZZ/api/user_tokens/generate", sonarQubeConfig.InternalHostname, sonarQubeConfig.InternalPort)
	rawUrl = strings.Replace(rawUrl, "/ZZZ", sonarQubeConfig.ContextRoot, 1)
	var apiUrl *url.URL
	apiUrl, err := url.Parse(rawUrl)
	if err != nil {
		panic("boom")
	}

	parameters := url.Values{}
	parameters.Add("name", sonarQubeConfig.APITokenName)
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
		if resp.StatusCode != http.StatusOK {
			fmt.Printf(" > call failed: %v", resp.StatusCode)
			return token
		}
		defer resp.Body.Close()
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		} else {
			json.Unmarshal([]byte(body), &token)
			fmt.Printf("  > Found token: %v\n", token)
		}
	}
	return token
}
