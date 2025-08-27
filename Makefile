build:
	go build -o ./dist/

test:
	go test ./... -v

docs-gen:
	tfplugindocs generate
	pre-commit run markdownlint --all-files
