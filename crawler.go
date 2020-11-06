package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getInput(question string, def string) string {
	// read from stdin
	reader := bufio.NewReader(os.Stdin)
	// ask the question
	fmt.Printf("%s Default %s: ", question, def)
	// read until the newline separator
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("can't read input: ", err)
	}
	// trim newline separator
	text = strings.TrimSpace(text)
	// if input is empty
	if text == "" {
		// assign the default value
		text = def
	}
	// return the final result
	return text
}

func main() {
	// ask about the query
	query := getInput("What should we search for?", "fantasy")
	// ask about the page count
	pageCount := getInput("How many pages should we crawl?", "10")
	log.Println(query, pageCount)
}
