package calendar

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

func GetCalendarConfig() Config {
	file, err := os.Open("./calendar.json")
	if errors.Is(err, os.ErrNotExist) {
		os.Create("./calendar.json")
		os.WriteFile("./calendar.json", []byte(`{"notion_api_key": "NOTION_API_KEY", "database_id": "DATABASE_ID", "calendars": [{"name":"somename", "url":"https://yourCalendarWebcalUrl"}]}`), 0644)
		log.Fatal("Please fill out the calendar.json file with your Notion API key, database ID, and calendar IDs")
	} else if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var calendarConfig Config
	json.Unmarshal(byteValue, &calendarConfig)

	return calendarConfig
}

type Config struct {
	NotionAPIKey string           `json:"notion_api_key"`
	DatabaseID   string           `json:"database_id"`
	Calendars    []CalendarConfig `json:"calendars"`
}

type CalendarConfig struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
