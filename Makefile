SUB_PACKAGE  := $(subst $(shell git rev-parse --show-toplevel),,$(CURDIR))
PACKAGE_ROOT := $(shell git remote -v | grep '^origin\s.*(fetch)$$' | awk '{print $$2}' | sed -E 's/^.*(\/\/|@)(.*)\.git\/?$$/\2/' | sed 's/:/\//g')
PACKAGE_NAME = $(PACKAGE_ROOT)$(SUB_PACKAGE)

APP_ROOT := $(shell dirname $(PACKAGE_NAME))
APP      := $(shell basename $(PACKAGE_NAME))
GROUP    := $(shell dirname $(APP_ROOT))

OUTPUT              = $(CURDIR)/output
CONF                = $(CURDIR)/conf
TEMPLATE_DIR        = $(CURDIR)/tpl
STATIC_DIR        	= $(CURDIR)/static

OUTPUT_LIB_DIR = $(OUTPUT)/lib

OUTPUT_DIRS = conf tpl bin static

BUILD_ROOT   := $(shell git rev-parse --show-toplevel)/build
BUILD_TARGET = src/$(PACKAGE_ROOT)
BUILD_DIR    = $(BUILD_ROOT)/$(BUILD_TARGET)

export GOPATH=$(BUILD_ROOT)
export GOBIN=$(BUILD_ROOT)/bin

.DEFAULT: all
all: build

build: clean prepare
	cd "$(BUILD_DIR)" && go build -o "$(OUTPUT)/bin/$(APP)" "$(BUILD_DIR)$(SUB_PACKAGE)/main.go"

clean:
	for i in $(OUTPUT_DIRS); do rm -rf "$(OUTPUT)/$$i"; done
	git checkout -- $(RANK_SEARCH_GO_FILE) $(RANK_REC_GO_FILE)

prepare:
	for i in $(OUTPUT_DIRS); do mkdir -p "$(OUTPUT)/$$i"; done
	cp -vr "$(CONF)" "$(OUTPUT)"
	cp -vr "$(TEMPLATE_DIR)" "$(OUTPUT)"
	cp -vr "$(STATIC_DIR)" "$(OUTPUT)"
	cp -v "$(CURDIR)/control.sh" "$(OUTPUT)"
	cp -v "$(CURDIR)/deploy.sh" "$(OUTPUT)"

run:
	cd "$(OUTPUT)" && bin/$(APP)

init:
	sed -i'' -e 's/^package:.*/package: $(subst /,\/,$(PACKAGE_NAME))/' "$(GLIDE_YAML)"

	mkdir -p "$(shell dirname $(BUILD_DIR))"
	if [ ! -e "$(BUILD_DIR)" ]; then ln -s "$(shell echo $(BUILD_TARGET) | sed -E 's/[a-zA-Z0-9_.-]+/../g')" "$(BUILD_DIR)"; fi

glide-up: glide-update
glide-update:
	glide update

glide-i: glide-install
glide-install:
	glide install

.PHONY: all build clean prepare run init glide-up glide-update glide-i glide-install
$(VERBOSE).SILENT: