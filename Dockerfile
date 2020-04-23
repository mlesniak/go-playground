# ----------------------------------------------------
FROM golang:1.14.2-alpine3.11 AS build
WORKDIR /app
RUN apk update && apk add upx
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o maincmd/go-demo/main.go
RUN go build -o main cmd/go-demo/main.go
RUN upx main

# ----------------------------------------------------
FROM alpine:3.7
COPY --from=build /app/main /app/main
WORKDIR /app
EXPOSE 8080
CMD ["/app/main"]