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
REGISTRY ?= steveww
VERSION ?= 0.0.1
DOCKER ?= docker
KUBECTL ?= kubectl
DIR := ${CURDIR}


.PHONY: all build image push deploy test vet clean

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

deploy:
	$(KUBECTL) apply -f manifests/eventrouter.yaml

test:
	go test ./...

vet:
	go vet ./...


clean:
	rm -f $(TARGET)
	$(DOCKER) rmi $(REGISTRY)/$(TARGET):latest
	$(DOCKER) rmi $(REGISTRY)/$(TARGET):$(VERSION)
