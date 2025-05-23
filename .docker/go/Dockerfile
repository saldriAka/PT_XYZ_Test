FROM golang:1.24.3

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

CMD ["go", "run", "main.go"]

