FROM golang

RUN apt-get install -y git && \
    go get github.com/czerwonk/atlas_exporter

CMD atlas_exporter 
EXPOSE 9400
