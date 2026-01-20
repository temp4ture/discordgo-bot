.PHONY: _ _no_param _chalk dependencies exec b
# make with no target defaults to "default".
_: _no_param linux

SRC_FILE := "./main.go"
OUTPUT_PATH := "./.bin"
EXPORT_FILENAME := "discordgo-bot"

_no_param:
	@echo "no target provided; building for linux."

_chalk:
	@command -v go >/dev/null 2>&1 || { \
		echo "golang is not installed."; \
		echo "install go here: https://golang.org/dl/"; \
		exit 1; \
	}
	@echo Checking for dependencies...
	@go get
	@go mod tidy
.env:
	@echo "\".env\" file does not exist in root! cannot continue."
	@echo "BOT_TOKEN = token_goes_here # https://discord.com/developers/applications" > .env
	@echo "created \".env\", update the file and run 'make' once more."
	@echo "                  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ "
	@echo 
	@echo 
	@echo 
	@exit 1

# get dependencies
dependencies: _chalk .env

exec: dependencies
	@echo "Executing \"main.go\"..."
	-go run main.go ${PARAMETERS} ||:
b: dependencies
	@echo "Building to \"${OUT_PATH}\"" 
	@${ENV} go build -o "${OUT_PATH}" ${SRC_FILE}

.PHONY: run run-no-terminal run-verbose
# note: i suggest using these only for checks and debugging
# 	 	do "go run main.go (parameters)" otherwise.
run:
	@$(MAKE) exec TARGET=$@ PARAMETERS="${pmt}"
run-no-terminal:
	@$(MAKE) exec TARGET=$@ PARAMETERS="--no-terminal"
run-verbose:
	@$(MAKE) exec TARGET=$@ PARAMETERS="--verbose"

.PHONY: build build-exe
linux:
	@echo Building project into linux executable...
	@$(MAKE) b TARGET=$@ ENV="" OUT_PATH="${OUTPUT_PATH}/linux/${EXPORT_FILENAME}"
windows:
	@echo Building project into windows executable...
	@$(MAKE) b TARGET=$@ ENV="GOOS=windows GOARCH=386" OUT_PATH="${OUTPUT_PATH}/windows/${EXPORT_FILENAME}.exe"