package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Book struct {
	title string
	price string
}

func main() {

	//Create a csv file here
	file, err := os.Create("exportBooks.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//Write csv file headers
	headers := []string{"Title", "Price"}
	writer.Write(headers)

	//Create a colly collector becoz its the one sends http requests and transpile html
	collector := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	collector.OnHTML(".product_pod", func(h *colly.HTMLElement) {
		book := Book{}
		book.title = h.ChildAttr(".image_container img", "alt")
		book.price = h.ChildText(".price_color")
		csvRow := []string{book.title, book.price}
		writer.Write(csvRow)
	})

	collector.OnHTML(".next > a", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		collector.Visit(nextPage)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	collector.Visit("https://books.toscrape.com/")
	//Scaping is done
	fmt.Println("Done")

}
