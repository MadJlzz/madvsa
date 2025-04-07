ARG TRIVY_VERSION=0.60.0
ARG ALPINE_VERSION=3.21.3
ARG GOLANG_VERSION=1.24

FROM golang:${GOLANG_VERSION} AS builder

WORKDIR /builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o trivyvsa cmd/trivy/main.go

FROM aquasec/trivy:${TRIVY_VERSION} AS trivy

FROM alpine:${ALPINE_VERSION}

WORKDIR /app
ENV PATH="$PATH:/app/bin"

COPY --from=trivy /usr/local/bin/trivy bin/trivy
COPY --from=builder /builder/trivyvsa .

ENTRYPOINT ["./trivyvsa"]
CMD ["--help"]