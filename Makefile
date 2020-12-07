.PHONY: build clean all binary
.SILENT: build binary

ARCH ?= amd64
CONTAINER_KSERVICE_HOME ?= /go/src/github.com/jinghzhu/kservice
ALL_ARCH = amd64 dawrin
GO_VERSION ?= 1.15
GO_IMAGE ?= golang
CONTAINER_KSERVICE_GOBIN ?= $(CONTAINER_KSERVICE_HOME)/bin/$(ARCH)
KSERVICE_CMDS = github.com/jinghzhu/kservice/cmd/kservice

all: build docs binary

build:
	for command in $(KSERVICE_CMDS) ; do \
		echo "building $$command......."; \
		docker run --rm -u $$(id -u):$$(id -g) -v $$(pwd):$(CONTAINER_KSERVICE_HOME) \
			-it $(GO_IMAGE):$(GO_VERSION) \
			/bin/sh -c "\
				mkdir -p $(CONTAINER_KSERVICE_GOBIN) && \
				GOBIN=$(CONTAINER_KSERVICE_GOBIN) go install $$command " && \
		BIN=$$(basename $$command) && echo "Generated bin/$(ARCH)/$$BIN" ; \
	done

	echo "copy swagger.json to bin/$(ARCH)";
	cp $(SWAGGER_SPEC) bin/$(ARCH)
	
binary:
	for command in $(KSERVICE_CMDS) ; do  \
		echo "building container image for target $$command" ; \
		BIN=$$(basename $$command) && \
	    		docker build --build-arg SWAGGER_SPEC=cmd/$$BIN/docs/swagger.json --build-arg BIN=bin/$(ARCH)/$$BIN  -t kservice-$$BIN . && echo "kservice-$$BIN image was built successfully" ; \
	done

clean:
	rm -rf $(pwd)/bin/$(ARCH)
	docker rmi kservice
