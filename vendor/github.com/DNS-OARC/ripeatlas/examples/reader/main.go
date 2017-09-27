package main

import (
    "flag"
    "log"
    "time"

    "github.com/DNS-OARC/ripeatlas"
    "github.com/DNS-OARC/ripeatlas/measurement"
)

var start int
var stop int
var last int
var file bool
var frag bool

func init() {
    flag.IntVar(&start, "start", 0, "start unixtime for results")
    flag.IntVar(&stop, "stop", 0, "stop unixtime for results")
    flag.IntVar(&last, "last", 0, "last N seconds of results, not used if start/stop are used")
    flag.BoolVar(&file, "file", false, "arguments given are files to read (default measurement ids to query for over HTTP)")
    flag.BoolVar(&frag, "frag", false, "if true, use/input is fragmented JSON")
}

func main() {
    flag.Parse()

    var startTime, stopTime time.Time
    var latest bool

    if last > 0 {
        stopTime = time.Now()
        startTime = stopTime.Add(time.Duration(-last) * time.Second)
    } else if start > 0 && stop > 0 {
        startTime = time.Unix(int64(start), 0)
        stopTime = time.Unix(int64(stop), 0)
    } else {
        latest = true
    }

    var msm ripeatlas.Atlaser
    if file {
        msm = ripeatlas.NewFile()
    } else {
        msm = ripeatlas.NewHttp()
    }

    for _, arg := range flag.Args() {
        var results <-chan *measurement.Result
        var err error

        if file {
            results, err = msm.MeasurementResults(ripeatlas.Params{
                "file":       arg,
                "fragmented": frag,
            })
            if err != nil {
                log.Fatalf(err.Error())
            }
        } else {
            if latest {
                results, err = msm.MeasurementLatest(ripeatlas.Params{
                    "pk":         arg,
                    "fragmented": frag,
                })
            } else {
                results, err = msm.MeasurementResults(ripeatlas.Params{
                    "start":      startTime.Unix(),
                    "stop":       stopTime.Unix(),
                    "pk":         arg,
                    "fragmented": frag,
                })
            }

            if err != nil {
                log.Fatalf(err.Error())
            }
        }

        for r := range results {
            if r.ParseError != nil {
                log.Println(r.ParseError.Error())
                break
            }
            log.Printf("%d %s", r.MsmId(), r.Type())

            switch r.Type() {
            case "dns":
                if r.DnsResult() != nil {
                    m, _ := r.DnsResult().UnpackAbuf()
                    if m != nil {
                        log.Printf("%v", m)
                    }
                }
                for _, s := range r.DnsResultsets() {
                    if s.Result() != nil {
                        m, _ := s.Result().UnpackAbuf()
                        if m != nil {
                            log.Printf("%v", m)
                        }
                    }
                }
            case "ping":
                for _, p := range r.PingResults() {
                    log.Printf("%v", p.Rtt())
                }
            case "traceroute":
                for _, t := range r.TracerouteResults() {
                    log.Printf("%v", t.Hop())
                }
            case "http":
                log.Printf("%s", r.Uri())
                for _, h := range r.HttpResults() {
                    log.Printf("header %d body %d", h.Hsize(), h.Bsize())
                }
            case "ntp":
                for _, n := range r.NtpResults() {
                    log.Printf("%f", n.ReceiveTs())
                }
            case "sslcert":
                if len(r.Cert()) > 0 {
                    log.Printf("%s", r.Cert()[0])
                }
            case "wifi":
                for k, v := range r.WpaSupplicant() {
                    log.Printf("%s: %v", k, v)
                }
            }
        }
    }
}
