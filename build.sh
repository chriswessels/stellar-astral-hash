#!/bin/sh

exec go build -o bin/astral-hash -ldflags "-X main.version=$(git describe --tags --always)" main.go