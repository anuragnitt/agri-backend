#!/bin/bash

protoc --go_out=generated --proto_path=proto --go_opt=paths=source_relative --go-grpc_out=generated --proto_path=proto --go-grpc_opt=paths=source_relative     proto/*.proto

go run *.go
