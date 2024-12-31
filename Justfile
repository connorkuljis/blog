cli := "./cli"
gen := "./gen"

public := "./public"
port := "8080"

clean:
    rm -f {{ cli }} {{ gen}}

build-cli: 
    go build -o {{ cli }} ./cmd/cli/main.go

build-gen: 
    go build -o {{ gen }} ./cmd/gen/main.go

serve:
    python3 -m http.server -d {{ public }} {{ port }}

