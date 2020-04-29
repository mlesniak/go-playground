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
	#docker build -f Dockerfile-filebeat -t mlesniak/filebeat .
	docker build --build-arg COMMIT=`git rev-parse HEAD` -t mlesniak/go-demo .

push:
	#docker push mlesniak/filebeat
	docker push mlesniak/go-demo

docker-run:
	docker-compose up

