main: clean build

build: main.go
	golint -set_exit_status
	go test .
	go build .

dependencies:
	go get -u github.com/labstack/echo/...

compress:
	upx go-demo

clean:
	rm -f go-demo

run: 
	go run .
