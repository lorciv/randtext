// Package randtext generates random text that reads well.
package randtext

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
)

var suffixTable = make(map[string][]string)

// Feed reads the text coming from the specified Reader and feeds it to the random text engine.
func Feed(in io.Reader) error {
	var words []string
	scan := bufio.NewScanner(in)
	scan.Split(bufio.ScanWords)
	for scan.Scan() {
		words = append(words, scan.Text())
	}
	if err := scan.Err(); err != nil {
		return err
	}

	// prefix of fixed length 2 for simplicity
	for i := 0; i < len(words)-2; i++ {
		p := strings.Join(words[i:i+2], " ")
		suffixTable[p] = append(suffixTable[p], words[i+2])
	}

	return nil
}

// Emit emits random text to the specified Writer. It makes use of the input previously fed to the engine
// to generate text that reads well. Output text contains at most the given number of words.
func Emit(out io.Writer, words int) error {
	if len(suffixTable) == 0 {
		return errors.New("cannot emit: some text must be fed first")
	}

	var pref []string
	for p := range suffixTable {
		pref = strings.Fields(p) // randomly select start prefix to improve randomness
	}
	if _, err := fmt.Fprintf(out, "%s %s", strings.Title(pref[0]), pref[1]); err != nil {
		return err
	}

	for i := 0; i < words; i++ {
		s := suffixTable[strings.Join(pref, " ")]
		if len(s) == 0 {
			break
		}
		word := s[rand.Intn(len(s))]

		if _, err := fmt.Fprintf(out, " %s", word); err != nil {
			return err
		}

		pref[0], pref[1] = pref[1], word
	}

	if _, err := fmt.Fprint(out, "\n"); err != nil {
		return err
	}

	return nil
}

func print() {
	const format = "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stderr, 0, 8, 2, '\t', 0)
	fmt.Fprintf(tw, format, "PREFIX", "SUFFIX")
	fmt.Fprintf(tw, format, "------", "------")
	for p, s := range suffixTable {
		fmt.Fprintf(tw, format, p, s)
	}
	tw.Flush()
}
