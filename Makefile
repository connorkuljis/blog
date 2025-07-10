list-commands:
	@cat Makefile

site-build-release: site-clean
	@go run -v ./cmd/site/main.go

site-build-draft:
	@go run -v ./cmd/site/main.go -d

site-debug:
	go tool dlv debug ./cmd/site -- -d

site-watch:
	@find assets/ templates/ cmd/ internal/ | entr make site-build-draft

site-deploy: site-build-release
	@echo "Syncing local assets to staging directory..."
	@rsync --delete -avz \
		dist/ \
		prod@kuljis.xyz:/home/prod/www/kuljis.xyz/dist/

site-clean:
	@rm -rf ./dist

cms-build:
	@rm -f ./cms
	go build -v ./cmd/cms

cms-debug:
	go tool dlv debug ./cmd/cms

goimports:
	go tool goimports -w -l .

serve:
	python3 -m http.server -d ./dist 3000

