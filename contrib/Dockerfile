FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev libc-dev git &&\
    git clone --depth 1 https://github.com/mstcl/legit /build

WORKDIR /build

ENV CGO_ENABLED=0

RUN go build -o legit

FROM scratch

WORKDIR /config
COPY --chown=1000:1000 --from=builder /build/config.yaml /config

WORKDIR /app
COPY --chown=1000:1000 --from=builder /build/legit /app/legit

EXPOSE 8080

ENTRYPOINT ["/app/legit"]
CMD ["--config", "/config/config.yaml"]
