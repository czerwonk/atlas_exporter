BINS = measurements/measurements \
 probes/probes \
 reader/reader \
 streamer/streamer

all: $(BINS)

clean:
	rm -f $(BINS)

measurements/measurements:
	cd measurements && go build

probes/probes:
	cd probes && go build

reader/reader:
	cd reader && go build

streamer/streamer:
	cd streamer && go build
