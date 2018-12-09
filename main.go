package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// logger settings
	log.SetFlags(0)
	log.SetPrefix("")

	// set subcommand 'help' if it does not exist.
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, "help")
	}

	// switch processing by subcommand.
	subCmd := args[0]
	switch subCmd {
	case "server":
		fs := newServerFlagSet(flag.ExitOnError)
		fs.Parse(args[1:])
		serverMain()
	case "client":
		fs := newClientFlagSet(flag.ExitOnError)
		fs.Parse(args[1:])
		if err := clientMain(); err != nil {
			log.Fatalf("%v %v: %v", os.Args[0], subCmd, err)
		}
	case "help":
		subargs := args[1:]
		if len(subargs) == 0 {
			fmt.Printf("Usage of %v: %v [server|client|help]\n", os.Args[0], os.Args[0])
			return
		}
		subSubCmd := subargs[0]
		var fs *flag.FlagSet
		switch subSubCmd {
		case "server":
			fs = newServerFlagSet(flag.ExitOnError)
		case "client":
			fs = newClientFlagSet(flag.ExitOnError)
		case "help":
			fmt.Printf(`Usage of %v: %v %v
               %v %v [server|client|help]
`, subCmd, os.Args[0], subCmd, os.Args[0], subCmd)
			return
		default:
			log.Fatalf("%v %v: invalid subsubcommand: %v", os.Args[0], subCmd, subSubCmd)
		}
		fs.SetOutput(os.Stdout)
		fmt.Printf("Usage of %v:\n", subSubCmd)
		fs.PrintDefaults()
	default:
		log.Fatalf("%v: invalid subcommand: %v", os.Args[0], subCmd)
	}
}
