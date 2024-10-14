FROM golang:1.22

WORKDIR /server

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

CMD ["./app"]