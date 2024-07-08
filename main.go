package main

import (
	"fmt"
	"slices"
	"sync"
)

func main() {
	fmt.Println("Starting notion calendar refresh")

	calendarConfig := GetCalendarConfig()

	var uids []string
	var cursor any = nil
	for true {
		calendarResponse := QueryNotionDatabase(calendarConfig, cursor)
		for _, item := range calendarResponse.Results {
			uidRichText := item.Properties["uid"].(map[string]interface{})["rich_text"].([]interface{})
			if len(uidRichText) == 0 {
				continue
			}
			uids = append(uids, item.Properties["uid"].(map[string]interface{})["rich_text"].([]interface{})[0].(map[string]interface{})["plain_text"].(string))
		}
		if calendarResponse.HasMore == false {
			break
		}
		cursor = calendarResponse.NextCursor
	}

	fmt.Println("Items already present: ", len(uids))
	var wg sync.WaitGroup
	for _, cal := range calendarConfig.Calendars {
		events := GetCalendar(cal.Url)
		fmt.Println("Events in ", cal.Name, "-", len(events))
		for _, event := range events {
			if !slices.Contains(uids, event.Uid) {
				wg.Add(1)
				go AddItemToCalendar(calendarConfig, cal.Name, event, &wg)
				uids = append(uids, event.Uid)
			}
		}
	}

	wg.Wait()
	fmt.Println("Finished notion calendar refresh")

}
