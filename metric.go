package main

import "io"

type Metric interface {
	Write(w io.Writer, pk string)

	SetAsn(asn int)
}
