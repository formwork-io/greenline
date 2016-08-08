prefix		= /usr/local
exec_prefix	= $(prefix)
bindir		= $(exec_prefix)/bin
top 		= $(abspath $(dir .))
builddir	= $(top)/build
depdir		= $(builddir)/deps
libdepdir	= $(builddir)/deps/libs
pkgdir		= $(depdir)/pkgs
grnl_bin 	= $(builddir)/grnl
grnlctl_bin = $(builddir)/grnlctl
go_install 	= go install -v -x -pkgdir $(pkgdir)
go_build 	= go build -v -x -pkgdir $(pkgdir)
core_object = $(pkgdir)/github.com/formwork-io/greenline/src/core.a

INSTALL			= /usr/bin/install
INSTALL_PROGRAM	= $(INSTALL)

export GOPATH = $(top)
SHELL := /usr/bin/env bash

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

ifndef PKG_CONFIG_PATH
export PKG_CONFIG_PATH = $(libdepdir)/lib/pkgconfig
endif

ifndef CFLAGS
OS := $(shell uname)
ifeq ($(OS),Darwin)
cflags = -L$(libdepdir)/lib -lsodium -lzmq -lczmq -lstdc++
else
cflags = -L$(libdepdir)/lib -lsodium -lzmq -lczmq -lstdc++ -static
endif
else
cflags = $(CFLAGS)
endif
go_ldflags = -ldflags '--extldflags "$(cflags)"'

all: grnl grnlctl

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) \
		| sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1:\2/' \
		| column -c2 -t -s :)"

$(libdepdir)/.done: env.sh scripts/get-lib-deps.sh
	./scripts/get-lib-deps.sh

libdeps: $(libdepdir)/.done ## Get library dependencies via get-lib-deps.sh

$(grnl_bin): src/grnl/*.go $(core_object)
	cd src/grnl && $(go_build) $(go_ldflags) -o $(@)
	@touch $(@)
grnl: $(grnl_bin) ## Builds the grnl binary

$(grnlctl_bin): src/grnlctl/*.go $(core_object)
	cd src/grnlctl && $(go_build) $(go_ldflags) -o $(@)
	@touch $(@)
grnlctl: $(grnlctl_bin) ## Builds the grnlctl binary

$(core_object): src/core/*.go
	cd src/core && $(go_install)
core: $(core_object) ## Builds the core code

clean: grnlctl-clean grnl-clean core-clean ## Cleans the build

grnlctl-clean:
	@rm -fv $(grnlctl_bin)

grnl-clean:
	@rm -fv $(grnl_bin)

core-clean:
	@rm -fv $(core_object)

mostlyclean: clean ## Cleans the build and local packages
	$(info Removing $(pkgdir))
	@rm -fr $(pkgdir)

distclean: clean ## Cleans everything
	$(info Removing $(builddir))
	@rm -fr $(builddir)

install: grnl grnlctl ## Installs the binaries
	@# see Makefile install command categories
	$(NORMAL_INSTALL)
	$(INSTALL_PROGRAM) -m 0755 \
		$(grnl_bin) $(grnlctl_bin) $(DESTDIR)$(bindir)

install-strip: ## Installs stripped versions of the binaries
	@# see Makefile install command categories
	$(NORMAL_INSTALL)
	$(INSTALL_PROGRAM) -m 0755 -s \
		$(grnl_bin) $(grnlctl_bin) $(DESTDIR)$(bindir)

uninstall: ## Uninstalls the binaries
	@rm -fv $(DESTDIR)$(bindir)/grnl \
		    $(DESTDIR)$(bindir)/grnlctl

.PHONY: clean distclean grnlctl-clean grnl-clean core-clean help install \
		install-strip uninstall
