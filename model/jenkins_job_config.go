package model

// JenkinsJobConfig information about a Jenkins job we want to interact with.
type JenkinsJobConfig struct {
	User        string
	Token       string
	Host        string
	Port        string
	ContextRoot string
	JobUrl      string
	JobParams   string
}
