FROM golang:1.17.1-alpine3.14 AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY app.env .
RUN go mod download

COPY . /build
WORKDIR /build/
RUN go build

FROM alpine:3.14
COPY --from=builder /build/tradingdata /build/app.env /app/
COPY --from=builder /build/internal/db/migration /app/migration/
WORKDIR /app
EXPOSE 8080
CMD ["./tradingdata"]