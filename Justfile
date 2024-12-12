build: clean
    go build -o ./bin/content .

clean:
    rm -f ./bin/content

run: build
    ./bin/content entries
    
serve:
    python3 -m http.server -d public 8080

