FROM golang:alpine as build
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apk update && apk add --no-cache ca-certificates
RUN go build -o app
EXPOSE 8080
CMD ["./app"]