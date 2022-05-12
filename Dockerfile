FROM golang:1.17.1-alpine3.14 AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY app.env .
RUN go mod download

COPY . /build
WORKDIR /build/
RUN go build -o /tradingdata

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /tradingdata .
COPY app.env .
COPY internal/db/migration /app/migration/
EXPOSE 8080
CMD ["./tradingdata"]