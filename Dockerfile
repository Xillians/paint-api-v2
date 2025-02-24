ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


FROM debian:bookworm

# Need to add ca-certificates or we will see the following error:
# x509: certificate signed by unknown authority
RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
