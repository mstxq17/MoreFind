# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
LDFLAGS := -s -w

ifneq ($(shell go env GOOS),darwin)
LDFLAGS := -extldflags "-static"
endif

all: build
build:
	$(GOBUILD) $(GOFLAGS) -ldflags '$(LDFLAGS)' -trimpath -o "MoreFind" main.go
