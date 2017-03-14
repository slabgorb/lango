// Markov chain generator for making up pretend languages

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type CharacterList struct {
	characters map[rune]int
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	START = '^'
	END   = '$'
)

type Chain struct {
	lookback       int
	characterLists map[string]*CharacterList
}

func NewCharacterList() *CharacterList {
	return &CharacterList{characters: make(map[rune]int)}
}

func NewChain(lookback int) *Chain {
	return &Chain{lookback: lookback, characterLists: make(map[string]*CharacterList)}
}

func (chain *Chain) AddRune(key string, char rune) {
	_, ok := chain.characterLists[key]
	if !ok {
		chain.characterLists[key] = NewCharacterList()
	}
	chain.characterLists[key].characters[char]++
}

func (chain *Chain) AddWord(word string) {
	keyRunes, key := chain.NewKey()
	for _, char := range strings.ToLower(word) {
		chain.AddRune(key, char)
		keyRunes, key = rotateKey(keyRunes, char)
	}
	chain.AddRune(key, END)

}

func (chain *Chain) NewKey() ([]rune, string) {
	keyRunes := make([]rune, chain.lookback)
	for index, _ := range keyRunes {
		keyRunes[index] = START
	}
	return keyRunes, string(keyRunes)

}

func rotateKey(key []rune, char rune) ([]rune, string) {
	key = append(key, char)[1:]
	return key, string(key)
}

func (chain *Chain) MakeWord() (word string) {
	keyRunes, key := chain.NewKey()
	for {
		char := chain.characterLists[key].Choose(len(word))
		if char == END {
			break
		}
		word += string(char)
		keyRunes, key = rotateKey(keyRunes, char)

	}
	return
}

func (characterList *CharacterList) TotalCounts() (total int) {
	total = 0
	for _, value := range characterList.characters {
		total = total + value
	}
	return
}

func (characterList *CharacterList) Choose(wordLength int) (choice rune) {
	randomNumber := random.Intn(characterList.TotalCounts())
	position := 0
	for char, count := range characterList.characters {
		position += count
		if randomNumber <= position {
			choice = char
		}

	}
	return
}

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

func main() {
	wordSlices := make(chan []string)
	var words []string
	chain := NewChain(2)
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
