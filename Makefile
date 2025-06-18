all: site cms

release:
	go run -v ./cmd/site/main.go

draft:
	go run -v ./cmd/site/main.go -d

cms: clean
	go build -v ./cmd/cms

clean:
	rm -f ./cms
goimports:
	go tool goimports -w -l .

	rm -rf ./dist

serve:
	python3 -m http.server -d ./dist 3000

deploy:
	echo "Syncing local assets to staging directory..."
	rsync --delete -avz \
		dist/ \
		prod@kuljis.xyz:/home/prod/www/kuljis.xyz/dist/
