ARG ALPINE_VERSION=3.21.3
ARG GOLANG_VERSION=1.24

FROM golang:${GOLANG_VERSION} AS builder

WORKDIR /builder

COPY go.mod go.sum ./
RUN go mod download

COPY *.go .
RUN CGO_ENABLED=0 go build -o madvsa *.go

FROM alpine:${ALPINE_VERSION}

WORKDIR /app
ENV PATH="$PATH:/app/bin"

COPY --from=builder /builder/madvsa .

ENTRYPOINT ["./madvsa"]
