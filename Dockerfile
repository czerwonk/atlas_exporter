FROM --platform=$BUILDPLATFORM golang:1.27.2-alpine3.24@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

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


FROM gcr.io/distroless/static-debian13:nonroot@sha256:e2e927ec666bae08560abb3c55d0659eceabb657f56b6782ab500a9fc7f555e3

WORKDIR /app
COPY --from=builder /out/atlas_exporter /app/atlas_exporter

EXPOSE 9400
ENTRYPOINT ["/app/atlas_exporter"]
