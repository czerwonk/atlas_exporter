package main

import (
    "flag"
    "log"
    "fmt"

    "github.com/DNS-OARC/ripeatlas"
)

var all bool
var page int
var file bool

func init() {
    flag.BoolVar(&all, "all", false, "list all measurements")
    flag.IntVar(&page, "page", 1, "pagination page to load (default first, only used without -file and arguments)")
    flag.BoolVar(&file, "file", false, "arguments given are files to read (default measurement ids to query for over HTTP)")
}

func print(m *ripeatlas.Measurement) {
    s := ""
    if m.IsOneoff() {
        s += " Online"
    } else {
        s += " Offline"
    }
    if m.IsPublic() {
        s += " Public"
    } else {
        s += " Private"
    }
    if m.TargetIp() != "" {
        s += " TargetIp="+m.TargetIp()
    }
    if m.TargetAsn() > 0 {
        s += " TargetAsn="+fmt.Sprintf("%d", m.TargetAsn())
    }
    log.Printf("%d %s%s %s", m.Id(), m.Type(), s, m.Description())
}

func main() {
    flag.Parse()

    var a ripeatlas.Atlaser
    if file {
        a = ripeatlas.NewFile()
    } else {
        a = ripeatlas.NewHttp()
    }

    if !file && all {
        page = 1
        for {
            measurements, err := a.Measurements(ripeatlas.Params{
                "page": int64(page),
            })
            if err != nil {
                log.Fatalf(err.Error())
            }

            got := 0
            for m := range measurements {
                got++
                if m.ParseError != nil {
                    log.Println(m.ParseError.Error())
                    break
                }
                print(m)
            }
            if got == 0 {
                break
            }
            page++
        }
        return
    } else if !file && flag.NArg() <= 0 {
        measurements, err := a.Measurements(ripeatlas.Params{
            "page": int64(page),
        })
        if err != nil {
            log.Fatalf(err.Error())
        }

        for m := range measurements {
            if m.ParseError != nil {
                log.Println(m.ParseError.Error())
                break
            }
            print(m)
        }
        return
    }

    for _, arg := range flag.Args() {
        var measurements <-chan *ripeatlas.Measurement
        var err error

        if file {
            measurements, err = a.Measurements(ripeatlas.Params{
                "file": arg,
            })
            if err != nil {
                log.Fatalf(err.Error())
            }
        } else {
            measurements, err = a.Measurements(ripeatlas.Params{
                "pk": arg,
            })
            if err != nil {
                log.Fatalf(err.Error())
            }
        }

        for m := range measurements {
            if m.ParseError != nil {
                log.Println(m.ParseError.Error())
                break
            }
            print(m)
        }
    }
}
