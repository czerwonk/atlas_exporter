# atlas_exporter [![Build Status](https://travis-ci.org/czerwonk/atlas_exporter.svg)][travis]
Metric exporter for RIPE atlas measurement results

## Remarks
* this is an early version, more features will be added step by step
* at the moment only the last result of an measurement is used
* the required Go version is 1.8+.

## Install
```
go get github.com/czerwonk/atlas_exporter
```

## Usage
### Start server
```
./atlas_exporter
```

### Call metrics URI
for measurement with id 7924888:
```
curl http://127.0.0.1:9400/metrics?measurement_id=7924888
```
the result should look similar to this one:
```
atlas_ping_success{measurement="7012288",probe="10563"} 1
atlas_ping_min_latency{measurement="7012288",probe="10563"} 153.079780
atlas_ping_max_latency{measurement="7012288",probe="10563"} 154.925535
atlas_ping_avg_latency{measurement="7012288",probe="10563"} 153.697307
atlas_ping_sent{measurement="7012288",probe="10563"} 3
atlas_ping_received{measurement="7012288",probe="10563"} 3
...
```

## Features
* ping measurements (success, min/max/avg latency)
* traceroute measurements (success, hop count)

## Third Party Components
This software uses components of the following projects
* Go bindings for RIPE Atlas API (https://github.com/DNS-OARC/ripeatlas)

## The RIPE Atlas Project
see http://atlas.ripe.net

[travis]: https://travis-ci.org/czerwonk/atlas_exporter
