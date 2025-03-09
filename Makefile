site:
	go run ./cmd/site

cms: clean
	go build -v ./cmd/cms

clean:
	rm -f ./cms
