package main

import (
	"./dockerprobe"
	"./jenkins"
	"./model"
	"./sonarqube"
	"./util"
	"./webserver"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	dockerHost := flag.String("dockerHost", "unix:///var/run/docker.sock", "The docker host to use")
	filterDockerRegister := flag.Bool("filterDockerRegister", false, "Whether or not to filter the docker image names for the serve action")
	filterDockerRegisterName := flag.String("filterDockerRegisterName", "N/A", "The name of the docker registry to filter of the image name (iff filterDockerRegister = true)")
	labelPrefix := flag.String("labelPrefix", "com.abnamro", "The (docker) label prefix")
	namespace := flag.String("namespace", "cicd-sandbox", "The namespace to filter for in Docker containers")
	keycloakHost := flag.String("keycloakHost", "localhost", "The keycloak hostname")
	keycloakPort := flag.String("keycloakPort", "8280", "Port number for the Keycloak server")
	keycloakClientId := flag.String("keycloakClientId", "", "The clientId for to use for the tool you're trying to configure")
	securityRealm := flag.String("securityRealm", "ci", "The security realm to use for Keycloak configuration")
	externalPort := flag.String("externalPort", "0", "The external port for the tool you're trying to configure (for connection through proxy)")
	externalHost := flag.String("externalHost", "localhost", "The external hostname for the tool you're trying to configure (for connection through proxy)")
	internalPort := flag.String("internalPort", "0", "The internal port for the tool you're trying to configure (for connection within docker network)")
	internalHost := flag.String("internalHost", "localhost", "The internal hostname for the tool you're trying to configure (for connection within docker network)")
	contextRoot := flag.String("contextRoot", "", "The context root of the tool you're trying to configure")
	sonarQubeTokenName := flag.String("sonarQubeTokenName", "ci", "The tokenName for the API token for SonarQube (for example, for Jenkins)")
	sonarQubeUser := flag.String("sonarQubeUser", "admin", "The username for SonarQube API")
	sonarQubePass := flag.String("sonarQubePass", "admin", "The password for SonarQube API")
	jenkinsUser := flag.String("jenkinsUser", "barbossa", "The username for Jenkins API")
	jenkinsToken := flag.String("jenkinsToken", "JENKINS_ROCKS", "The user's token for Jenkins API")
	jenkinsHost := flag.String("jenkinsHost", "localhost", "The host for Jenkins API")
	jenkinsPort := flag.String("jenkinsPort :=", "8282", "The port for Jenkins API")
	jenkinsContextRoot := flag.String("jenkinsContextRoot :=", "/jenkins", "The context root of the Jenkins")
	jenkinsJobUrl := flag.String("jenkinsJobUrl", "/configs/sonar", "The url of the jenkins job to trigger")
	jenkinsJobParams := flag.String("jenkinsJobParams", "", "comma delimited key:value pairs")

	serverPort := flag.String("serverPort", "7777", "The Port number of the webserver when action is 'serve'")
	action := flag.String("action", "generate-config", `
		- generate-config: Generate configuration files such as keycloak configuration for Jenkins
		- sonar-init: initialize the configuration of SonarQube, such as the keycloak configuration
  		- sonar-token: generates a token for ci systems (such as jenkins) in SonarQube
		- jenkins-sonar-token: generates a SonarQube security token (internalHost/Port...) and triggers a configuration job in Jenkins 
		- sonar-token-list: lists sonar tokens
		- list-docker: list docker containers part of this stack
		- serve: serve as web server serving a html page with the docker container listing (same source as list-docker)
		`)
	flag.Parse()

	sonarQubeConfig := model.SonarQubeConfig{
		ExternalHostname: *externalHost,
		ExternalPort:     *externalPort,
		SecurityRealm:    *securityRealm,
		KeycloakClientId: *keycloakClientId,
		ContextRoot:      *contextRoot,
		InternalHostname: *internalHost,
		InternalPort:     *internalPort,
		SonarQubeUser:    *sonarQubeUser,
		SonarQubePass:    *sonarQubePass,
		APITokenName:     *sonarQubeTokenName,
	}

	jenkinsJobConfig := model.JenkinsJobConfig{
		User:        *jenkinsUser,
		Token:       *jenkinsToken,
		Host:        *jenkinsHost,
		Port:        *jenkinsPort,
		ContextRoot: *jenkinsContextRoot,
		JobUrl:      *jenkinsJobUrl,
		JobParams:   *jenkinsJobParams,
	}

	// add functions
	// SonarQube -> generate security token for Jenkins
	// Nexus LDAP config
	// Nexus repo?

	hostname, _ := os.Hostname()
	if len(os.Getenv("EXTERNAL_HOSTNAME")) > 0 {
		hostname = os.Getenv("EXTERNAL_HOSTNAME")
	}
	fmt.Printf("== EXTERNAL_HOSTNAME=%v\n", os.Getenv("EXTERNAL_HOSTNAME"))

	hostname = strings.ToLower(hostname) // windows might not care, but Keycloak certainly does!
	fmt.Printf("== Hostname to use: %s\n", hostname)
	fmt.Printf("== Action to perform: %s\n", *action)
	switch *action {
	case "generate-config":
		fmt.Printf("== Keycloak Config for Jenkins\n-----------------\n")
		jenkins.KeycloakConfig(*keycloakHost, *keycloakPort, *securityRealm)
	case "sonar-init":
		fmt.Printf("== Keycloak Config for SonarQube\n-----------------\n")
		sonarKeycloakConfig := sonarqube.GenerateSonarKeycloakConfig(*keycloakHost, *keycloakPort, sonarQubeConfig)
		fmt.Printf("-----------------\n")
		fmt.Printf("== Update Config of SonarQube\n-----------------\n")
		fmt.Printf("== keycloakPort=%s\n", *keycloakPort)
		fmt.Printf("== keycloakHost=%s\n", *keycloakHost)
		printSonarQubeConfig(sonarQubeConfig)
		sonarqube.UpdateSonarQubeConfig(sonarKeycloakConfig, sonarQubeConfig)
	case "sonar-token":
		fmt.Printf("== Generate CI token in Sonarqube -- start\n-----------------\n")
		printSonarQubeConfig(sonarQubeConfig)
		sonarqube.GenerateAPIToken(sonarQubeConfig)
		fmt.Printf("== Generate CI token in Sonarqube -- end\n-----------------\n")
	case "sonar-token-list":
		fmt.Printf("== Generate CI token in Sonarqube -- start\n-----------------\n")
		printSonarQubeConfig(sonarQubeConfig)
		sonarqube.RetrieveAPITokens(sonarQubeConfig)
		fmt.Printf("== Generate CI token in Sonarqube -- end\n-----------------\n")
	case "jenkins-sonar-token":
		fmt.Printf("== Configure Sonarqube in Jenkins -- start\n-----------------\n")
		printSonarQubeConfig(sonarQubeConfig)
		printJenkinsConfig(jenkinsJobConfig)
		tokens := sonarqube.RetrieveAPITokens(sonarQubeConfig)
		fmt.Printf("  > Found %d tokens\n", len(tokens.UserTokens))
		tokenExists := false
		for _, userToken := range tokens.UserTokens {
			if userToken.Name == *sonarQubeTokenName {
				tokenExists = true
			}
		}
		if tokenExists {
			fmt.Println(" > Token already exists, we cannot retrieve it, so not updating Jenkins")
		} else {
			fmt.Println(" > Token does not exist yet, we will create it")
			token := sonarqube.GenerateAPIToken(sonarQubeConfig)
			jenkins.SonarConfiguration(sonarQubeConfig, jenkinsJobConfig, token)
		}
		fmt.Printf("== Configure Sonarqube in Jenkins -- end\n-----------------\n")
	case "list-docker":
		labelFilter := fmt.Sprintf("%s=%s", *labelPrefix+util.LabelNamespace, *namespace)
		containersList, err := dockerprobe.ContainerList(labelFilter, *dockerHost)
		containers := dockerprobe.ContainerInfoList(containersList, *filterDockerRegister, *filterDockerRegisterName, *labelPrefix)

		if err != nil {
			fmt.Println(err.Error())
		} else if len(containers) == 0 {
			fmt.Printf("   > No Containers found with label filter %s\n", labelFilter)
		} else {
			fmt.Printf(" > We found these containers: \n")
			for _, container := range containers {
				if container.Name != "" {
					fmt.Printf("   > %s\n", container.String())
				}
			}
		}
	case "serve":

		fmt.Printf("=== STARTING WEB SERVER @%s\n", *serverPort)
		fmt.Println("=============================================")

		labelFilter := fmt.Sprintf("%s=%s", *labelPrefix+util.LabelNamespace, *namespace)
		containersList, _ := dockerprobe.ContainerList(labelFilter, *dockerHost)
		containers := dockerprobe.ContainerInfoList(containersList, *filterDockerRegister, *filterDockerRegisterName, *labelPrefix)
		webserverData := &webserver.WebserverData{Containers: containers, Title: "CICD Sandbox Containers"}

		c := make(chan bool)
		go webserver.StartServer(*serverPort, webserverData, c)
		fmt.Println("> Started the web server, now polling swarm")

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		for i := 1; ; i++ { // this is still infinite
			t := time.NewTicker(time.Second * 30)
			select {
			case <-stop:
				fmt.Println("> Shutting down polling")
				break
			case <-t.C:
				fmt.Println("  > Updating Stacks")
				containersList, _ := dockerprobe.ContainerList(labelFilter, *dockerHost)
				containers := dockerprobe.ContainerInfoList(containersList, *filterDockerRegister, *filterDockerRegisterName, *labelPrefix)
				webserverData.UpdateContainers(containers)
				continue
			}
			break // only reached if the quitCh case happens
		}
		fmt.Println("> Shutting down webserver")
		c <- true
		if b := <-c; b {
			fmt.Println("> Webserver shut down")
		}
		fmt.Println("> Shut down app")
	default:
		panic(fmt.Sprintf(
			"Action '%v' not recognized\n", *action))
	}
	fmt.Printf("-----------------\n")
}

