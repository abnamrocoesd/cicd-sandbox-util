package model

type JenkinsJobConfig struct {
	User        string
	Token       string
	Host        string
	Port        string
	ContextRoot string
	JobUrl      string
	JobParams   string
}
