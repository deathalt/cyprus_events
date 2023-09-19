package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// This will set up a tick to trigger every hour
	for range time.Tick(1 * time.Hour) {
		scrapeEvents()
	}
}

func notify(title string, category string, location string,
	price string, eventDate string, buyTicketsURL string) {

	// Structure the data
	message := fmt.Sprintf("#%s\n\n Event: %s\n\nLocation: %s\nPrice: %s\nEvent Date: %s\n\nBuy Tickets URL: %s",
		category, title, location, price, eventDate, buyTicketsURL)

	data := map[string]string{
		"chat_id": os.Getenv("TG_CHAT_ID"), // Assuming chatID is defined somewhere in your code
		"text":    message,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	// Send the JSON data as a POST request
	apiEndpoint := "https://api.telegram.org/bot" + os.Getenv("TG_BOT_TOKEN") + "/sendMessage" // Replace with your API endpoint
	_, err = http.Post(apiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Event sent: %s\n", title)

}

func scrapeEvents() {
	// Open or create an SQLite database file
	db, err := sql.Open("sqlite3", os.Getenv("SQLITE_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a table to store the scraped data
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		title TEXT,
		category TEXT,
		location TEXT,
		price TEXT,
		event_date TEXT,
		buy_tickets_url TEXT,
		PRIMARY KEY (title, category, location, price, event_date, buy_tickets_url)
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// URL of the web page to scrape
	baseUrl := "https://www.soldoutticketbox.com"
	homeEn := baseUrl + "/en/home"
	// Make an HTTP GET request to the URL
	response, err := http.Get(homeEn)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find all elements with class "homeBoxEvent"
	document.Find("div.homeBoxEvent").Each(func(index int, element *goquery.Selection) {
		// Extract data from each "homeBoxEvent" div
		title := element.Find("h2 a").Text()
		category := element.Find("a.h3Style").First().Text()
		location := element.Find("p.blackSmall span").First().Text()
		price := element.Find("p.blackSmall span").Last().Text()

		// Extract event date and clean up the text
		eventDate := element.Find("p.blackSmall span").Eq(1).Text()
		eventDate = strings.TrimSpace(strings.Split(eventDate, "(")[0])
		eventDate = strings.ReplaceAll(eventDate, " -", "-")

		// Use a regular expression to replace multiple spaces with a single space
		re := regexp.MustCompile(`\s+`)
		eventDate = re.ReplaceAllString(eventDate, " ")

		buyTicketsURL, _ := element.Find("a.buyTicket").Attr("href")
		buyTicketsURL = baseUrl + buyTicketsURL

		// Insert the scraped data into the SQLite database
		insertDataSQL := `
		INSERT OR IGNORE INTO events (title, category, location, price, event_date, buy_tickets_url)
		VALUES (?, ?, ?, ?, ?, ?)
		`
		result, err := db.Exec(insertDataSQL, title, category, location, price, eventDate, buyTicketsURL)
		if err != nil {
			log.Fatal(err)
		} else {
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				log.Fatal(err)
			} else if rowsAffected > 0 {
				notify(title, category, location, price, eventDate, buyTicketsURL)
			} else if rowsAffected == 0 {
				log.Printf("No new events found")
			}
		}

	})
}
