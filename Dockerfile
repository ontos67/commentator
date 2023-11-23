FROM golang:1.15 AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download

COPY . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/commentator.go


FROM debian:stretch

COPY --from=builder /code/commentator /usr/local/bin/commentator

RUN chmod +x /usr/local/bin/commentator

ENTRYPOINT [ "commentator" ]