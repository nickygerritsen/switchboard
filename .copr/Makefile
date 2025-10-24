srpm:
	@echo "=== DEBUG: Git environment ==="
	@echo "Current directory: $$(pwd)"
	@echo "Git branch/commit: $$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'detached')"
	@echo "Git commit SHA: $$(git rev-parse HEAD 2>/dev/null)"

	@echo ""
	@echo "=== DEBUG: Fetching tags ==="
	@git fetch --tags 2>&1 || echo "git fetch --tags failed"

	@echo ""
	@echo "=== DEBUG: Available tags ==="
	@git tag -l 2>&1 || echo "git tag -l failed"

	@echo ""
	@echo "=== DEBUG: Tags pointing at HEAD ==="
	@git tag --points-at HEAD 2>&1 || echo "git tag --points-at HEAD failed"

	@echo ""
	@echo "=== DEBUG: git describe attempts ==="
	@echo "Method 1 (exact match): $$(git describe --tags --exact-match 2>&1 | head -1)"
	@echo "Method 2 (describe): $$(git describe --tags 2>&1 | head -1)"
	@echo "Method 3 (points-at): $$(git tag --points-at HEAD | grep '^v' | head -1)"

	@echo ""
	@echo "=== DEBUG: Extracting version ==="
	@VERSION=$$(git describe --tags --exact-match 2>/dev/null | sed 's/^v//') && \
	if [ -z "$$VERSION" ]; then \
		VERSION=$$(git describe --tags 2>/dev/null | sed 's/^v//' | sed 's/-.*//'); \
	fi && \
	if [ -z "$$VERSION" ]; then \
		VERSION=$$(git tag --points-at HEAD 2>/dev/null | grep '^v' | sed 's/^v//' | head -1); \
	fi && \
	if [ -z "$$VERSION" ]; then \
		echo "ERROR: VERSION could not be determined from git tags"; \
		exit 1; \
	fi && \
	echo "Extracted VERSION: $$VERSION" && \
	echo "" && \
	echo "=== Building SRPM for version: $$VERSION ===" && \
	sed "s/__VERSION__/$$VERSION/g" switchboard.spec > $(outdir)/switchboard.spec && \
	rpmbuild -bs --define "_sourcedir $(CURDIR)" \
		--define "_specdir $(outdir)" \
		--define "_builddir $(CURDIR)" \
		--define "_srcrpmdir $(outdir)" \
		--define "_rpmdir $(outdir)" \
		$(outdir)/switchboard.spec
