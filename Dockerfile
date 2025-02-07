FROM golang:1.22.3-alpine AS builder

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

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs .

# static files, unlikely to change
RUN mkdir -p ./assets
COPY ./assets ./assets

# copy the config
COPY ./config.yaml .