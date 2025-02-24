FROM golang:latest as base

WORKDIR /app

FROM base as build 

COPY go.mod go.sum ./

COPY . .

RUN go build -o main

FROM alpine:latest as final

WORKDIR /app

COPY --from=build /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]