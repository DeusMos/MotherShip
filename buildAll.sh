#!/bin/bash
go build bin/portal/portal.go 
go build bin/client/client.go 
mv client rmd
GOARCH=386 GOOS=linux go build bin/remote/remote.go 