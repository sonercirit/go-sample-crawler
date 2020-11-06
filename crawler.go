package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type book struct {
	Title           string   `json:"title,omitempty"`
	Authors         []string `json:"authors,omitempty"`
	AverageRating   float32  `json:"average_rating,omitempty"`
	NumberOfRatings int      `json:"number_of_ratings,omitempty"`
	Published       int      `json:"published,omitempty"`
	Editions        int      `json:"editions,omitempty"`
}

var regexes struct {
	averageRatingRegex   *regexp.Regexp
	numberOfRatingsRegex *regexp.Regexp
	publishedRegex       *regexp.Regexp
	editionsRegex        *regexp.Regexp
}

func getInput(question string, def string) string {
	// read from stdin
	reader := bufio.NewReader(os.Stdin)
	// ask the question
	fmt.Printf("%s Default \"%s\": ", question, def)
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
	return text
}

func compileRegexes() {
	// average rating regex
	regexes.averageRatingRegex = regexp.MustCompile("([\\d|.]+) avg rating")
	// number of ratings regex
	regexes.numberOfRatingsRegex = regexp.MustCompile(" ([\\d|,]+) rating")
	// published regex - extra whitespace check because some date are separated by newline (\n)
	regexes.publishedRegex = regexp.MustCompile("published\\s+(\\d+)")
	// editions regex
	regexes.editionsRegex = regexp.MustCompile("(\\d+) edition")
}

func getAverageRating(text string) (float64, error) {
	// assign average rating
	averageRating := regexes.averageRatingRegex.FindStringSubmatch(text)[1]
	// convert to float
	averageRatingFloat, err := strconv.ParseFloat(averageRating, 32)
	if err != nil {
		return 0, err
	}
	return averageRatingFloat, nil
}

func getNumberOfRatings(text string) (int, error) {
	// assign number of ratings
	numberOfRatings := regexes.numberOfRatingsRegex.FindStringSubmatch(text)[1]
	// remove commas from number of ratings
	numberOfRatings = strings.ReplaceAll(numberOfRatings, ",", "")
	// convert to integer
	numberOfRatingsInt, err := strconv.Atoi(numberOfRatings)
	if err != nil {
		return 0, err
	}
	return numberOfRatingsInt, err
}

func getPublished(text string) (int, error) {
	// get matches
	publishedMatches := regexes.publishedRegex.FindStringSubmatch(text)
	// init publishedInt variable
	var publishedInt int
	// also init error to avoid compile error
	var err error
	// there might not be a publish date
	if publishedMatches != nil {
		// assign publish year
		published := publishedMatches[1]
		// parse to int
		publishedInt, err = strconv.Atoi(published)
		if err != nil {
			return 0, err
		}
	}
	return publishedInt, nil
}

func getEditions(text string) (int, error) {
	// assign editions
	editions := regexes.editionsRegex.FindStringSubmatch(text)[1]
	// convert to integer
	editionsInt, err := strconv.Atoi(editions)
	if err != nil {
		return 0, err
	}
	return editionsInt, nil
}

func getAuthors(e *colly.HTMLElement) []string {
	// create the authors array
	var authors []string
	// for each author
	e.ForEach(".authorName__container", func(i int, e *colly.HTMLElement) {
		// add the author to array
		authors = append(authors, strings.TrimSpace(e.Text))
	})
	return authors
}

func handleBooks(c *colly.Collector, books *[]book, bookCount *int) {
	// for each book result
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		// increment the counter for each book
		*bookCount++
		// get the book name
		name := e.ChildText(".bookTitle")

		// get authors array
		authors := getAuthors(e)

		// get text
		text := e.ChildText(".uitext.greyText.smallText")

		// parse and get details
		averageRatingFloat, err := getAverageRating(text)
		if err != nil {
			log.Println("error while parsing averageRating to float:", err, name)
			return
		}
		numberOfRatingsInt, err := getNumberOfRatings(text)
		if err != nil {
			log.Println("error while parsing numberOfRatings to int:", err, name)
			return
		}
		publishedInt, err := getPublished(text)
		if err != nil {
			log.Println("error while parsing publish date to int:", err, name)
			return
		}
		editionsInt, err := getEditions(text)
		if err != nil {
			log.Println("error while parsing editions to int:", err, name)
			return
		}

		// generate and add struct to books array
		*books = append(*books, book{
			Title:           name,
			Authors:         authors,
			AverageRating:   float32(averageRatingFloat),
			NumberOfRatings: numberOfRatingsInt,
			Published:       publishedInt,
			Editions:        editionsInt,
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

func writeToFile(data []book) {
	// generate json from struct
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("error while generating json: ", err)
	}
	// write results to file
	err = ioutil.WriteFile("results.json", file, 0644)
	if err != nil {
		log.Fatal("error while writing to file: ", err)
	}
	log.Println("You can find the results at results.json")
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
	handleBooks(c, &books, &bookCount)

	// start scraping
	startScraping(pageCountInt, query, c)

	// wait for all the threads to finish
	c.Wait()
	// write results to file
	writeToFile(books)
	// print the final book count
	log.Println("Parsed book count:", bookCount)
}
