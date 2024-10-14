FROM golang:1.22

WORKDIR /client

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/client

CMD ["./app", "-addr", "server:9999"]