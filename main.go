package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/YTapi/config"
)

type YouTubeAPI struct {
	Query         string
	APIKey        string
	NextPageToken string
	Data          map[string][]VideoInfo
	Count         int
}

type VideoInfo struct {
	Index       int
	Title       string
	Description string
	PublishTime string
	Channel     string
	Image       string
}

func (y *YouTubeAPI) getLoadingStats() error {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&maxResults=100&type=video&eventType=completed&order=date&key=%s", y.Query, y.APIKey)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return err
	}

	if nextPageToken, ok := data["nextPageToken"].(string); ok {
		y.NextPageToken = nextPageToken
	}

	if items, ok := data["items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				snippet := itemMap["snippet"].(map[string]interface{})
				videoTitle := snippet["title"].(string)
				image := snippet["thumbnails"].(map[string]interface{})["high"].(map[string]interface{})["url"].(string)
				publishTime := snippet["publishedAt"].(string)
				channel := snippet["channelTitle"].(string)
				description := snippet["description"].(string)

				y.Data["Video List"] = append(y.Data["Video List"], VideoInfo{
					Index:       y.Count,
					Title:       videoTitle,
					Description: description,
					PublishTime: publishTime,
					Channel:     channel,
					Image:       image,
				})

				y.Count++
			}
		}
	}

	return nil
}

func main() {

	// connect to DB
	config.ConnectDB()

	youtubeAPI := YouTubeAPI{
		Query:         "mrbeast",
		APIKey:        "AIzaSyBPaQHapf1F_NDU_Y73tKPrE457Gb-gKjM",
		Data:          make(map[string][]VideoInfo),
		Count:         0,
		NextPageToken: "",
	}

	err := youtubeAPI.getLoadingStats()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	jsonData, err := json.Marshal(youtubeAPI)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
