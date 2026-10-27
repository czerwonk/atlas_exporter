FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4 AS builder

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


FROM gcr.io/distroless/static-debian13:nonroot@sha256:18f4cc26c9ffa7f20a359298be7deadac5c607c8e03a02b8119f1a5fc2f82c99

WORKDIR /app
COPY --from=builder /out/atlas_exporter /app/atlas_exporter

EXPOSE 9400
ENTRYPOINT ["/app/atlas_exporter"]
