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
core_object = $(pkgdir)/core.a

INSTALL			= /usr/bin/install
INSTALL_PROGRAM	= $(INSTALL)

export GOPATH = $(top)
SHELL := /usr/bin/env bash

ifndef PKG_CONFIG_PATH
export PKG_CONFIG_PATH = $(libdepdir)/lib/pkgconfig
endif

HAVE_PKG_CONFIG := $(shell command -v pkg-config)
HAVE_GO := $(shell command -v go)

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

all: cmddeps grnl grnlctl

cmddeps:
ifndef HAVE_PKG_CONFIG
	$(error "pkg-config is not available")
endif

ifndef HAVE_GO
	$(error "go is not available")
endif

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) \
		| sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1:\2/' \
		| column -c2 -t -s :)"

$(libdepdir)/.done: cmddeps env.sh scripts/get-lib-deps.sh
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
	-rm -f $(grnlctl_bin)

grnl-clean:
	-rm -f $(grnl_bin)

core-clean:
	-rm -f $(core_object)

mostlyclean: clean ## Cleans the build and local packages
	-rm -fr $(pkgdir)

distclean: ## Cleans everything
	-rm -fr $(builddir)

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
	-rm -f $(DESTDIR)$(bindir)/grnl \
		    $(DESTDIR)$(bindir)/grnlctl

.PHONY: clean distclean grnlctl-clean grnl-clean core-clean help install \
		install-strip uninstall
