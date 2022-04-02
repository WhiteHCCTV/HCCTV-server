FROM golang:1.18

WORKDIR /go/apps/HCCTV/src

CMD ["go", "run", "main.go"]
