package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getInput(question string, def string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s Default %s: ", question, def)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("can't read input: ", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		text = def
	}
	return text
}

func main() {
	query := getInput("What should we search for?", "fantasy")
	pageCount := getInput("How many pages should we crawl?", "10")
	log.Println(query, pageCount)
}
