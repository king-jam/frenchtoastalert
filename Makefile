build:
	go build -o bin/scraper main.go

build-server:
	go build -o bin/server data/main.go