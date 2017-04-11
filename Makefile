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

# Used only on function 'release', generate one binary per couple os/arch
RELEASE_OS			?= linux darwin
RELEASE_ARCH		?= 386 amd64

# Temporary directories to use to generate binaries and documentation
DIR_PROJECT			:= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DIR_BUILD			:= $(DIR_PROJECT)/build
DIR_COVERAGE		:= $(DIR_PROJECT)/coverage
DIR_RELEASE			:= $(DIR_PROJECT)/release
DIR_RELEASE_TMP		:= $(DIR_PROJECT)/.tmp/

# GTK version to use, see https://github.com/gotk3/gotk3/wiki
GTK_VERSION			?= $(shell pkg-config --modversion gtk+-3.0 | sed -E 's/([0-9]+)\.([0-9]+).*/gtk_\1_\2/')
GTK_TAG				:= -tags $(GTK_VERSION)

# Used to fill the /version api endpoint
BUILD_VERSION		:= $(shell git describe --tags --always --dirty="-dev")
BUILD_TIME			:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
BUILD_FLAGS			:= -ldflags='-X "main.BuildVersion=$(BUILD_VERSION)" -X "main.BuildTime=$(BUILD_TIME)"' $(GTK_TAG)

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
	$Q echo -e '$(COLOR_PRINT)Building $(DIR_BUILD)/bin/$(BINARY_NAME)...$(COLOR_RESET)'
	$Q mkdir -p $(DIR_BUILD)/bin
	$Q go build -i -v -o $(DIR_BUILD)/bin/$(BINARY_NAME) $(BUILD_FLAGS)
	$Q echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

# Synchronize vendors and tools
vendor:
	$Q echo -e '$(COLOR_PRINT)Syncing tools...$(COLOR_RESET)'
	$Q retool sync
	$Q echo -e '$(COLOR_PRINT)Syncing vendors...$(COLOR_RESET)'
	$Q retool do govendor sync -v
	$Q retool do govendor test -i $(GTK_TAG)
	$Q echo -e '$(COLOR_PRINT)Syncing linters...$(COLOR_RESET)'
	$Q retool do gometalinter --install --update --force
	$Q echo -e '$(COLOR_SUCCESS)Synchronization done without errors$(COLOR_RESET)'

