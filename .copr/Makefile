srpm:
	@echo "=== Building SRPM ==="
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION environment variable not set"; \
		exit 1; \
	fi
	@echo "Using VERSION: $(VERSION)"
	@sed "s/__VERSION__/$(VERSION)/g" switchboard.spec > $(outdir)/switchboard.spec && \
	rpmbuild -bs --define "_sourcedir $(CURDIR)" \
		--define "_specdir $(outdir)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		$(outdir)/switchboard.spec
