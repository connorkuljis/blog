cli := "./cli"
gen := "./ssg"

public := "./public"
port := "3000"

all: clean build-cli build-gen

build-cli: 
    go build -v -o {{ cli }} ./cmd/cli/main.go

build-gen: 
    go build -v -o {{ gen }} ./cmd/ssg/main.go

clean:
    rm -f {{ cli }} {{ gen}}

serve:
    python3 -m http.server -d {{ public }} {{ port }}

