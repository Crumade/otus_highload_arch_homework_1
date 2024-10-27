FROM golang:1.23.2-bookworm
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /social_network main.go
EXPOSE 8080
CMD ["/social_network"]