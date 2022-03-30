FROM golang:1.18

WORKDIR /go/src/app

COPY main.go go.mod go.sum ./

CMD ["go", "run", "main.go"]
