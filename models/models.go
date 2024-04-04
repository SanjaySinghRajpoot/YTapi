package models

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
