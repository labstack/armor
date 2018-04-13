VERSION = 0.1.3

publish:
	git tag $(VERSION)
	git push origin --tags
	goreleaser --rm-dist

.PHONY: publish 
