GO_SOURCES = ../measurement.go \
 ../measurement/measurement_status.go \
 ../measurement/probe.go \
 ../measurement/probe_source.go \
 ../request/probe.go \
 ../request/probe/number.go \
 ../request/probe/probe_tags.go \
 ../request/probe/probe_status.go

all: $(GO_SOURCES)

../measurement.go: gen.pl measurement.model
	./gen.pl measurement.model > ../measurement.go

../measurement/measurement_status.go: gen.pl measurement/measurement_status.model
	cd measurement && ../gen.pl measurement_status.model > ../../measurement/measurement_status.go

../measurement/probe.go: gen.pl measurement/probe.model
	cd measurement && ../gen.pl probe.model > ../../measurement/probe.go

../measurement/probe_source.go: gen.pl measurement/probe_source.model
	cd measurement && ../gen.pl probe_source.model > ../../measurement/probe_source.go

../request/probe.go: gen.pl request/probe.model
	cd request && ../gen.pl probe.model > ../../request/probe.go

../request/probe/number.go: gen.pl request/probe/number.model
	cd request/probe && ../../gen.pl number.model > ../../../request/probe/number.go

../request/probe/probe_tags.go: gen.pl request/probe/probe_tags.model
	cd request/probe && ../../gen.pl probe_tags.model > ../../../request/probe/probe_tags.go

../request/probe/probe_status.go: gen.pl request/probe/probe_status.model
	cd request/probe && ../../gen.pl probe_status.model > ../../../request/probe/probe_status.go
