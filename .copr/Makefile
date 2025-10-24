srpm:
	@echo "=== Building SRPM ==="
	@echo "DEBUG: Checking for git..."
	@which git || echo "ERROR: git not found in PATH"
	@git --version || echo "ERROR: git command failed"
	@echo ""
	@echo "DEBUG: Current commit and tags..."
	@git rev-parse HEAD || echo "ERROR: git rev-parse failed"
	@git tag -l | head -5 || echo "ERROR: git tag -l failed"
	@echo ""
	@echo "DEBUG: Attempting version extraction..."
	@VERSION=$$(git describe --tags --exact-match 2>/dev/null | sed 's/^v//') && \
	echo "Method 1 (exact-match): $$VERSION" && \
	if [ -z "$$VERSION" ]; then \
		VERSION=$$(git describe --tags 2>/dev/null | sed 's/^v//' | sed 's/-.*//'); \
		echo "Method 2 (describe): $$VERSION"; \
	fi && \
	if [ -z "$$VERSION" ]; then \
		VERSION=$$(git tag --points-at HEAD 2>/dev/null | grep '^v' | sed 's/^v//' | head -1); \
		echo "Method 3 (points-at): $$VERSION"; \
	fi && \
	if [ -z "$$VERSION" ]; then \
		echo "ERROR: VERSION could not be determined from git tags"; \
		exit 1; \
	fi && \
	echo "" && \
	echo "=== Using VERSION: $$VERSION ===" && \
	sed "s/__VERSION__/$$VERSION/g" ../switchboard.spec > $(outdir)/switchboard.spec && \
	rpmbuild -bs --define "_sourcedir $$HOME/rpmbuild/SOURCES" \
		--define "_specdir $(outdir)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		$(outdir)/switchboard.spec
