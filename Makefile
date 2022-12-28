APP = chat
SERVICE_PATH = github.com/Dsmit05
BR = `git rev-parse --symbolic-full-name --abbrev-ref HEAD`
VER = `git describe --tags --abbrev=0`
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
FORMAT = $(VER)-$(TIMESTM)
DOCTAG = $(VER)-$(BR)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP) cmd/chat/main.go

.PHONY: build-image
build-image:
	docker build -t $(APP):$(DOCTAG) .

.PHONY: run-app
run-app:
	docker run -d --name=$(APP)-$(VER) $(APP):$(DOCTAG)

.PHONY: del-app
del-app:
	docker rm $(APP)-$(VER)
