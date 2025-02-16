all: clean build-cms build-site

build-cms: 
    go build -v ./cmd/cms/cms.go

build-site: 
    go build -v ./cmd/site/site.go

clean:
    rm -f ./cms ./site

serve:
    python3 -m http.server -d ./public 3000

