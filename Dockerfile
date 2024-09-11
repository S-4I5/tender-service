FROM golang:1.22.5-alpine AS builder

RUN apk  update && \
    apk  add --no-cache gcc g++ libc-dev

WORKDIR /build

ADD ../go.mod .

COPY .. .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]