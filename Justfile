clean:
    rm -f ./main

build: clean
    go build -o ./main .

run: build
    ./main

