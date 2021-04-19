// Randtext generates random text that reads well. It reads from stdin or from a list of named files.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/lorciv/randtext"
)

var (
	numWords = flag.Int("n", 100, "number of words to generate")
	//startWords = flag.String("s", "", "starting words for generated text")
)

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()

	rand.Seed(time.Now().Unix())

	if flag.NArg() == 0 {
		if err := randtext.Feed(os.Stdin); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, arg := range flag.Args() {
			f, err := os.Open(arg)
			if err != nil {
				log.Print(err)
				continue
			}

			if err := randtext.Feed(f); err != nil {
				log.Print(err)
				f.Close()
				continue
			}
			f.Close()
		}
	}

	if err := randtext.Generate(os.Stdout, *numWords); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	const descr = "Randtext generates random text that reads well. It reads from stdin or from a list of named files."
	fmt.Fprintln(flag.CommandLine.Output(), descr)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [file ...]\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
	flag.PrintDefaults()
}
