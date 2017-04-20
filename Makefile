# Binary name to generate
BINARY_NAME			:= nebulo-client-desktop

# Used for generation and copy in archive
CONFIGURATION_FILE	:= config.sample.json

# Used for go coverage tools
TEST_COVERAGE_MODE	?= count

# Overload this variable on make call 'make <function> CI=1' to add debug information
#	and remove terminal colors
CI					?= 0

# Overload this variable on make call `make <function> ARGS="help" to run with custom arguments`
ARGS				?= -c config.json run

# Temporary directories to use to generate binaries and documentation
DIR_PROJECT			:= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DIR_BUILD			:= $(DIR_PROJECT)/build
DIR_COVERAGE		:= $(DIR_PROJECT)/coverage
DIR_RELEASE			:= $(DIR_PROJECT)/release
DIR_RELEASE_TMP		:= $(DIR_PROJECT)/.tmp/

# GTK version to use, see https://github.com/gotk3/gotk3/wiki
GTK_VERSION			?= $(shell pkg-config --modversion gtk+-3.0 2>/dev/null | sed -E 's/([0-9]+)\.([0-9]+).*/gtk_\1_\2/' || echo "gtk_unknown")

# Used to fill the version and some useful flags
BUILD_VERSION		:= $(shell git describe --tags --always --dirty="-dev")
BUILD_TIME			:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
BUILD_FLAGS			:= -ldflags='-X "main.BuildVersion=$(BUILD_VERSION)" -X "main.BuildTime=$(BUILD_TIME)"'
BUILD_TAGS			:= -tags $(GTK_VERSION)

