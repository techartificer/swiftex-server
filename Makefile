#!make
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o swiftex

.PHONY: runbuild
runbuild:
	./swiftex
