main: clean build compress

build: main.go
	golint
	go test .
	go build .


compress:
	upx go-demo

clean:
	rm -f go-demo

