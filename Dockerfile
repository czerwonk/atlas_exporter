FROM golang as builder
RUN go get -d -v github.com/czerwonk/atlas_exporter
WORKDIR /go/src/github.com/czerwonk/atlas_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/czerwonk/atlas_exporter/app atlas_exporter
CMD ./atlas_exporter
EXPOSE 9400
