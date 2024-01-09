FROM golang:latest AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/app ./cmd/server

FROM mysql:latest
ENV MYSQL_ROOT_PASSWORD=1111
ENV MYSQL_DATABASE=user_posts_db
ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=1111
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
EXPOSE 8080
CMD ["app"]