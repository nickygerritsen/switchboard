VERSION ?= $(shell git describe --tags --always 2>/dev/null | sed 's/^v//' || echo "1.0.2")

srpm:
	rpmbuild -bs --define "_sourcedir $(CURDIR)" \
		--define "_specdir $(CURDIR)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		--define "version $(VERSION)" \
		switchboard.spec
