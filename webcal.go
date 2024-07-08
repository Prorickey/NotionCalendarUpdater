package main

import (
	"github.com/apognu/gocal"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetCalendar(webcalURL string) []gocal.Event {
	webcalURL = strings.Replace(webcalURL, "webcal://", "https://", 1)
	req, err := http.NewRequest("GET", webcalURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	calendar := gocal.NewParser(strings.NewReader(string(body)))
	calendar.Strict = gocal.StrictParams{
		Mode: gocal.StrictModeFailAttribute,
	}
	start, end := time.Now(), time.Now().AddDate(1, 0, 0)
	calendar.Start, calendar.End = &start, &end
	calendar.Parse()
	return calendar.Events
}
