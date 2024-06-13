package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	link         = "https://www.annauniv.edu/events.php"
	databaseName = "AnnaUnivJobAlert.db"
	botToken     = "Add Your BOT Token"
	chatIDStr       = "Add your chat ID"
)

func main() {
	// Setup logging
	logFile, err := os.OpenFile("job_alerts.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Convert chatIDStr to int64
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing chat ID: %v", err)
	}

	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Scrape job data from the website
	err = scrapeJobSite(db, chatID)
	if err != nil {
		log.Printf("Error scraping job site: %v", err)
		sendTelegramMessage(fmt.Sprintf("Error during execution: %v", err), chatID)
	}
}

func scrapeJobSite(db *sql.DB, chatID int64) error {
	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("error fetching job site: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("error parsing HTML: %w", err)
	}

	doc.Find("#graphic-design-2 tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() == 4 {
			title := strings.TrimSpace(cells.Eq(1).Text())
			departmentName := strings.TrimSpace(cells.Eq(2).Text())
			lastDate := strings.TrimSpace(cells.Eq(3).Text())
			linkString, _ := cells.Eq(1).Html()
			pdfLink := extractLink(linkString)
			encodedPDFLink := url.QueryEscape(pdfLink)

			existingJob, err := checkExistingJob(db, pdfLink)
			if err != nil {
				log.Printf("Error checking existing job: %v", err)
				return
			}

			if !existingJob {
				message := fmt.Sprintf(`
Anna University Job Alert

Title: %s
Department Name: %s
Last Date: %s
PDF LINK: %s
`, title, departmentName, lastDate, encodedPDFLink)

				err := sendTelegramMessage(message, chatID)
				if err != nil {
					log.Printf("Error sending telegram message: %v", err)
					return
				}

				err = insertJobEntry(db, title, departmentName, lastDate, pdfLink)
				if err != nil {
					log.Printf("Error inserting job entry: %v", err)
					return
				}

				log.Printf("Job entry added: %s", title)
			} else {
				log.Printf("Duplicate job entry found: %s", pdfLink)
			}
		}
	})

	return nil
}

func extractLink(linkString string) string {
	start := strings.Index(linkString, "href=\"") + len("href=\"")
	end := strings.Index(linkString[start:], "\"") + start
	return linkString[start:end]
}

func checkExistingJob(db *sql.DB, pdfLink string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM AnnaJobs WHERE pdf_link=? LIMIT 1)"
	err := db.QueryRow(query, pdfLink).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func insertJobEntry(db *sql.DB, title, departmentName, lastDate, pdfLink string) error {
	query := "INSERT INTO AnnaJobs (title, department_name, last_date, pdf_link) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, title, departmentName, lastDate, pdfLink)
	return err
}

func sendTelegramMessage(message string, chatID int64) error {
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("error creating telegram bot: %w", err)
	}

	_, err = b.Send(&tb.Chat{ID: chatID}, message)
	return err
}
