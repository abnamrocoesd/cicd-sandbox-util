#!/usr/bin/env bash
TOKEN_NAME_DEFAULT=jenkins
TOKEN_NAME="${1:-$TOKEN_NAME_DEFAULT}"

cicd-sandbox-util -action jenkins-sonar-token -internalPort 8289 -contextRoot /sonar -sonarQubeTokenName ${TOKEN_NAME}
