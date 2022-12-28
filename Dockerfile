ARG GITHUB_PATH=github.com/Dsmit05/chat

FROM golang:1.19-alpine3.16 AS builder

RUN apk add --no-cache ca-certificates git make

WORKDIR /home/${GITHUB_PATH}

COPY . .

RUN make build

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /home/${GITHUB_PATH}/chat .

CMD ["./chat"]