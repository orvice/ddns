FROM golang:1.22 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build


FROM ghcr.io/orvice/go-runtime:master

LABEL org.opencontainers.image.description "DDNS"

ENV PROJECT_NAME ddns

COPY --from=builder /home/app/bin/${PROJECT_NAME} .
