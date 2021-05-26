FROM golang:1.16 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build


FROM orvice/go-runtime

ENV PROJECT_NAME ddns

COPY --from=builder /home/app/bin/${PROJECT_NAME} .

ENTRYPOINT "./ddns"
