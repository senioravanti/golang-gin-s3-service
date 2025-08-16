FROM alpine:3.22

RUN apk add --no-cache --update \
  go

COPY . /usr/local/src/credgen/
WORKDIR /usr/local/src/credgen/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o /usr/local/bin/credgen ./cmd/main.go

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/credgen" ]