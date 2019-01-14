FROM golang as builder
RUN go get -d -v github.com/czerwonk/atlas_exporter
WORKDIR /go/src/github.com/czerwonk/atlas_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /app
COPY --from=builder /go/src/github.com/czerwonk/atlas_exporter/app atlas_exporter
EXPOSE 9400

ADD entrypoint /entrypoint
RUN chmod 0755 /entrypoint
ENTRYPOINT ["/entrypoint"]
