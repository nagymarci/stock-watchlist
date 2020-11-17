FROM golang:alpine AS builder

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./cmd/server

FROM alpine:latest

RUN apk add --no-cache bash tzdata

COPY --from=builder /build/main .

EXPOSE 3300
CMD ["./main"]