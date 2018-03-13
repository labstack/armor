IMAGE = labstack/armor
VERSION = 0.4.0-dev

clean:
	rm -rf build

build: clean
	xgo --pkg cmd/armor -out build/armor-$(VERSION) github.com/labstack/armor
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE) .

push: build
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

.PHONY: clean build push
