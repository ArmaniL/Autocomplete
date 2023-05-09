package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"unicode"
)

type CharNode struct {
	letter   rune
	children [26]*CharNode
}

type AutoComplete struct {
	prefix   string
	limit    int
	amount   int
	children [26]*CharNode
}

func (ac *AutoComplete) isInvalidCharacter(c rune) bool {
	return !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z')
}

func (ac *AutoComplete) AddWords(wordBank []string) {
	for _, word := range wordBank {
		ac.AddWord(word)
	}
}

func (ac *AutoComplete) AddWord(word string) {

	var index int
	children := &ac.children
	for _, ch := range word {
		if ac.isInvalidCharacter(ch) {

			continue
		}
		ch = unicode.ToLower(ch)
		index = charHash(ch)

		if children[index] == nil {
			children[index] = new(CharNode)
			children[index].letter = ch
			children[index].children = [26]*CharNode{}
		}
		children = &children[index].children
	}

}

func (ac *AutoComplete) GuessWord(prefix string) ([]string, error) {
	var listOfWords []string
	var sb string
	ac.prefix = prefix
	children := &ac.children
	ac.limit = -1

	for _, ch := range prefix {
		ch = unicode.ToLower(ch)
		index := charHash(ch)
		if children[index] == nil {
			return nil, errors.New("no answer")
		}
		children = &children[index].children
	}

	for _, child := range children {

		ac.depthFirstAccumulator(child, &sb, &listOfWords)
	}
	return listOfWords, nil
}

func (ac *AutoComplete) GuessNWords(prefix string, n int) ([]string, error) {

	var listOfWords []string
	var sb string
	ac.prefix = prefix
	children := &ac.children
	ac.limit = n
	ac.amount = 0

	for _, ch := range prefix {
		ch = unicode.ToLower(ch)
		index := charHash(ch)
		if children[index] == nil {
			return nil, errors.New("no answer")
		}
		children = &children[index].children
	}

	for _, child := range children {

		ac.depthFirstAccumulator(child, &sb, &listOfWords)
	}
	return listOfWords, nil
}

func (ac *AutoComplete) depthFirstAccumulator(node *CharNode, sb *string, words *[]string) {
	if ac.limit != -1 && ac.amount >= ac.limit {
		return
	}

	if node == nil {
		return
	}

	*sb += string(node.letter)
	lastLetter := true
	for _, child := range node.children {
		lastLetter = lastLetter && (child == nil)
	}
	if lastLetter {
		*words = append(*words, ac.prefix+*sb)
		ac.amount = ac.amount + 1
	} else {
		for _, child := range node.children {
			if child != nil {

				ac.depthFirstAccumulator(child, sb, words)
			}
		}
	}
	str := *sb
	newStr := str[:len(str)-1]
	*sb = newStr

}

func charHash(c rune) int {

	return int(c - 'a')
}

func (ac *AutoComplete) AddWordsFromFile(path string) error {
	uniqueWords := make(map[string]bool)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		for _, word := range words {
			if word != "" {
				uniqueWords[word] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	for word := range uniqueWords {
		ac.AddWord(word)
	}
	return nil
}

func NewAutoComplete() AutoComplete {
	return AutoComplete{
		children: [26]*CharNode{},
	}
}

// func main() {
// 	trial := NewAutoComplete()
// 	err := trial.AddWordsFromFile("/Users/armaniweise/Desktop/words.csv")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for {
// 		fmt.Println("Enter a prefix")
// 		scanner.Scan()
// 		prefix := scanner.Text()
// 		words, err := trial.GuessNWords(prefix, 10)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		fmt.Println(words)
// 	}

// }
