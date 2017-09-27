package main

import (
    "flag"
    "log"

    "github.com/DNS-OARC/ripeatlas"
)

var Msm = flag.Int("msm", 0, "stream specific measurement id (optional)")
var Type = flag.String("type", "", "stream specific measurement type")

func main() {
    flag.Parse()

    if Type == nil {
        log.Fatalf("Need -type")
    }

    p := ripeatlas.Params{}
    if Msm != nil {
        p["msm"] = *Msm
    }
    if Type != nil {
        p["type"] = *Type
    }

    a := ripeatlas.Atlaser(ripeatlas.NewStream())

    results, err := a.MeasurementResults(p)
    if err != nil {
        log.Fatalf(err.Error())
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
