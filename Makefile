VERSION ?= $(shell cat VERSION)

build:
	docker build --build-arg VERSION=$(VERSION) -t davidcollom/zwiftcal:$(VERSION)  .

local: build
	docker run --rm -ti -v $(PWD)/.cache:/app/cache -p 3000:3000 davidcollom/zwiftcal:$(VERSION)

debug: build
	docker run --rm -ti -v $(PWD):/app/ -v $(PWD)/.cache:/app/cache davidcollom/zwiftcal:$(VERSION) bash

publish: build
	docker push davidcollom/zwiftcal:$(VERSION)
