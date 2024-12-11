build: clean
    go build -o ./bin/content .

clean:
    rm -f ./bin/content

run: build
    ./bin/content entries
    
       

