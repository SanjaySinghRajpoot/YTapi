package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SanjaySinghRajpoot/YTapi/models"
)

func InsertVideosToDatabase(db *sql.DB, videos []models.VideoInfo) error {

	// Bulk Insert Operation
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO videos (video_title, description, publish_time, thumbnail_url, channel, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
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
