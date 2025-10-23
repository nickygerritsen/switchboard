srpm:
	dnf install -y git golang
	git archive --format=tar.gz --prefix=switchboard-$(VERSION)/ HEAD > switchboard-$(VERSION).tar.gz
	rpmbuild -bs --define "_sourcedir $(CURDIR)" \
		--define "_specdir $(CURDIR)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		--define "version $(VERSION)" \
		switchboard.spec
