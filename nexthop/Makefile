# Copyright (c) Tetrate, Inc 2021 All Rights Reserved.

BIN := nexthop
OUT := bin

$(BIN): $(OUT)/$(BIN)
$(BIN)-static: $(OUT)/$(BIN)-static

$(OUT)/$(BIN): $(OUT)
	go build -o $@ main.go

$(OUT)/$(BIN)-static: $(OUT)
	CGO_ENABLED=0 GOOS=linux go build \
		-a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo \
		-o $@ main.go

$(OUT):
	mkdir -p $@

docker-build: $(BIN)-static
	docker build -t $(HUB)/$(BIN) -f Dockerfile .

docker-push:
	docker push $(HUB)/$(BIN)

clean:
	rm -rf $(OUT)

.PHONY: docker-build docker-push clean
