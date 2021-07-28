FROM golang:alpine AS builder

WORKDIR $GOPATH/src/target
COPY . $GOPATH/src/target

RUN apk update && apk add --no-cache git ca-certificates
RUN adduser -D -g '' gouser

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o /go/bin ./...

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/ /go/bin/

USER gouser

CMD ["/go/bin/arch.homework"]
