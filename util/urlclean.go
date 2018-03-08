package util

import "strings"

func ReplaceHostnamePortAndRealm(originalString string, hostname string, port string, securityRealm string) string {
	returnString := ReplaceHostnameAndPort(originalString, hostname, port)
	return strings.Replace(returnString, "ZZZ", securityRealm, 1)
}

func ReplaceHostnameAndPort(originalString string, hostname string, port string) string {
	returnString := strings.Replace(originalString, "XXX", hostname, 1)
	return strings.Replace(returnString, "YYY", port, 1)
}