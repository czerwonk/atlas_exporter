SOURCES := $(wildcard *.go) \
    $(wildcard measurement/*.go) \
    $(wildcard measurement/dns/*.go) \
    $(wildcard measurement/ping/*.go) \
    $(wildcard measurement/traceroute/*.go) \
    $(wildcard measurement/http/*.go) \
    $(wildcard measurement/ntp/*.go) \
    $(wildcard measurement/sslcert/*.go)

EXAMPLE_SOURCES = $(wildcard example/*/*.go)

all: .built

fmt: format

clean:
	rm -f .built

format:
	gofmt -w $(SOURCES) $(EXAMPLE_SOURCES)
	sed -i -e 's%	%    %g' $(SOURCES) $(EXAMPLE_SOURCES)

.built: $(SOURCES)
	go build -v -x
	touch .built
