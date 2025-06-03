FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./

# Uncomment lines below if using external dependencies
# COPY go.sum ./
# RUN go mod download

# COPY application code
COPY . .

# Build Go application
RUN go build -o main ./cmd/server/main.go

EXPOSE 4000 

CMD ["./main"]