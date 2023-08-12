// Package main is simply a container for the exporter
//
// You can find all relevant code in [github.com/jjack/herpstat_spyderweb_exporter/exporter]
package main

import (
	"github.com/jjack/herpstat_spyderweb_exporter/exporter"
)

func main() {
	exporter.Run()
}
