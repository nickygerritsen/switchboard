VERSION ?= $(shell git describe --tags --always 2>/dev/null | sed 's/^v//')
VERSION := $(if $(VERSION),$(VERSION),1.0.2)

srpm:
	@echo "Building SRPM for version: $(VERSION)"
	sed 's/__VERSION__/$(VERSION)/g' switchboard.spec > switchboard.spec.tmp
	rpmbuild -bs --define "_sourcedir $(CURDIR)" \
		--define "_specdir $(CURDIR)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		switchboard.spec.tmp
	rm -f switchboard.spec.tmp
