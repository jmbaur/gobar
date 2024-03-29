# vim: ft=make

help:
	@just --list

build:
	go build -o $out/gobar ./cmd/gobar

check: build
	go test ./...
	revive -set_exit_status=1 ./...
	staticcheck ./...

update:
	#!/usr/bin/env bash
	go get -u all
	go mod tidy
	export NIX_PATH="nixpkgs=$(nix flake prefetch nixpkgs --json | jq --raw-output '.storePath')"
	newvendorHash="$(nix-prefetch \
		 "{ sha256 }: ((import <nixpkgs> {}).callPackage ./. {}).goModules.overrideAttrs (_: { vendorHash = sha256; })")"
	sed -i "s|vendorHash.*|vendorHash = \"$newvendorHash\";|" default.nix

run module="full":
	go run ./cmd/gobar --config {{justfile_directory()}}/examples/{{module}}.yaml

test: build
	sway --config ./test/sway.config
