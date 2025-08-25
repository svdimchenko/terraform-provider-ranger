build:
	go build -o ./dist/

docs-gen:
	tfplugindocs generate
	pre-commit run markdownlint --all-files
