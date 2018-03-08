package sonarqube

import (
	"../model"
	"../util"
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"os"
	"io/ioutil"
)

func RetrieveAPIToken(sonarQubeConfig model.SonarQubeConfig) {
	fmt.Printf(">> Retrieving SonarQube API Tokens\n")
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
		defer resp.Body.Close()
		bodyString := "N/A"
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			check(err2)
			bodyString = string(bodyBytes)
		}
		fmt.Printf(" < %v\n", bodyString)
		writeTokenToFile("./sonar-api-tokens.txt", bodyString)
	}
}

func GenerateAPIToken(sonarQubeConfig model.SonarQubeConfig) {
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
		defer resp.Body.Close()
		bodyString := "N/A"
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			check(err2)
			bodyString = string(bodyBytes)
		}
		fmt.Printf(" < %v\n", bodyString)
		writeTokenToFile("./sonar-api-token.txt", bodyString)
	}
}
func writeTokenToFile(filename string, body string) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()
	_, err2 := f.WriteString(body)
	check(err2)
	//	Issue a Sync to flush writes to stable storage.
	f.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}