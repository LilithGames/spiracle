.PHONY: build
build:
	@GOOS=linux go build -o bin/ github.com/LilithGames/spiracle/...

.PHONY: run
run: build
	@wsl -e bin/spiracle
