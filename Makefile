site-release: clean-dist
	go run -v ./cmd/site/main.go

site-draft: clean-dist
	go run -v ./cmd/site/main.go -d


site-watch:
	find assets/ templates/ cmd/ internal/ | entr make site-draft

cms-build:
	rm -f ./cms
	go build -v ./cmd/cms

goimports:
	go tool goimports -w -l .

clean-dist:
	rm -rf ./dist

serve:
	python3 -m http.server -d ./dist 3000

deploy:
	echo "Syncing local assets to staging directory..."
	rsync --delete -avz \
		dist/ \
		prod@kuljis.xyz:/home/prod/www/kuljis.xyz/dist/
