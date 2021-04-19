// Package randtext generates random text that reads well.
package randtext

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type prefix []string

func (p prefix) shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

func (p prefix) string() string {
	return strings.Join(p, " ")
}

// Rand is a source of random text.
type Rand struct {
	stateTab  map[string][]string
	prefixLen int
}

// New returns a new Rand with given prefix length.
func New(prefixLen int) *Rand {
	return &Rand{
		stateTab:  make(map[string][]string),
		prefixLen: prefixLen,
	}
}

// Feed feeds the random source with text coming from the specified Reader.
func (r *Rand) Feed(in io.Reader) error {
	scan := bufio.NewScanner(in)
	scan.Split(bufio.ScanWords)
	pref := make(prefix, r.prefixLen)
	for scan.Scan() {
		word := scan.Text()
		key := pref.string()
		r.stateTab[key] = append(r.stateTab[key], word)
		pref.shift(word)
	}
	if err := scan.Err(); err != nil {
		return fmt.Errorf("could not read text to feed: %v", err)
	}
	return nil
}

// Generate generates random text to the specified Writer.
//
// It makes use of the input previously fed to the source to generate text that reads well. Output text
// contains at most the given number of words.
func (r *Rand) Generate(out io.Writer, words int) error {
	if len(r.stateTab) == 0 {
		return errors.New("could not generate: no text has been fed")
	}

	pref := make(prefix, r.prefixLen)
	for i := 0; i < words; i++ {
		choices := r.stateTab[pref.string()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		if _, err := fmt.Fprintf(out, " %s", next); err != nil {
			return fmt.Errorf("could not generate: %v", err)
		}
		pref.shift(next)
	}

	return nil
}

/*
 * Top-level convenience functions
 */

var globalRand = New(2)

// Feed feeds the default random source with text coming from the specified Reader.
func Feed(in io.Reader) error {
	return globalRand.Feed(in)
}

// Generate generates random text from the default source to the specified Writer.
//
// It makes use of the input previously fed to the default source to generate text that reads well.
// Output text contains at most the given number of words.
func Generate(out io.Writer, words int) error {
	return globalRand.Generate(out, words)
}

// func print() {
// 	const format = "%v\t%v\t\n"
// 	tw := new(tabwriter.Writer).Init(os.Stderr, 0, 8, 2, '\t', 0)
// 	fmt.Fprintf(tw, format, "PREFIX", "SUFFIX")
// 	fmt.Fprintf(tw, format, "------", "------")
// 	for p, s := range globalRand.stateTab {
// 		fmt.Fprintf(tw, format, p, s)
// 	}
// 	tw.Flush()
// }
