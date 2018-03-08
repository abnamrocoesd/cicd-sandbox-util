#!/usr/bin/env bash
docker run --rm -e EXTERNAL_HOSTNAME=localhost --network=cidc_default cicd-sandbox-util:0.1.0 cicd-util -action sonar-token \
    -sonarQubeTokenName jenkins3 \
    -externalPort 8289\
    -contextRoot "/sonar"\
    -internalPort 9000 \
    -internalHost=sonar