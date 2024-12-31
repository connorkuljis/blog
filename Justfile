cli := "./cli"
gen := "./gen"

public := "./public"
port := "8080"

all: clean build-cli build-gen

build-cli: 
    go build -o {{ cli }} ./cmd/cli/main.go

build-gen: 
    go build -o {{ gen }} ./cmd/gen/main.go

clean:
    rm -f {{ cli }} {{ gen}}

serve:
    python3 -m http.server -d {{ public }} {{ port }}

