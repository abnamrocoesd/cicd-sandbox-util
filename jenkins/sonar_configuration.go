package jenkins

import (
	"../model"
	"../util"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SonarConfiguration assumes we can call a job in Jenkins for configuring SonarQube.
// It also assumed we can do this via the Build Token Root plugin and the jobConfig's token is the build token.
func SonarConfiguration(sonarQubeConfig model.SonarQubeConfig, jobConfig model.JenkinsJobConfig, token model.ApiToken) {
	fmt.Printf(">> Updating Jenkins' SonarQube Configuration\n")
	// http://localhost:8282/jenkins/job/configs/job/sonar/buildWithParameters
	// when used with root build token (https://wiki.jenkins.io/display/JENKINS/Build+Token+Root+Plugin)
	//		it will be like this: http://localhost:8282/jenkins/buildWithParameters?job=configs/sonar&token=<build_token>&<further params>
	// form data: sonar_token
	// success = 201 Created
	rawUrl := util.ReplaceHostnameAndPort("http://XXX:YYY/ZZZ/buildByToken/buildWithParameters",
		jobConfig.Host,
		jobConfig.Port)
	rawUrl = strings.Replace(rawUrl, "/ZZZ", jobConfig.ContextRoot, 1)
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