# CI variable define how informations are displayed on console
ifeq ($(CI),0)
	# Don't show command, and make everything 'pretty'
	Q				:= @
	COLOR_SUCCESS	:= \e[0;32m # green
	COLOR_FAIL		:= \e[0;31m # red
	COLOR_PRINT		:= \e[0;33m # orange
	COLOR_RESET		:= \e[0m
else
	# Show everything, make everything readable in a file
	Q				:= $(shell echo "")
	COLOR_SUCCESS	:= $Q
	COLOR_FAIL		:= $Q
	COLOR_PRINT		:= $Q
	COLOR_RESET		:= $Q
endif


all : clean vendor build test

# Compile for current os/arch and save binary in $DIR_BUILD folder
$(BINARY_NAME):
	$Q echo -e '$(COLOR_PRINT)Building $(DIR_BUILD)/bin/$(BINARY_NAME) with $(GTK_VERSION)...$(COLOR_RESET)'
	$Q mkdir -p $(DIR_BUILD)/bin
	$Q go build -i -v -o $(DIR_BUILD)/bin/$(BINARY_NAME) $(BUILD_FLAGS) $(BUILD_TAGS)
	$Q echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

build: $(BINARY_NAME)

# Generate configuration file
config: $(BINARY_NAME)
	$Q echo -e '$(COLOR_PRINT)Generating $(CONFIGURATION_FILE)...$(COLOR_RESET)'
	$Q $(shell $(DIR_BUILD)/bin/$(BINARY_NAME) -v quiet config-gen -d $(CONFIGURATION_FILE))
	$Q echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

# Compile for current os/arch and run binary
run: $(BINARY_NAME)
	$Q echo -e '$(COLOR_PRINT)Running $(BINARY_NAME):$(COLOR_RESET)'
	$Q $(DIR_BUILD)/bin/$(BINARY_NAME) ${ARGS}
	$Q echo -e '$(COLOR_PRINT)Terminated$(COLOR_RESET)'

# Synchronize vendors and tools
vendor:
	$Q echo -e '$(COLOR_PRINT)Syncing tools...$(COLOR_RESET)'
	$Q retool sync
	$Q echo -e '$(COLOR_PRINT)Syncing vendors...$(COLOR_RESET)'
	$Q retool do govendor sync -v
	$Q retool do govendor test -i $(BUILD_TAGS)
	$Q echo -e '$(COLOR_PRINT)Syncing linters...$(COLOR_RESET)'
	$Q retool do gometalinter --install --update --force
	$Q echo -e '$(COLOR_SUCCESS)Synchronization done without errors$(COLOR_RESET)'

# Clean vendors and tools
vendor-clean:
	$Q echo -e '$(COLOR_PRINT)Cleaning vendors...$(COLOR_RESET)'
	$Q echo -e '$(COLOR_PRINT)Cleaning vendored tools...$(COLOR_RESET)'
	$Q rm -rf _tools
	$Q echo -e '$(COLOR_PRINT)Cleaning vendored sources...$(COLOR_RESET)'
	$Q find vendor/* -maxdepth 0 -type d -exec rm -r {} +
	$Q echo -e '$(COLOR_SUCCESS)Cleaned$(COLOR_RESET)'

# Remove all non-essentials directories and files
clean: vendor-clean
	$Q echo -e '$(COLOR_PRINT)Cleaning...$(COLOR_RESET)'
	$Q rm -rf $(DIR_BUILD) $(DIR_RELEASE) $(DIR_RELEASE_TMP) $(DIR_COVERAGE)
	$Q echo -e '$(COLOR_SUCCESS)Cleaned$(COLOR_RESET)'

# Create and run a docker image containing all needed libraries
docker-build:
	$Q docker build -t nebulo:client-desktop .
	$Q docker stop nebulo_client_desktop || true
	$Q docker run --name nebulo_client_desktop -d --rm -e DISPLAY=$(DISPLAY) --net=host -v /tmp/.X11-unix:/tmp/.X11-unix:ro -v $(DIR_PROJECT):/go/src/github.com/krostar/nebulo-client-desktop nebulo:client-desktop tail -f /dev/null
	$Q xhost +

# Run the code inside the docker container
docker-run:
	$Q docker exec nebulo_client_desktop make run

# Exec an arbitrary command in the container
docker-exec:
	$Q docker exec -it nebulo_client_desktop $(CMD)

# Check for useless and missing dependencies
test-dependencies:
	$Q echo -e '$(COLOR_PRINT)Testing dependencies...$(COLOR_RESET)'
	$Q retool do govendor list +unused +missing
	@[ $(shell retool do govendor list +unused +missing | wc -l) = 0 ]
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check syntax, format, useless, and non-optimized code
test-code:
	$Q echo -e '$(COLOR_PRINT)Testing code with linters...$(COLOR_RESET)'
	$Q find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec gofmt -d {} \;
	@[ $(shell find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec gofmt -d {} \; | wc -l) = 0 ]
	$Q ./_scripts/linter.sh $(BUILD_TAGS)
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check unit tests
test-unit:
	$Q echo -e '$(COLOR_PRINT)Testing code with unit tests...$(COLOR_RESET)'
	$Q retool do govendor test +local -v -timeout 5s $(BUILD_TAGS) ./...
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# TODOs should never exist
test-todo:
	$Q echo -e '$(COLOR_PRINT)Testing presence of TODOs in code...$(COLOR_RESET)'
	$Q find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec grep -Hn "//TODO:" {} \;
	@[ $(shell find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec grep -Hn "//TODO:" {} \; | wc -l) = 0 ]
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check all kind of tests
test: test-dependencies test-code test-unit test-todo

# Compute coverage and create coverage files
coverage:
	$Q echo -e '$(COLOR_PRINT)Generating test converage...$(COLOR_RESET)'
	$Q rm -rf $(DIR_COVERAGE)
	$Q mkdir -p $(DIR_COVERAGE)
	$Q echo "mode: $(TEST_COVERAGE_MODE)" > $(DIR_COVERAGE)/coverage.out
	$Q # go test -covermode="$(TEST_COVERAGE_MODE)" -coverprofile="$(DIR_COVERAGE)/coverage.tmp" $(BUILD_TAGS) "$$pkg" 2>&1 > /dev/null; \
	$Q #grep -h -v "^mode:" $(DIR_COVERAGE)/coverage.tmp >> $(DIR_COVERAGE)/coverage.out 2> /dev/null; \
	$Q for pkg in $(shell retool do govendor list -no-status +local); do \
		go test -covermode="$(TEST_COVERAGE_MODE)" -coverprofile="$(DIR_COVERAGE)/coverage.tmp" -tags $(BUILD_TAGS) "$$pkg"; \
		(test -e $(DIR_COVERAGE)/coverage.tmp && grep -h -v "^mode:" $(DIR_COVERAGE)/coverage.tmp >> $(DIR_COVERAGE)/coverage.out) || true; \
	done
	$Q rm -f $(DIR_COVERAGE)/coverage.tmp
	$Q go tool cover -func=$(DIR_COVERAGE)/coverage.out
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

.PHONY: all $(BINARY_NAME) build config run vendor vendor-clean clean docker-build docker-run docker-exec test-dependencies test-code test-unit test-todo test coverage
