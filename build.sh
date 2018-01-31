#!/bin/sh

exec go build -o build/astral-hash -ldflags "-X main.version=$(git describe --tags --always)" main.go