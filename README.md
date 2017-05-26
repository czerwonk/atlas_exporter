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
atlas_traceroute_hops{measurement="7924888",probe="15072",asn="20375"} 13
atlas_traceroute_success{measurement="7924888",probe="15072",asn="20375"} 1
atlas_traceroute_hops{measurement="7924888",probe="15093",asn="3265"} 8
atlas_traceroute_success{measurement="7924888",probe="15093",asn="3265"} 1
...
```

## Features
* ping measurements (success, min/max/avg latency, dups, size)
* traceroute measurements (success, hop count, rtt)
* ntp
* dns (succress, rtt)

## Configuration (Prometheus)
```
  - job_name: 'atlas_exporter'
    scrape_interval: 5m
    static_configs:
      - targets:
        - 7924888
        - 7924886
    relabel_configs:
      - source_labels: [__address__]
        regex: (.*)(:80)?
        target_label: __param_measurement_id
        replacement: ${1}
      - source_labels: [__param_measurement_id]
        regex: (.*)
        target_label: instance
        replacement: ${1}
      - source_labels: []
        regex: .*
        target_label: __address__
        replacement: atlas-exporter.mytld:9400

```

## Third Party Components
This software uses components of the following projects
* Go bindings for RIPE Atlas API (https://github.com/DNS-OARC/ripeatlas)

## The RIPE Atlas Project
see http://atlas.ripe.net

[travis]: https://travis-ci.org/czerwonk/atlas_exporter
