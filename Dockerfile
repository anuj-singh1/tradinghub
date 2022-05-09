FROM golang:1.14.9-alpine AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY app.env .
RUN go mod download

COPY . /build
WORKDIR /build/
RUN go build

FROM alpine
COPY --from=builder /build/tradingdata /build/app.env /app/
WORKDIR /app
EXPOSE 8080
CMD ["./tradingdata"]