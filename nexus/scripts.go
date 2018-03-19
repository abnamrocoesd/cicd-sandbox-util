package nexus

import (
	"../util"
	"fmt"
	"net/url"
	"strings"
)

//
//
// POST /service/rest/v1/script
//curl -u admin:admin123 -X POST --header 'Content-Type: application/json' \
//http://localhost:8081/service/rest/v1/script \  -d @helloWorld.json
func AddScript() {
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
