
.PHONY: list-outdated-deps
list-outdated-deps:
	@go list -m -u -f '{{if and (not .Indirect) .Update }}{{.}}{{end}}' all 2> /dev/null

.PHONY: bump-deps
bump-deps:
	@scripts/bump-dependencies.sh

.PHONY: image
image:
	if [ ! -f version ]; then echo "1.0.0" > version; fi
	docker build --tag "gstack/keyval-resource:v`< version`" .
	docker tag "gstack/keyval-resource:v`< version`" "gstack/keyval-resource:latest"

.PHONY: publish
publish:
	docker push "gstack/keyval-resource:v`< version`"
	docker push "gstack/keyval-resource:latest"
