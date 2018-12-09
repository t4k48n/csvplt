package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var sPort string

func newServerFlagSet(flags flag.ErrorHandling) *flag.FlagSet {
	fs := flag.NewFlagSet("server", flags)
	fs.StringVar(&sPort, "port", "", "communication port (e.g. ':8080')")
	return fs
}

type State struct {
	*Packet
	sync.RWMutex
}

var state State

func init() {
	state = State{
		Packet: &Packet{
			Data: &Data{
				Datasets: make([]Dataset, 0, 8),
			},
			Options: &Options{
				ShowLines: true,
				Elements: Elements{
					Line: Line{
						Fill:    false,
						Tension: 0,
					},
					Point: Point{
						Radius: 3,
					},
				},
			},
		},
	}
}

func serverMain() {
	if sPort == "" {
		fs := newServerFlagSet(flag.ExitOnError)
		fmt.Fprintf(os.Stderr, "flag not specified: -port\nUsage of %v:\n", "server")
		fs.PrintDefaults()
		os.Exit(2)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, IndexHead + sPort + IndexTail)
	})
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			state.RLock()
			packBytes, err := json.Marshal(state.Packet)
			state.RUnlock()
			if err != nil {
				log.Fatal(err)
			}
			w.Header().Set("Content-Type", `application/json;charset="utf-8"`)
			fmt.Fprintf(w, "%s", packBytes)
		case http.MethodPost:
			state.Lock()
			packBytes, err := ioutil.ReadAll(r.Body)
			state.Unlock()
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(packBytes, &state)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	log.Printf("Listening and serving (port %v) ...\n", sPort)
	log.Fatal(http.ListenAndServe(sPort, nil))
}
