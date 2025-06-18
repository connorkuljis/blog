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
	./scripts/serve.sh

watch:
	./scripts/watch.sh

