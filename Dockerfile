FROM golang:alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app ./cmd/bot


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
