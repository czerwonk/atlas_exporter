FROM golang:1.26.6-alpine3.24@sha256:af8d6740070b8906d12eae1c3e3ea0957fb63f492051ea05e354c38ef9fe88df AS builder
ADD . /go/atlas_exporter/
WORKDIR /go/atlas_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/atlas_exporter


FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
RUN apk --no-cache add ca-certificates bash
WORKDIR /app
COPY --from=builder /go/bin/atlas_exporter .
EXPOSE 9400

ADD entrypoint /entrypoint
RUN chmod 0755 /entrypoint
ENTRYPOINT ["/entrypoint"]
