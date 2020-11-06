package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
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
	// parse pageCount to int
	pageCountInt, err := strconv.Atoi(pageCount)
	if err != nil {
		log.Fatal("can't parse page count to int: ", err)
	}
	// print the results for user
	log.Println("Detected inputs:", query, pageCountInt)

	// generate new collector
	c := colly.NewCollector()

	// for every page - i starts at 1 and goes till including pageCountInt
	for i := 1; i <= pageCountInt; i++ {
		// generate url
		url := fmt.Sprintf("https://www.goodreads.com/search?page=%d&q=%s", i, query)
		// print url that is going to be parsed
		log.Println("Going to parse", url)
		// start scraping
		err = c.Visit(url)
		if err != nil {
			log.Fatal("error while doing the request: ", err)
		}
	}
}
