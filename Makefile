IMAGE = labstack/armor
VERSION = 0.3.7
WORKDIR = /go/src/github.com/labstack/armor

clean:
	rm -rf build

build: clean
	docker run -v `pwd`:$(WORKDIR) -w $(WORKDIR) golang:1.9 sh -c '\
		GOOS=linux GOARCH=amd64 go build -o build/armor-$(VERSION)_linux-64 cmd/armor/main.go \
		&& GOOS=linux GOARCH=arm go build -o build/armor-$(VERSION)_linux-arm32 cmd/armor/main.go \
    && GOOS=linux GOARCH=arm64 go build -o build/armor-$(VERSION)_linux-arm64 cmd/armor/main.go \
    && GOOS=darwin GOARCH=amd64 go build -o build/armor-$(VERSION)_darwin-64 cmd/armor/main.go \
    && GOOS=windows GOARCH=amd64 go build -o build/armor-$(VERSION)_windows-64.exe cmd/armor/main.go'
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE) .

push: build
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

.PHONY: clean build push