# Synchronize vendors and tools
vendor_clean:
	$Q echo -e '$(COLOR_PRINT)Cleaning vendors...$(COLOR_RESET)'
	$Q echo -e '$(COLOR_PRINT)Cleaning vendored tools...$(COLOR_RESET)'
	$Q rm -rf _tools
	$Q echo -e '$(COLOR_PRINT)Cleaning vendored sources...$(COLOR_RESET)'
	$Q find vendor/* -maxdepth 0 -type d -exec rm -r {} +
	$Q echo -e '$(COLOR_SUCCESS)Cleaned$(COLOR_RESET)'

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

# Remove all non-essentials directories and files
clean: vendor_clean
	$Q echo -e '$(COLOR_PRINT)Cleaning...$(COLOR_RESET)'
	$Q rm -rf $(DIR_BUILD) $(DIR_RELEASE) $(DIR_RELEASE_TMP) $(DIR_COVERAGE)
	$Q echo -e '$(COLOR_SUCCESS)Cleaned$(COLOR_RESET)'

# Run code documentation server
godoc:
	$Q echo -e '$(COLOR_PRINT)Open a web browser and load 127.0.0.1:6060 ...$(COLOR_RESET)'
	$Q godoc -http=:6060 -index
	$Q echo -e '$(COLOR_PRINT)Terminated$(COLOR_RESET)'

# Check for useless and missing dependencies
test_dependencies:
	$Q echo -e '$(COLOR_PRINT)Testing dependencies...$(COLOR_RESET)'
	$Q retool do govendor list +unused +missing
	@[ "$(shell retool do govendor list +unused +missing | wc -l)" = "0" ]
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check syntax, format, useless, and non-optimized code
test_code:
	$Q echo -e '$(COLOR_PRINT)Testing code with linters...$(COLOR_RESET)'
	$Q find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec gofmt -d {} \;
	@[ $(shell find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec gofmt -d {} \; | wc -l) = 0 ]
	$Q ./_scripts/linter.sh $(GTK_TAG)
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check unit tests
test_unit:
	$Q echo -e '$(COLOR_PRINT)Testing code with unit tests...$(COLOR_RESET)'
	$Q retool do govendor test +local -v -timeout 5s $(GTK_TAG) ./...
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# TODOs should never exist
test_todo:
	$Q echo -e '$(COLOR_PRINT)Testing presence of TODOs in code...$(COLOR_RESET)'
	$Q find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec grep -Hn "//TODO:" {} \;
	@[ "$(shell find . -name vendor -prune -o -name _tools -prune -o -name "*.go" -exec grep -Hn "//TODO:" {} \; | wc -l)" = "0" ]
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

# Check all kind of tests
test: test_dependencies test_code test_unit test_todo

# Compute coverage and create coverage files
coverage:
	$Q echo -e '$(COLOR_PRINT)Generating test converage...$(COLOR_RESET)'
	$Q rm -rf $(DIR_COVERAGE)
	$Q mkdir -p $(DIR_COVERAGE)
	$Q echo "mode: $(TEST_COVERAGE_MODE)" > $(DIR_COVERAGE)/coverage.out
	$Q for pkg in $(shell retool do govendor list -no-status +local); do \
		go test -covermode="$(TEST_COVERAGE_MODE)" -coverprofile="$(DIR_COVERAGE)/coverage.tmp" "$$pkg" 2>&1 > /dev/null; \
		grep -h -v "^mode:" $(DIR_COVERAGE)/coverage.tmp >> $(DIR_COVERAGE)/coverage.out 2> /dev/null; \
	done
	$Q rm -f $(DIR_COVERAGE)/coverage.tmp
	$Q go tool cover -func=$(DIR_COVERAGE)/coverage.out
	$Q echo -e '$(COLOR_SUCCESS)Done$(COLOR_RESET)'

coverage_show: coverage
	$Q go tool cover -html=$(DIR_COVERAGE)/coverage.out

# Compile and save all necessaries file for each couple os/arch
release_build: test coverage config
	$Q echo -e '$(COLOR_PRINT)List of files beeing compiled:$(COLOR_RESET)'
	$Q go list -f '{{.GoFiles}}' ./...
	$Q mkdir -p $(DIR_BUILD)/bin
	$Q echo -e '$(COLOR_PRINT)Building...$(COLOR_RESET)'
	$Q for os in $(RELEASE_OS); do \
		for arch in $(RELEASE_ARCH); do \
			echo "Building $(BINARY_NAME) for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -o $(DIR_BUILD)/bin/$(BINARY_NAME)-$$os$$arch $(BUILD_FLAGS); \
		done; \
	done
	$Q echo -e '$(COLOR_SUCCESS)Compilation done without errors$(COLOR_RESET)'

# Generate a release ; multiple files are going to be generated:
#	a documentation archive, a runnable environment archive for each couple os/arch
release: config clean doc release_build
	$Q [ -n "$(TAG)" ] || (echo "Please add the release tag with the TAG=x.x.x environment variable" && false)
	$Q mkdir -p $(DIR_RELEASE) $(DIR_RELEASE_TMP)/$(BINARY_NAME)
	$Q cp $(DIR_BUILD)/bin/* $(DIR_RELEASE_TMP)/$(BINARY_NAME)
	$Q cp $(CONFIGURATION_FILE) $(DIR_RELEASE_TMP)/$(BINARY_NAME)/config.ini
	$Q tar --create --gzip --file=$(DIR_RELEASE)/$(BINARY_NAME)-doc-$(TAG).tar.gz -C $(DIR_BUILD)/ doc/
	$Q for os in $(RELEASE_OS); do \
		echo "Creating $(BINARY_NAME)-$$os archives..."; \
		cd $(DIR_RELEASE_TMP) && tar --verbose --create --gzip --file=$(DIR_RELEASE)/$(BINARY_NAME)-bin-$(TAG)-$$os.tar.gz \
			$(BINARY_NAME)/config.ini $(BINARY_NAME)/$(BINARY_NAME)-$$os*; \
	done


.PHONY: all $(BINARY_NAME) vendor build config run clean godoc test_dependencies test_code test_unit test_todo test coverage coverage_show release_build release
