build: clean
    go build -o ./content .

clean:
    rm -f ./content

run *args: build
    ./content {{ args }}

ssg:
    ./content render
    
serve:
    python3 -m http.server -d public 8080

