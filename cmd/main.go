// Markov chain generator for making up pretend languages

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/slabgorb/lango"
)

func readCorpus(path string) []string {
	re := regexp.MustCompile(`[^\p{L}]`)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	rawWords := strings.Split(string(dat), " ")
	var words []string
	for _, word := range rawWords {
		if word != "" {
			words = append(words, word)
		}
	}
	for index, word := range words {
		words[index] = re.ReplaceAllString(word, "")
	}
	return words
}

type arrayFlags struct {
	vals []string
}

func (a *arrayFlags) String() string {
	return strings.Join(a.vals, ":")
}

func (a *arrayFlags) Set(value string) error {
	a.vals = append(a.vals, value)
	return nil
}

var corpora arrayFlags
var lookback = 2

func init() {
	flag.Var(&corpora, "corpus", "path to corpus file, can be specified multiple times to mix corpora")
	flag.Parse()
}

func main() {
	wordSlices := make(chan []string)
	var words []string
	chain := lango.NewChain(2)
	go func() {
		wordSlices <- readCorpus("/Users/keithavery/Projects/fantasy-language-maker/corpus/french.txt")
	}()
	go func() {
		wordSlices <- readCorpus("/Users/keithavery/Projects/fantasy-language-maker/corpus/english.txt")
	}()
	for i := 0; i < 2; i++ {
		result := <-wordSlices
		words = append(words, result...)
	}
	for _, word := range words {
		chain.AddWord(word)
	}
	fmt.Println(chain.MakeWord())

}
