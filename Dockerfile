FROM golang:1.15.7 as build

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest
COPY --from=build /app/ledger-auth /app/

EXPOSE 8080

ENTRYPOINT ["/app/ledger-auth"]
