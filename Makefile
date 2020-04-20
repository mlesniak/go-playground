main: clean build compress

build: main.go
	golint -set_exit_status
	go test .
	go build .


compress:
	upx go-demo

clean:
	rm -f go-demo

