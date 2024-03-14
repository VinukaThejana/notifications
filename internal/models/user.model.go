package models

// User defines the shape of the user
type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
