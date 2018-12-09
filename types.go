package main

import "strconv"

type Packet struct {
	Data    *Data    `json:"data"`
	Options *Options `json:"options"`
}

type Data struct {
	Datasets []Dataset `json:"datasets"`
}

func NewDataFromTable(table [][]float64, xcol int) *Data {
	nr := len(table)
	nc := len(table[0])

	xs := make([]float64, nr)
	if xcol < 0 || xcol >= nc {
		for i := 0; i < nr; i++ {
			xs[i] = float64(i)
		}
	} else {
		for i := 0; i < nr; i++ {
			xs[i] = table[i][xcol]
		}
	}

	d := new(Data)
	d.Datasets = make([]Dataset, 0, nc)
	cIdx := 0
	cLen := len(Colors)
	for j := 0; j < nc; j++ {
		if j == xcol {
			continue
		}
		xys := make([]XY, nr)
		for i := 0; i < nr; i++ {
			xys[i].X = xs[i]
			xys[i].Y = table[i][j]
		}
		ds := Dataset{
			Label:           strconv.Itoa(cIdx),
			BackgroundColor: Colors[cIdx%cLen],
			BorderColor:     Colors[cIdx%cLen],
			Data:            xys,
		}
		d.Datasets = append(d.Datasets, ds)
		cIdx++
	}
	return d
}

type Dataset struct {
	Label           string `json:"label"`
	BackgroundColor string `json:"backgroundColor"`
	BorderColor     string `json:"borderColor"`
	Data            []XY   `json:"data"`
}

const (
	Blue    = "#0000ff"
	Red     = "#ff0000"
	Green   = "#008000"
	Cyan    = "#00bfbf"
	Magenta = "#bf00bf"
	Yello   = "#bfbf00"
	Black   = "#000000"
)

var Colors = []string{
	Blue,
	Red,
	Green,
	Cyan,
	Magenta,
	Yello,
	Black,
}

type XY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Options struct {
	ShowLines bool     `json:"showLines"`
	Elements  Elements `json:"elements"`
}

type Elements struct {
	Line  Line  `json:"line"`
	Point Point `json:"point"`
}

type Line struct {
	Fill    bool    `json:"fill"`
	Tension float64 `json:"tension"`
}

type Point struct {
	Radius float64 `json:"radius"`
}
