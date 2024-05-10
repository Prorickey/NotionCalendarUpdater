package main

import (
	"fmt"
	"notionCalendarUpdater/calendar"
	"slices"
)

func main() {
	fmt.Println("Starting notion calendar refresh")

	calendarConfig := calendar.GetCalendarConfig()

	var uids []string
	var cursor any = nil
	for true {
		calendarResponse := calendar.QueryNotionDatabase(calendarConfig, cursor)
		for _, item := range calendarResponse.Results {
			uids = append(uids, item.Properties["uid"].(map[string]interface{})["rich_text"].([]interface{})[0].(map[string]interface{})["plain_text"].(string))
		}
		if calendarResponse.HasMore == false {
			break
		}
		cursor = calendarResponse.NextCursor
	}

	fmt.Println("Items already present: ", len(uids))
	for _, cal := range calendarConfig.Calendars {
		events := calendar.GetCalendar(cal.Url)
		fmt.Println("Events in ", cal.Name, "-", len(events))
		for _, event := range events {
			if !slices.Contains(uids, event.Uid) {
				good := calendar.AddItemToCalendar(calendarConfig, cal.Name, event)
				uids = append(uids, event.Uid)
				if good {
					fmt.Printf("Added %s to %s\n", event.Summary, cal.Name)
				} else {
					fmt.Printf("Failed to add %s to %s\n", event.Summary, cal.Name)
					break
				}
			}
		}
	}

}
