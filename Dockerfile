FROM golang:1.11 as builder

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ddns


FROM orvice/go-runtime:lite

COPY --from=builder /home/app/ddns .

ENTRYPOINT [ "./ddns" ]
