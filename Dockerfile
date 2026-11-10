FROM --platform=$BUILDPLATFORM golang:1.27.4-alpine3.24@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
COPY atlas ./atlas
COPY config ./config
COPY dns ./dns
COPY exporter ./exporter
COPY http ./http
COPY ntp ./ntp
COPY ping ./ping
COPY probe ./probe
COPY sslcert ./sslcert
COPY traceroute ./traceroute

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
  GOARM=${TARGETVARIANT#v} \
  go build -trimpath \
  -ldflags="-s -w" \
  -o /out/atlas_exporter .


FROM gcr.io/distroless/static-debian13:nonroot@sha256:25ac8f0c3b2334a9eb05452266bd4697864a8c8ef0ef71e5e6480d7f92868329

WORKDIR /app
COPY --from=builder /out/atlas_exporter /app/atlas_exporter

EXPOSE 9400
ENTRYPOINT ["/app/atlas_exporter"]
