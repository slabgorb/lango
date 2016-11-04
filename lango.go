// Markov chain generator for making up pretend languages

package main

import (
	"fmt"
	"strings"
	"math/rand"
	"time"
	"io/ioutil"
	"regexp"
)

var source = rand.NewSource(time.Now().UnixNano())
var random = rand.New(source)

type CharacterList struct {
	characters map[rune]int
}

type Chain struct {
	lookback int
	start rune
	end rune
	characterLists map[string]*CharacterList
}

func NewCharacterList() (*CharacterList) {
	var characterList CharacterList
	characterList.characters = make(map[rune]int)
	return &characterList
}

func NewChain(lookback int) (*Chain) {
	chain := Chain{lookback: lookback, start: '^', end: '$'}
	chain.characterLists = make(map[string]*CharacterList)
	return &chain
}

func (chain *Chain) AddRune(key string, char rune) {
	_,ok := chain.characterLists[key]
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
	chain.AddRune(key, chain.end)
	
}

func (chain *Chain) NewKey() ([]rune, string) {
	keyRunes := make([]rune, chain.lookback)
	for index,_ := range keyRunes {
		keyRunes[index] = chain.start
	}
	return keyRunes, string(keyRunes)
	
}

func rotateKey(key []rune, char rune) ([]rune, string) {
	key = append(key, char)
	key = key[1:]
	return key, string(key)
}

func (chain *Chain) MakeWord() (word string) {
	keyRunes, key :=  chain.NewKey()
	for {
		char := chain.characterLists[key].Choose(len(word))
		if char == chain.end { break }
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


func readCorpus(path string)  []string  {
	re := regexp.MustCompile(`[^\p{L}]`)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(dat), " ")
	for index, word := range words {
		words[index] = re.ReplaceAllString(word, "")
	}
	return words
}


func main() {
	wordSlices := make(chan []string)
	var words []string
	chain := NewChain(2)
	go func() { wordSlices <- readCorpus("/Users/keithavery/Projects/fantasy-language-maker/corpus/french.txt") }()
	go func() { wordSlices <- readCorpus("/Users/keithavery/Projects/fantasy-language-maker/corpus/english.txt") }()
	for i := 0; i < 2; i++ {
		result := <- wordSlices
		words = append(words, result...)
	}
	for _,word := range words {
		chain.AddWord(word)
	}
	fmt.Println(chain.MakeWord())
	
}
