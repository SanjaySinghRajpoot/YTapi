package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SanjaySinghRajpoot/YTapi/config"
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

	fmt.Println(videoList)

	return videoList, nil
}

// get endpoint
// Handler function for the API endpoint
func GetVideosHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := 10 // Change the limit according to your requirement

	offset := (page - 1) * limit

	rows, err := config.DB.Query("SELECT id, video_title, description, publish_time, thumbnail_url, channel, created_at FROM videos ORDER BY publish_time DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch videos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	videos := make([]models.Video, 0)
	for rows.Next() {
		var v models.Video
		err := rows.Scan(&v.ID, &v.VideoTitle, &v.Description, &v.PublishTime, &v.ThumbnailURL, &v.Channel, &v.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to scan video row", http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error in fetching videos", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(videos)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
