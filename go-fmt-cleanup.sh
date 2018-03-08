#!/usr/bin/env bash
gofmt -s -w -l main.go
gofmt -s -w -l model/*.*
gofmt -s -w -l webserver/*.*
gofmt -s -w -l jenkins/*.*
gofmt -s -w -l dockerprobe/*.*
gofmt -s -w -l sonarqube/*.*
gofmt -s -w -l util/*.*