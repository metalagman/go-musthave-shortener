FROM golang:1.17 AS dependencies
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS build
COPY . /go/src/app
WORKDIR /go/src/app
RUN make build-shortener

FROM debian:buster-slim
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/bin/* /app/
RUN chmod +x /app/*
EXPOSE 80/tcp
CMD /app/shortener -v
