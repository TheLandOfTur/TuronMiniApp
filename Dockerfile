FROM golang:1.22

WORKDIR ./

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/main.go

EXPOSE 3032

CMD ["/main"]
