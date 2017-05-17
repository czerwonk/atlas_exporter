# atlas_exporter [![Build Status](https://travis-ci.org/czerwonk/atlas_exporter.svg)][travis]
Metric exporter for RIPE atlas measurement results

## Remarks
this is an early version, more features will be added step by step

## Install
```
go get github.com/czerwonk/atlas_exporter
```

## Features
* ping measurements (success, min/max/avg latency)
* traceroute measurements (success, hop count)

## Third Party Components
This software uses components of the following projects
* Go bindings for RIPE Atlas API (https://github.com/DNS-OARC/ripeatlas)

## The RIPE Atlas Project
see htt://atlas.ripe.net

[travis]: https://travis-ci.org/czerwonk/atlas_exporter