REGISTRY         ?= docker.io
PROJECT          ?= ansibleplaybookbundle
TAG              ?= latest
BROKER_IMAGE     = $(REGISTRY)/$(PROJECT)/ansible-service-broker
BUILD_DIR        = "${GOPATH}/src/github.com/openshift/ansible-service-broker/build"

install: $(shell find cmd pkg)
	go install ./cmd/broker

${GOPATH}/bin/mock-registry: $(shell find cmd/mock-registry)
	go install ./cmd/mock-registry

# Will default run to dev profile
run: install vendor
	@${GOPATH}/src/github.com/openshift/ansible-service-broker/scripts/runbroker.sh dev

deploy:
	@${GOPATH}/src/github.com/openshift/ansible-service-broker/scripts/deploy.sh

run-mock-registry: ${GOPATH}/bin/mock-registry vendor
	@${GOPATH}/src/github.com/openshift/ansible-service-broker/cmd/mock-registry/run.sh

prepare-build: install
	cp "${GOPATH}"/bin/broker build/broker

build: prepare-build
	docker build ${BUILD_DIR} -t ${BROKER_IMAGE}:${TAG}
	@echo
	@echo "Remember you need to push your image before calling make deploy"
	@echo "    docker push ${BROKER_IMAGE}:${TAG}"

clean:
	@rm -f ${GOPATH}/bin/broker
	@rm -f build/broker

vendor:
	@glide install -v

test: vendor
	go test ./pkg/...

.PHONY: run run-mock-registry clean test build asb-image install prepare-build
