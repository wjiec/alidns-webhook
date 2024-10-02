FROM golang:1.22-alpine AS build_deps

RUN apk add --no-cache git

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download


FROM build_deps AS builder

COPY . .

RUN CGO_ENABLED=0 go build -o alidns-webhook -ldflags '-s -w' .


FROM alpine:3.17

COPY --from=builder /workspace/alidns-webhook /usr/local/bin/alidns-webhook

ENTRYPOINT ["alidns-webhook"]
