#!/bin/sh

exec go build -ldflags "-X main.version=$(git describe --tags --always)" main.go -o ./build/astral-hash