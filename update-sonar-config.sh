#!/usr/bin/env bash
docker run --rm -e EXTERNAL_HOSTNAME=localhost --network=cidc_default cicd-sandbox-util:0.1.0 cicd-util -action sonar-init \
    -keycloakClientId sonarqube \
    -keycloakHost 172.17.0.1\
    -externalPort 8289\
    -contextRoot "/sonar"\
    -internalPort 9000 \
    -internalHost=sonar