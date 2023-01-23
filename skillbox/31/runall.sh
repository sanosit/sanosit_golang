#!/bin/sh

export PORT=8082
go run cmd/main.go &
export PORT=8084
go run cmd/main.go &
export PORT=8080
export WORKERS=8082:8084
go run cmd/proxy.go &
