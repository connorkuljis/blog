# Variables
GO := go
PYTHON := python3
DIST_DIR := ./dist
SITE_CMD := ./cmd/site/main.go
CMS_CMD := ./cmd/cms
CMS_BINARY := ./cms
DEPLOY_HOST := prod@connorkuljis.com:/home/prod/www/connorkuljis.com/dist/
SERVE_PORT := 3000
WATCH_DIRS := assets/ templates/ cmd/ internal/

# Targets
.PHONY: site-build-draft site-build-release site-debug site-watch site-deploy site-clean cms-build cms-debug goimports serve

site-build-draft: site-clean
	@$(GO) run -v $(SITE_CMD) -d

site-build-release: site-clean
	@$(GO) run -v $(SITE_CMD)

site-debug:
	$(GO) tool dlv debug $(SITE_CMD) -- -d

site-watch:
	@find $(WATCH_DIRS) | entr make site-build-draft

site-deploy: site-build-release
	@echo "Syncing local assets to staging directory..."
	@rsync --delete -avz \
		$(DIST_DIR)/ \
		$(DEPLOY_HOST)

site-clean:
	@rm -rf $(DIST_DIR)

cms-build:
	@rm -f $(CMS_BINARY)
	$(GO) build -v $(CMS_CMD)

cms-debug:
	$(GO) tool dlv debug $(CMS_CMD)

goimports:
	$(GO) tool goimports -w -l .

serve:
	$(PYTHON) -m http.server -d $(DIST_DIR) $(SERVE_PORT)
