FROM golang:1.10 AS builder

WORKDIR /go/src/simple-dns

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM alpine
COPY --from=builder /go/src/simple-dns/simple-dns /usr/local/bin/

EXPOSE 53

CMD ["simple-dns"]
