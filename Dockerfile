FROM golang:1.11 as builder

## Create a directory and Add Code
RUN mkdir -p /go/src/github.com/orvice/ddns
WORKDIR /go/src/github.com/orvice/ddns
ADD .  /go/src/github.com/orvice/ddns

RUN CGO_ENABLED=0 go build


FROM orvice/go-runtime:lite

COPY --from=builder /go/src/github.com/orvice/ddns/ddns .


ENTRYPOINT [ "./ddns" ]