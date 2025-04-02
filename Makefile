all: site cms

site:
	go run -v ./cmd/site

cms: clean
	go build -v ./cmd/cms

clean:
	rm -f ./cms

serve:
	./scripts/serve.sh
