package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var cPort string
var xc int
var line, point bool
var cfs *flag.FlagSet

var NoData = fmt.Errorf("NoData")

func newClientFlagSet(flags flag.ErrorHandling) *flag.FlagSet {
	fs := flag.NewFlagSet("client", flags)
	fs.StringVar(&cPort, "port", "", "communication port (e.g. ':8080')")
	fs.IntVar(&xc, "xc", -1, "index of x-column")
	fs.BoolVar(&line, "l", false, "")
	fs.BoolVar(&line, "line", false, "show line")
	fs.BoolVar(&point, "p", false, "")
	fs.BoolVar(&point, "point", false, "show point")
	cfs = fs
	return cfs
}

func clientMain() error {
	if cPort == "" {
		fs := newClientFlagSet(flag.ExitOnError)
		fmt.Fprintf(os.Stderr, "flag not specified: -port\nUsage of %v:\n", "client")
		fs.PrintDefaults()
		os.Exit(2)
	}
	var reader io.Reader
	switch n := len(cfs.Args()); n {
	case 0:
		reader = os.Stdin
	case 1:
		var err error
		p, err := filepath.Abs(cfs.Arg(0))
		if err != nil {
			return err
		}
		reader, err = os.Open(p)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("zero or one file acceptable but %d files given", n)
	}

	table, err := parseTable(reader)
	if err != nil {
		return fmt.Errorf("parseTable: %v", err)
	}
	p := Packet{
		Data:    NewDataFromTable(table, xc),
		Options: new(Options),
	}
	if !line && !point {
		line = true
		point = true
	}
	p.Options.ShowLines = line
	if point {
		p.Options.Elements.Point.Radius = 3
	} else {
		p.Options.Elements.Point.Radius = 0
	}

	packBytes, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}
	jsonReader := bytes.NewReader(packBytes)
	_, err = http.Post(fmt.Sprintf("http://localhost%v/request", cPort), `application/json;charset="utf-8"`, jsonReader)
	return err
}

func parseTable(r io.Reader) ([][]float64, error) {
	strTab, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	nr := len(strTab)
	if nr == 0 {
		return nil, NoData
	}
	nc := len(strTab[0])
	if nc == 0 {
		return nil, NoData
	}
	table := make([][]float64, 0, nr)
	for i := 0; i < nr; i++ {
		table = append(table, make([]float64, nc))
		for j := range table[i] {
			var err error
			table[i][j], err = strconv.ParseFloat(strTab[i][j], 64)
			if err != nil {
				return nil, err
			}
		}
	}
	return table, nil
}
