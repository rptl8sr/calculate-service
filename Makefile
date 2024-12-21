ENV_FILE := .env

ifeq ("$(wildcard $(ENV_FILE))","")
  ENV_FILE := dumb.env
endif

ifeq ("$(wildcard $(ENV_FILE))","")
  $(error "Neither .env nor dumb.env found!")
endif

include $(ENV_FILE)
export

REPO_NAME := $(shell basename $(CURDIR))
PROJECT := $(CURDIR)
LOCAL_BIN := $(CURDIR)/bin

# Workflow
.PHONY: init
init: go-init git-init blueprint-init

.PHONY: blueprint-init
blueprint-init:
	touch README.md
	touch Dockerfile
	touch Docker-compose.yaml
	mkdir -p bin
	mkdir -p cmd/server && echo 'package main' >> cmd/server/main.go
	mkdir -p internal/app && echo 'package app' >> internal/app/app.go
	mkdir -p internal/config && echo 'package config' >> internal/config/config.go
	mkdir -p internal/constants && echo 'package constants' >> internal/constants/constants.go
	mkdir -p internal/logger && echo 'package logger' >> internal/logger/logger.go
	mkdir -p internal/server && echo 'package server' >> internal/server/server.go
	mkdir -p internal/router && echo 'package router' >> internal/router/router.go
	mkdir -p internal/middleware && echo 'package middleware' >> internal/middleware/middleware.go
	mkdir -p internal/controller && echo 'package controller' >> internal/controller/controller.go
	mkdir -p internal/handlers && echo 'package handlers' >> internal/handlers/handlers.go
	mkdir -p internal/models && echo 'package models' >> internal/models/models.go

# Go
.PHONY: go-init
go-init:
	go mod init $(REPO_NAME)

# GIT
.PHONY: git-init
git-init:
	gh repo create $(USER)/$(REPO_NAME) --private
	git init
	git config user.name "$(USER)"
	git config user.email "$(EMAIL)"
	git add Makefile go.mod
	git commit -m "Init commit"
	git remote add origin git@github.com:$(USER)/$(REPO_NAME).git
	git remote -v
	git push -u origin master


BN ?= dev
.PHONY: git-checkout
git-checkout:
	git checkout -b $(BN)

# LINT
.PHONY: golangci-lint-install
lint-install:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run ./...

# App
.PHONY: mod
mod:
	go mod download

.PHONY: build
build:
	go build -o calculate-app ./cmd/server/main.go

.PHONY: run
run:
	go run ./cmd/server/main.go