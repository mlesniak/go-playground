main: clean dependencies build

build:
	go build cmd/go-demo/main.go

dependencies:
	go get -u github.com/labstack/echo/...
	go get "github.com/sirupsen/logrus"

compress:
	upx go-demo

clean:
	rm -f go-demo

run:
	go run .

docker:
	docker build -t mlesniak/go-demo .

docker-run:
	docker run --rm -it -p 8080:8080 mlesniak/go-demo

