# CICD Sandbox Util

The CICD Sandbox Utility is a tool that can help with automating a Continuous Integration (CI) and/or Continuous Deliver (CD) environment.

It is aimed at being used with such a sandbox created with [Docker Compose](https://docs.docker.com/compose/) (such as [CIDC](http://github.com/joostvdg/cidc)). 

Although there are better alternatives such as doing so with [Docker Swarm (Viktor Farcic)](http://vfarcic.github.io/jenkins-swarm/#/cover), [Kubernetes (fabric8)](https://fabric8.io/) or with [OpenShift's MiniShift](https://www.openshift.org/minishift/).

The goal of this is to support a local sandbox environment that is very easy to setup on windows, mac or linux.
The above mentioned alternatives do not satisfy this requirement - or at least not to the degree we're interested in.

## What does it do

Currently it is capable of doing two things:

* **Tool Configuration**: help you automate your tool configurations. Including, but not limited to:
    * configure Keycloak in SonarQube
    * configure a SonarQube API token in Jenkins
* **Service Listing**: this allows you to create a home page listing a set of services started with docker compose and their correct URL. Taking into account things such as ports, contextroots.

For all the actions it does, you can run the following:

```bash
docker run --rm --name util-temp abnamrocoesd/cicd-sandbox-util cicd-util --help
```

## How to use

* build with go
* run with docker

