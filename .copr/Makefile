# Try multiple methods to extract version
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null | sed 's/^v//')
ifeq ($(VERSION),)
VERSION := $(shell git describe --tags 2>/dev/null | sed 's/^v//' | cut -d'-' -f1)
endif
ifeq ($(VERSION),)
VERSION := $(shell git tag --points-at HEAD 2>/dev/null | grep '^v' | sed 's/^v//' | head -1)
endif
ifeq ($(VERSION),)
$(error VERSION could not be determined from git tags)
endif

srpm:
	@echo "Building SRPM for version: $(VERSION)"
	sed 's/__VERSION__/$(VERSION)/g' switchboard.spec > $(outdir)/switchboard.spec
	rpmbuild -bs --define "_sourcedir $(CURDIR)" \
		--define "_specdir $(outdir)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		$(outdir)/switchboard.spec
