NAME = armor-web
IMAGE = labstack/armor-web 

build:
	rm -rf public
	hugo
	docker build -t $(IMAGE) .

push: build
	docker push $(IMAGE)

.PHONY: build push