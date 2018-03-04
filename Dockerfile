FROM golang:1.9 as builder

## Create a directory and Add Code
RUN mkdir -p /go/src/github.com/orvice/ddns
WORKDIR /go/src/github.com/orvice/ddns
ADD .  /go/src/github.com/orvice/ddns

# Download and install any required third party dependencies into the container.
RUN go-wrapper download
# RUN go-wrapper install
RUN CGO_ENABLED=0 go build


FROM alpine

COPY --from=builder /go/src/github.com/orvice/ddns/ddns .

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
# Change TimeZone
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
# Clean APK cache
RUN rm -rf /var/cache/apk/*

ENTRYPOINT [ "./ddns" ]