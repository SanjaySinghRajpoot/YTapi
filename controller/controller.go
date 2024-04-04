package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/YTapi/models"
)

func GetLoadingStats(y models.YouTubeAPI) ([]models.VideoInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&maxResults=100&type=video&eventType=completed&order=date&key=%s", y.Query, y.APIKey)

	var videoList []models.VideoInfo

	response, err := http.Get(url)
	if err != nil {
		return videoList, err
	}
	defer response.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return videoList, err
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

				videoList = append(videoList, models.VideoInfo{
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

	return videoList, nil
}
