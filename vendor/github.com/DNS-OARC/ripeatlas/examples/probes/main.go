package main

import (
    "flag"
    "log"

    "github.com/DNS-OARC/ripeatlas"
    "github.com/DNS-OARC/ripeatlas/request"
)

var all bool
var page int
var file bool

func init() {
    flag.BoolVar(&all, "all", false, "list all probes")
    flag.IntVar(&page, "page", 1, "pagination page to load (default first, only used without -file and arguments)")
    flag.BoolVar(&file, "file", false, "arguments given are files to read (default probe ids to query for over HTTP)")
}

func print(p *request.Probe) {
    log.Printf("%d %s %s/%s %s", p.Id(), p.Type(), p.AddressV4(), p.AddressV6(), p.Description())
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
            probes, err := a.Probes(ripeatlas.Params{
                "page": int64(page),
            })
            if err != nil {
                log.Fatalf(err.Error())
            }

            got := 0
            for p := range probes {
                got++
                if p.ParseError != nil {
                    log.Println(p.ParseError.Error())
                    break
                }
                print(p)
            }
            if got == 0 {
                break
            }
            page++
        }
        return
    } else if !file && flag.NArg() <= 0 {
        probes, err := a.Probes(ripeatlas.Params{
            "page": int64(page),
        })
        if err != nil {
            log.Fatalf(err.Error())
        }

        for p := range probes {
            if p.ParseError != nil {
                log.Println(p.ParseError.Error())
                break
            }
            print(p)
        }
        return
    }

    for _, arg := range flag.Args() {
        var probes <-chan *request.Probe
        var err error

        if file {
            probes, err = a.Probes(ripeatlas.Params{
                "file": arg,
            })
            if err != nil {
                log.Fatalf(err.Error())
            }
        } else {
            probes, err = a.Probes(ripeatlas.Params{
                "pk": arg,
            })
            if err != nil {
                log.Fatalf(err.Error())
            }
        }

        for p := range probes {
            if p.ParseError != nil {
                log.Println(p.ParseError.Error())
                break
            }
            print(p)
        }
    }
}
