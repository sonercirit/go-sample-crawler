package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getInput(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + ": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("can't read input: ", err)
	}
	return strings.TrimSpace(text)
}

func main() {
	query := getInput("What should we search for? For example \"business\"")
	pageCount := getInput("How many pages should we crawl?")
	log.Println(query, pageCount)
}
