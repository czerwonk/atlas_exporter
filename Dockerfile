FROM golang as builder
RUN go get github.com/czerwonk/atlas_exporter

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/bin/atlas_exporter .
CMD /app/atlas_exporter
EXPOSE 9400
