# Copyright 2017 Heptio Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TARGET = eventrouter
GOTARGET = github.com/heptiolabs/$(TARGET)
BUILDMNT = /src/
REGISTRY ?= steveww
VERSION ?= 0.0.1
IMAGE = $(REGISTRY)/$(BIN)
DOCKER ?= docker
DIR := ${CURDIR}

ifneq ($(VERBOSE),)
VERBOSE_FLAG = -v
endif
TESTARGS ?= $(VERBOSE_FLAG) -timeout 60s
TEST_PKGS ?= $(GOTARGET)/sinks/...
TEST = go test $(TEST_PKGS) $(TESTARGS)
VET_PKGS ?= $(GOTARGET)/...
VET = go vet $(VET_PKGS)


all: build

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)

image:
	$(DOCKER) build -t $(REGISTRY)/$(TARGET) .
	$(DOCKER) tag $(REGISTRY)/$(TARGET) $(REGISTRY)/$(TARGET):$(VERSION)

push:
	$(DOCKER) push $(REGISTRY)/$(TARGET)
	if git describe --tags --exact-match >/dev/null 2>&1; \
	then \
		$(DOCKER) push $(REGISTRY)/$(TARGET):$(VERSION); \
	fi

test:
	$(DOCKER_BUILD) '$(TEST)'

vet:
	$(DOCKER_BUILD) '$(VET)'

.PHONY: all local build image push

clean:
	rm -f $(TARGET)
	$(DOCKER) rmi $(REGISTRY)/$(TARGET):latest
	$(DOCKER) rmi $(REGISTRY)/$(TARGET):$(VERSION)
