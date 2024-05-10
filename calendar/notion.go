package calendar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apognu/gocal"
	"io"
	"net/http"
	"strings"
	"time"
)

func AddItemToCalendar(config Config, calendar string, event gocal.Event) bool {
	var notionCreatePageRequest NotionCreatePageRequest
	notionCreatePageRequest.Parent.DatabaseID = config.DatabaseID
	notionCreatePageRequest.Properties.Name.Title = append(notionCreatePageRequest.Properties.Name.Title, struct {
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
	}{Text: struct {
		Content string `json:"content"`
	}{Content: event.Summary}})
	notionCreatePageRequest.Properties.Date.Date.Start = event.Start.Format(time.RFC3339)
	notionCreatePageRequest.Properties.Date.Date.End = event.End.Format(time.RFC3339)
	notionCreatePageRequest.Properties.Calendar.MultiSelect = make([]struct {
		Name string `json:"name"`
	}, 0)
	notionCreatePageRequest.Properties.Calendar.MultiSelect = append(notionCreatePageRequest.Properties.Calendar.MultiSelect, struct {
		Name string `json:"name"`
	}{Name: calendar})
	notionCreatePageRequest.Properties.Tags.MultiSelect = make([]struct {
		Name string `json:"name"`
	}, 0)
	notionCreatePageRequest.Properties.Class.MultiSelect = make([]struct {
		Name string `json:"name"`
	}, 0)
	if event.URL == "" {
		notionCreatePageRequest.Properties.Link.URL = nil
	} else {
		notionCreatePageRequest.Properties.Link.URL = event.URL
	}
	var uidRichText RichText
	uidRichText.Type = "text"
	uidRichText.Text.Content = event.Uid
	notionCreatePageRequest.Properties.UID.RichText = append(notionCreatePageRequest.Properties.UID.RichText, uidRichText)
	description := strings.Split(event.Description, "\\n")
	for _, paragraph := range description {
		var child Child
		child.Object = "block"
		child.Type = "paragraph"
		var richText RichText
		richText.Type = "text"
		richText.Text.Content = paragraph
		child.Paragraph.RichText = append(child.Paragraph.RichText, richText)
		notionCreatePageRequest.Children = append(notionCreatePageRequest.Children, child)
	}
	jsonStr, _ := json.Marshal(notionCreatePageRequest)
	req, _ := http.NewRequest("POST", "https://api.notion.com/v1/pages", bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+config.NotionAPIKey)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

type NotionCreatePageRequest struct {
	Parent struct {
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Properties struct {
		Name struct {
			Title []struct {
				Text struct {
					Content string `json:"content"`
				} `json:"text"`
			} `json:"title"`
		} `json:"Name"`
		Date struct {
			Date struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"date"`
		} `json:"Date"`
		Tags struct {
			MultiSelect []struct {
				Name string `json:"name"`
			} `json:"multi_select"`
		} `json:"Tags"`
		Class struct {
			MultiSelect []struct {
				Name string `json:"name"`
			} `json:"multi_select"`
		} `json:"Class"`
		Calendar struct {
			MultiSelect []struct {
				Name string `json:"name"`
			} `json:"multi_select"`
		} `json:"Calendar"`
		Link struct {
			URL any `json:"url"`
		} `json:"Link"`
		UID struct {
			RichText []RichText `json:"rich_text"`
		} `json:"uid"`
	} `json:"properties"`
	Children []Child `json:"children"`
}

type RichText struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

type Child struct {
	Object    string `json:"object"`
	Type      string `json:"type"`
	Paragraph struct {
		RichText []RichText `json:"rich_text"`
	} `json:"paragraph"`
}

func QueryNotionDatabase(config Config, start any) NotionCalendarResponse {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", config.DatabaseID)

	var notionQueryDatabaseRequest any
	if start != nil {
		var req NotionQueryDatabaseRequest
		req.StartCursor = start
		notionQueryDatabaseRequest = req
	} else {
		notionQueryDatabaseRequest = struct{}{}
	}
	var jsonStr, _ = json.Marshal(notionQueryDatabaseRequest)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+config.NotionAPIKey)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var calendarResponse NotionCalendarResponse
	json.Unmarshal(body, &calendarResponse)
	return calendarResponse
}

type NotionQueryDatabaseRequest struct {
	StartCursor any `json:"start_cursor"`
}

type NotionCalendarResponse struct {
	Object         string `json:"object"`
	NextCursor     string `json:"next_cursor"`
	HasMore        bool   `json:"has_more"`
	Type           string `json:"type"`
	PageOrDatabase any    `json:"page_or_database"`
	Results        []struct {
		Object         string `json:"object"`
		ID             string `json:"id"`
		CreatedTime    string `json:"created_time"`
		LastEditedTime string `json:"last_edited_time"`
		Url            string `json:"url"`
		Archived       bool   `json:"archived"`
		CreatedBy      struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		Cover struct {
			Type     string `json:"type"`
			External struct {
				URL string `json:"url"`
			}
		} `json:"cover"`
		Icon struct {
			Type  string `json:"type"`
			Emoji string `json:"emoji"`
		} `json:"icon"`
		Parent struct {
			Type       string `json:"type"`
			DatabaseID string `json:"database_id"`
		} `json:"parent"`
		Properties map[string]any `json:"properties"`
	}
}
