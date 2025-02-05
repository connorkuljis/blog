cli := "./cli"
gen := "./ssg"

public := "./public"
port := "3000"

all: clean build-cli build-ssg

build-cli: 
    go build -v -o {{ cli }} ./cmd/cli/main.go

build-ssg: 
    go build -v -o {{ gen }} ./cmd/ssg/main.go

clean:
    rm -f {{ cli }} {{ gen}}

serve:
    python3 -m http.server -d {{ public }} {{ port }}

