docs-build:
	redocly bundle api/openapi.yaml --output api/bundled-schema.yml
	redocly build-docs api/bundled-schema.yml --output=ui/docs/index.html
	rm api/bundled-schema.yml

docs:
	redocly preview-docs ./api/openapi.yaml