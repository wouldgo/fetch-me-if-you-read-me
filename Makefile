SHELL := /bin/sh
OUT := $(shell pwd)/_out
BUILDARCH := $(shell uname -m)
GCC := $(OUT)/$(BUILDARCH)-linux-musl-cross/bin/$(BUILDARCH)-linux-musl-gcc
LD := $(OUT)/$(BUILDARCH)-linux-musl-cross/bin/$(BUILDARCH)-linux-musl-ld
VERSION := 0.0.1
ENTRYPOINT := cmd/server/main.go cmd/server/options.go

include LOCAL_ENV

env:
	$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' LOCAL_ENV))

test: deps
	rm -Rf _out/.coverage;
	go test -timeout 120s -cover -coverprofile=_out/.coverage -v ./...;
	go tool cover -html=_out/.coverage;

clean-compile: clean compile

docker-build:
	./_tools/docker-buildx \
		build \
			--platform linux/arm64 \
			-f cmd/server/Dockerfile \
			-t ghcr.io/wouldgo/fmiyrm:$(VERSION) .

run: deps env
	go run $(ENTRYPOINT)

compile: deps
	CGO_ENABLED=0 \
	go build \
		-ldflags='-extldflags=-static' \
		-a -o _out/fmiyrm $(ENTRYPOINT)

deps: musl
	go mod tidy -v
	go mod download

musl:
	if [ ! -d "$(OUT)/$(BUILDARCH)-linux-musl-cross" ]; then \
		(cd $(OUT); curl --limit-rate 1G -LOk https://musl.cc/$(BUILDARCH)-linux-musl-cross.tgz) && \
		tar zxf $(OUT)/$(BUILDARCH)-linux-musl-cross.tgz -C $(OUT); \
	fi

clean:
	rm -Rf $(OUT) $(BINARY_NAME)
	mkdir -p $(OUT)
	touch $(OUT)/.keep
