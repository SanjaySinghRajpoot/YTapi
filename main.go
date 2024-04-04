package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SanjaySinghRajpoot/YTapi/config"
	"github.com/SanjaySinghRajpoot/YTapi/controller"
	"github.com/SanjaySinghRajpoot/YTapi/models"
)

func insertVideosToDatabase(db *sql.DB, videos []models.VideoInfo) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO videos (video_title, description, publish_time, thumbnail_url, channel, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (video_title, created_at) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, video := range videos {
		_, err := stmt.Exec(video.Title, video.Description, video.PublishTime, video.Image, video.Channel, time.Now())
		if err != nil {
			return fmt.Errorf("error inserting video into database: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	fmt.Println("Data inserted into database successfully.")
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

// lets code a cron job which will run for 10 secs get the data and save it in the DB

func main() {

	// connect to DB
	config.ConnectDB()

	youtubeAPI := models.YouTubeAPI{
		Query:         "mrbeast",
		APIKey:        "AIzaSyBPaQHapf1F_NDU_Y73tKPrE457Gb-gKjM",
		Data:          make(map[string][]models.VideoInfo),
		Count:         0,
		NextPageToken: "",
	}

	// err := youtubeAPI.getLoadingStats()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// // videos := youtubeAPI.Data

	// // Begin transaction
	// tx, err := config.DB.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Prepare the bulk insert statement
	// stmt, err := tx.Prepare(`
	//  INSERT INTO videos (video_title, description, publish_time, thumbnail_url, channel, created_at)
	//  VALUES ($1, $2, $3, $4, $5, $6)
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, obj := range youtubeAPI.Data["Video List"] {

	// 	fmt.Println(obj)

	// 	_, err = stmt.Exec(obj.Title, obj.Description, obj.PublishTime, obj.Image, obj.Channel, time.Now())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// // Commit the transaction
	// err = tx.Commit()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// jsonData, err := json.Marshal(youtubeAPI)
	// if err != nil {
	// 	fmt.Println("Error marshalling JSON:", err)
	// 	return
	// }

	// fmt.Println(string(jsonData))

	// Creating sync channels and routines
	concurrencyLimit := 5
	resultChan := make(chan int, concurrencyLimit)

	// Start a Goroutine to continuously fetch data
	ticker := time.Tick(time.Minute)
	go func() {
		for range ticker {
			select {
			case resultChan <- 0:
				// Channel is available, proceed to fetch data
				vidInfo, err := controller.GetLoadingStats(youtubeAPI)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				// Insert fetched videos into the database
				err = insertVideosToDatabase(config.DB, vidInfo)
				if err != nil {
					log.Println(err)
				}
			default:
				// Channel is full, skip this fetch iteration
				fmt.Println("Skipping data fetch as the channel is full")
			}
		}
	}()

	http.HandleFunc("/videos", getVideosHandler)
	http.ListenAndServe(":8080", nil)
}
