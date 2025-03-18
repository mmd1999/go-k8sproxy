# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o k8sproxy .
RUN useradd -u 10001 nonroot

FROM scratch

COPY --from=builder /app/k8sproxy /app/k8sproxy
COPY --from=builder /etc/passwd /etc/passwd
USER nonroot
EXPOSE 5000
WORKDIR /app
ENTRYPOINT [ "/app/k8sproxy" ]