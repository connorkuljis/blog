bin := "./content"

list:
    @just --list

clean:
    rm -f {{ bin }}

build: clean
    go build -o {{ bin }} .

run *args: 
    {{ bin }} {{ args }}

serve:
    python3 -m http.server -d public 8080

