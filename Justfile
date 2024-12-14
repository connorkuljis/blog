build: clean
    go build -o ./bin/content .

clean:
    rm -f ./bin/content

run *sub: build
    ./bin/content {{ sub }}
    
serve:
    python3 -m http.server -d public 8080

