.PHONY: _ _no_param _chalk dependencies exec b
# make with no target defaults to "default".
_: _no_param linux

SRC_FILE := "./main.go"
OUTPUT_PATH := "./.bin"
EXPORT_FILENAME := "discordgo-bot"
PATH_LINUX := "${OUTPUT_PATH}/linux"
PATH_WIN := "${OUTPUT_PATH}/windows"

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
	@${MAKE} b_env PATH='.env' --silent
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
b_env:
	@echo "BOT_TOKEN = token_goes_here # https://discord.com/developers/applications" > "${PATH}"

.PHONY: run run-no-terminal run-verbose
# i suggest using these only for checks and debugging
# do "go run main.go (parameters)" otherwise.
run:
	@$(MAKE) exec TARGET=$@ PARAMETERS="${pmt}" --silent
run-no-terminal:
	@$(MAKE) exec TARGET=$@ PARAMETERS="--no-terminal" --silent
run-verbose:
	@$(MAKE) exec TARGET=$@ PARAMETERS="--verbose" --silent

.PHONY: build build-exe
linux:
	@echo Building project into linux executable
	@$(MAKE) b TARGET=$@ ENV="" OUT_PATH="${PATH_LINUX}/${EXPORT_FILENAME}" --silent
	@echo Creating \".env\" file
	@${MAKE} b_env PATH="${PATH_LINUX}/.env" --silent
windows:
	@echo Building project into windows executable
	@$(MAKE) b TARGET=$@ ENV="GOOS=windows GOARCH=386" OUT_PATH="${PATH_WIN}/${EXPORT_FILENAME}.exe" --silent
	@echo Creating \".env\" file
	@${MAKE} b_env PATH="${PATH_WIN}/.env" --silent