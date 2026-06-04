FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git
WORKDIR /workspace

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o alidns-webhook -ldflags '-s -w' .


FROM alpine

COPY --from=builder /workspace/alidns-webhook /usr/local/bin/alidns-webhook

ENTRYPOINT ["alidns-webhook"]
