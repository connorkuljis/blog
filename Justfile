bin := "./content"

clean:
    rm -f {{ bin }}

build: clean
    go build -o {{ bin }} .

run *args: build
    {{ bin }} {{ args }}

render:
    {{ bin }} render
    
serve:
    python3 -m http.server -d public 8080

