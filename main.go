package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/SanjaySinghRajpoot/YTapi/config"
	"github.com/SanjaySinghRajpoot/YTapi/controller"
	"github.com/SanjaySinghRajpoot/YTapi/models"
	"github.com/SanjaySinghRajpoot/YTapi/utils"
)

func main() {

	// connect to DB
	config.ConnectDB()

	var keys []string

	//  -> /REMOVE erase this from the string
	key1 := "AIzaSyDK3D7jgGxwEeGsWzWgMtgu3eLZPM5OeFA/REMOVE"
	key2 := "AIzaSyD1N1nQECXYAxmdwcRFV0tqkAa2zLkBFCQ/REMOVE"
	key3 := "AIzaSyBPaQHapf1F_NDU_Y73tKPrE457Gb-gKjM/REMOVE"

	keys = append(keys, key1, key2, key3)

	// Generate a random number between 0 and 2
	randomNumber := rand.Intn(len(keys))

	youtubeAPI := models.YouTubeAPI{
		Query:         "mrbeast",
		APIKey:        keys[randomNumber],
		Data:          make(map[string][]models.VideoInfo),
		Count:         0,
		NextPageToken: "",
	}

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
				err = utils.InsertVideosToDatabase(config.DB, vidInfo)
				if err != nil {
					log.Println(err)
				}
			default:
				// Channel is full, skip this fetch iteration
				fmt.Println("Skipping data fetch as the channel is full")
			}
		}
	}()

	http.HandleFunc("/videos", controller.GetVideosHandler)
	http.ListenAndServe(":8080", nil)
}
