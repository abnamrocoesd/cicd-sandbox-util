package jenkins

// TODO: sonar <> Jenkins
/* 		Params:	jenkinsUser/Pass
		jenkinsHost
		jenkinsPort
		jenkinsJobUrl
		jenkinsJobParam
		sonarQubeUser/Pass
		sonarQubeTokenName
		sonarQubeHost
		sonarQubePort
action:
	- retrieve tokens from sonar
	- if requested token (sonarQubeTokenName) does not exist, create it
	- trigger jenkinsJob with token param
*/
import (
	"../model"
	"../util"
	"fmt"
	"strings"
	"net/url"
	"net/http"
)

func SonarConfiguration(sonarQubeConfig model.SonarQubeConfig, jobConfig model.JenkinsJobConfig, token model.ApiToken) {
	fmt.Printf(">> Updating Jenkins' SonarQube Configuration\n")
	// http://localhost:8282/jenkins/job/configs/job/sonar/buildWithParameters
	// auth: jenkins User, jenkins User TOKEN
	// form data: sonar_token
	// success = 201 Created
	rawUrl := util.ReplaceHostnameAndPort("http://XXX:YYY/ZZZ/buildByToken/buildWithParameters",
		jobConfig.Host,
		jobConfig.Port)
	rawUrl = strings.Replace(rawUrl, "/ZZZ", jobConfig.ContextRoot, 1)
	//rawUrl = strings.Replace(rawUrl, "/AAA", jobConfig.JobUrl, 1)
	fmt.Printf(" > Using job url: %v\n", rawUrl)
	var apiUrl *url.URL
	apiUrl, err := url.Parse(rawUrl)
	if err != nil {
		panic("boom")
	}

	parameters := url.Values{}
	parameters.Add("job", jobConfig.JobUrl)
	parameters.Add("sonar_token", token.Token)
	parameters.Add("cause", "Auto configuration")
	parameters.Add("token", jobConfig.Token)
	apiUrl.RawQuery = parameters.Encode()
	urlStr := apiUrl.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader("")) // <-- URL-encoded payload
	//r.SetBasicAuth(jobConfig.User, jobConfig.Token)
	// if using https://plugins.jenkins.io/build-token-root we can skip basic auth

	fmt.Printf(" > URL: %v\n", urlStr)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf(" < Failed to update setting: %v", err)
	} else {
		if resp.StatusCode == http.StatusCreated {
			fmt.Printf(" > Sonar Config job successfully scheduled\n")
		} else {
			fmt.Printf(" > call failed: %v\n", resp.StatusCode)
		}
	}
}
