#!/usr/bin/env bash
gofmt -s -w -l main.go
gofmt -s -w -l model/*.*
gofmt -s -w -l webserver/*.*
