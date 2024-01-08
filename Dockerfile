FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/server
CMD ["./app"]