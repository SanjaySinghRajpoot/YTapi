package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

// Video represents the structure of a video
type Video struct {
	ID           int    `json:"id"`
	VideoTitle   string `json:"video_title"`
	Description  string `json:"description"`
	PublishTime  string `json:"publish_time"`
	ThumbnailURL string `json:"thumbnail_url"`
	Channel      string `json:"channel"`
	CreatedAt    string `json:"created_at"`
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

// get endpoint
// Handler function for the API endpoint
func getVideosHandler(w http.ResponseWriter, r *http.Request) {
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

	videos := make([]Video, 0)
	for rows.Next() {
		var v Video
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

// lets code a cron job which will run for 10 secs get the data and save it in the DB

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

	// videos := youtubeAPI.Data

	// Begin transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the bulk insert statement
	stmt, err := tx.Prepare(`
	 INSERT INTO videos (video_title, description, publish_time, thumbnail_url, channel, created_at)
	 VALUES ($1, $2, $3, $4, $5, $6)
    `)
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range youtubeAPI.Data["Video List"] {

		fmt.Println(obj)

		_, err = stmt.Exec(obj.Title, obj.Description, obj.PublishTime, obj.Image, obj.Channel, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(youtubeAPI)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))

	http.HandleFunc("/videos", getVideosHandler)
	http.ListenAndServe(":8080", nil)
}
