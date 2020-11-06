package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type book struct {
	title           string
	authors         []string
	averageRating   float32
	numberOfRatings int
}

var regexes struct {
	averageRatingRegex   *regexp.Regexp
	numberOfRatingsRegex *regexp.Regexp
}

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

func compileRegexes() {
	// average rating regex
	regexes.averageRatingRegex = regexp.MustCompile("([\\d|.]*) avg rating")
	// number of ratings regex
	regexes.numberOfRatingsRegex = regexp.MustCompile("([\\d|,]*) ratings")
}

func getDetails(e *colly.HTMLElement) (float64, int) {
	// get text
	text := e.ChildText(".uitext.greyText.smallText")

	// assign average rating
	averageRating := regexes.averageRatingRegex.FindStringSubmatch(text)[1]
	// convert to float
	averageRatingFloat, err := strconv.ParseFloat(averageRating, 32)
	if err != nil {
		log.Fatal("error while parsing averageRating to int: ", err)
	}

	// assign number of ratings
	numberOfRatings := regexes.numberOfRatingsRegex.FindStringSubmatch(text)[1]
	// remove commas from number of ratings
	numberOfRatings = strings.ReplaceAll(numberOfRatings, ",", "")
	// convert to integer
	numberOfRatingsInt, err := strconv.Atoi(numberOfRatings)
	if err != nil {
		log.Fatal("error while parsing numberOfRatings to int: ", err)
	}

	// return the final results
	return averageRatingFloat, numberOfRatingsInt
}

func getAuthors(e *colly.HTMLElement) []string {
	// create the authors array
	var authors []string
	// for each author
	e.ForEach(".authorName__container", func(i int, e *colly.HTMLElement) {
		// add the author to array
		authors = append(authors, strings.TrimSpace(e.Text))
	})
	// return the final results
	return authors
}

func handleBooks(c *colly.Collector, books []book, bookCount int) {
	// for each book result
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		// increment the counter for each book
		bookCount++
		// get the book name
		name := e.ChildText(".bookTitle")
		// log the founded book
		log.Println("Found book:", name)

		// get authors array
		authors := getAuthors(e)

		// parse and get details string
		averageRatingFloat, numberOfRatingsInt := getDetails(e)

		// generate and add struct to books array
		books = append(books, book{
			title:           name,
			authors:         authors,
			averageRating:   float32(averageRatingFloat),
			numberOfRatings: numberOfRatingsInt,
		})
	})
}

func startScraping(pageCountInt int, query string, c *colly.Collector) {
	// for every page - i starts at 1 and goes till including pageCountInt
	for i := 1; i <= pageCountInt; i++ {
		// generate url
		url := fmt.Sprintf("https://www.goodreads.com/search?page=%d&q=%s", i, query)
		// print url that is going to be parsed
		log.Println("Going to parse:", url)
		// start scraping
		err := c.Visit(url)
		if err != nil {
			log.Fatal("error while doing the request: ", err)
		}
	}
}

func main() {
	// compile regexes
	compileRegexes()

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

	// generate new collector - start the generator in async mode
	c := colly.NewCollector(colly.Async(true))

	// init books array
	var books []book
	// start the book counter
	bookCount := 0
	// register book handler
	handleBooks(c, books, bookCount)

	// start scraping
	startScraping(pageCountInt, query, c)

	// wait for all the threads to finish
	c.Wait()
	// print the final book count
	log.Println("Parsed book count:", bookCount)
}
