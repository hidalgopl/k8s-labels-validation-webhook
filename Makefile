NAME=$(shell basename $(TOP_DIR))
CMD=k8s-labels-validation-webhook
GIT_COMMIT=$(shell git rev-parse --short HEAD)
CURRENT_TIME=$(shell date +%Y%m%d%H%M%S)

TOP_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
MAIN_SRC=*.go
DKR=docker
#PKG_SRC=$(shell find $(TOP_DIR)pkg -type f -name '*.go')
VERSION=$(shell git describe --dirty 2>/dev/null || echo "dev")
LD_VERSION_FLAGS=-X main.BuildVersion=$(GIT_COMMIT) -X main.BuildTime=$(CURRENT_TIME)
LDFLAGS=-ldflags "$(LD_VERSION_FLAGS)"
REGISTRY_REPO=hidalgopl/k8s-labels-validation-webhook
IMAGE_TAG=$(GIT_COMMIT)
ifneq ($(findstring -,$(GIT_COMMIT)),)
IMAGE_DEV_OR_LATEST=dev
else
IMAGE_DEV_OR_LATEST=latest
endif

all: $(CMD)

$(CMD): $(MAIN_SRC)
	CGO_ENABLED=0 GOOS=linux go build $(LDFLAGS) -o $(TOP_DIR)$(CMD)

img: $(BINARIES)
	@echo "Checking if IMAGE_TAG is set" && test -n "$(IMAGE_TAG)"
	$(DKR) build -t $(REGISTRY_REPO):$(IMAGE_TAG) \
		-t $(REGISTRY_REPO):$(IMAGE_DEV_OR_LATEST) .

push-img: img
	@echo "Checking if IMAGE_TAG is set" && test -n "$(IMAGE_TAG)"
	$(DKR) push $(REGISTRY_REPO):$(IMAGE_TAG)
	$(DKR) push $(REGISTRY_REPO):$(IMAGE_DEV_OR_LATEST)

lint:
	golangci-lint run -v

test:
	go test ./... -v

clean:
	go clean

install-cert-manager:
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.2/cert-manager.yaml

deploy-cert:
	kubectl apply -f deployment/cert.yaml

deploy:
	kubectl apply -f deployment/deployment.yaml deployment/validating-webhook-template.yaml

.PHONY: all clean

