FROM golang:1.23.6-alpine3.21 AS builder

WORKDIR /app

# prepare the dependencies
COPY go.sum go.mod .
RUN go mod download

# copy the swagger docs
RUN mkdir -p ./docs
COPY ./docs docs

# copy the entire source code
RUN mkdir -p internal
COPY ./internal internal
COPY ./main.go .

# build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# final stage
FROM alpine:3.18.4

RUN apk add curl

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs .

# static files, unlikely to change
RUN mkdir -p ./assets ./db
COPY ./assets ./assets
COPY ./db ./db

HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
CMD curl --fail http://0.0.0.0:5000/ping || exit 1