// printJenkinsConfig prints the basic information captured for Jenkins.
func printJenkinsConfig(jenkinsJobConfig model.JenkinsJobConfig) {
	fmt.Printf("======= Jenkins Job configuration\n")
	fmt.Printf("== User=%s\n", jenkinsJobConfig.User)
	fmt.Printf("== Host=%s\n", jenkinsJobConfig.Host)
	fmt.Printf("== Port=%s\n", jenkinsJobConfig.Port)
	fmt.Printf("== ContextRoot=%s\n", jenkinsJobConfig.ContextRoot)
	fmt.Printf("== JobUrl=%s\n", jenkinsJobConfig.JobUrl)
	fmt.Printf("== JobParams=%s\n", jenkinsJobConfig.JobParams)
	fmt.Printf("======= Jenkins Job configuration\n")
}

// printSonarQubeConfig prints the basic information captured for SonarQube.
func printSonarQubeConfig(sonarQubeConfig model.SonarQubeConfig) {
	fmt.Printf("======= SonarQube configuration\n")
	fmt.Printf("== KeycloakClientId=%s\n", sonarQubeConfig.KeycloakClientId)
	fmt.Printf("== ExternalHostname=%s\n", sonarQubeConfig.ExternalHostname)
	fmt.Printf("== ExternalPort=%s\n", sonarQubeConfig.ExternalPort)
	fmt.Printf("== InternalHostname=%s\n", sonarQubeConfig.InternalHostname)
	fmt.Printf("== InternalPort=%s\n", sonarQubeConfig.InternalPort)
	fmt.Printf("== ContextRoot=%s\n", sonarQubeConfig.ContextRoot)
	fmt.Printf("== SonarQubeUser=%s\n", sonarQubeConfig.SonarQubeUser)
	fmt.Printf("== APITokenName=%s\n", sonarQubeConfig.APITokenName)
	fmt.Printf("======= SonarQube configuration\n")
}
