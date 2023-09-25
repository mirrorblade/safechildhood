package domain

import "time"

type Complaint struct {
	ID               any       `json:"id" db:"id"`
	Coordinates      string    `json:"coordinates" db:"coordinates"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`
	PhotosPath       string    `json:"photos_path" db:"photos_path"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
