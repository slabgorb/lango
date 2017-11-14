package lango

import (
	"math/rand"
	"strings"
	"time"
)

type CharacterList struct {
	characters map[rune]int
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	start = '^'
	end   = '$'
)

type Chain struct {
	lookback       int
	characterLists map[string]*CharacterList
}

func newCharacterList() *CharacterList {
	return &CharacterList{characters: make(map[rune]int)}
}

func NewChain(lookback int) *Chain {
	return &Chain{lookback: lookback, characterLists: make(map[string]*CharacterList)}
}

func (chain *Chain) addRune(key string, char rune) {
	_, ok := chain.characterLists[key]
	if !ok {
		chain.characterLists[key] = newCharacterList()
	}
	chain.characterLists[key].characters[char]++
}

func (chain *Chain) AddWord(word string) {
	keyRunes, key := chain.newKey()
	for _, char := range strings.ToLower(word) {
		chain.addRune(key, char)
		keyRunes, key = rotateKey(keyRunes, char)
	}
	chain.addRune(key, end)

}

func (chain *Chain) newKey() ([]rune, string) {
	keyRunes := make([]rune, chain.lookback)
	for index, _ := range keyRunes {
		keyRunes[index] = start
	}
	return keyRunes, string(keyRunes)

}

func rotateKey(key []rune, char rune) ([]rune, string) {
	key = append(key, char)[1:]
	return key, string(key)
}

func (chain *Chain) MakeWord() (word string) {
	keyRunes, key := chain.newKey()
	for {
		char := chain.characterLists[key].Choose(len(word))
		if char == end {
			break
		}
		word += string(char)
		keyRunes, key = rotateKey(keyRunes, char)

	}
	return
}

func (characterList *CharacterList) totalCounts() (total int) {
	total = 0
	for _, value := range characterList.characters {
		total = total + value
	}
	return
}

func (characterList *CharacterList) Choose(wordLength int) (choice rune) {
	randomNumber := random.Intn(characterList.totalCounts())
	position := 0
	for char, count := range characterList.characters {
		position += count
		if randomNumber <= position {
			choice = char
		}

	}
	return
}
