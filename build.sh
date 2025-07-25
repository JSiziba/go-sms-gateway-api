#!/bin/bash

swag init -g main.go -o docs

go build -o go-sms-gateway-api -ldflags="-s -w" main.go