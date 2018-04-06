IMAGE = labstack/armor
VERSION = 0.4.4

clean:
	rm -rf build

build: clean
	xgo --targets=darwin-10.8/amd64,linux/amd64,linux/arm-6,linux/arm-7,linux/arm64,windows-8.0/amd64 --pkg cmd/armor -out build/armor-$(VERSION) github.com/labstack/armor
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE) .

push: build
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

.PHONY: clean build push
