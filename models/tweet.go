package models

import "time"

type Tweet struct {
	Message string `bson:"message" json:"message"`
}

type SavedTweet struct {
	UserID  string    `bson:"user_id" json:"user_id,omitempty"`
	Message string    `bson:"message" json:"message,omitempty"`
	Date    time.Time `bson:"date" json:"date,omitempty"`
}
