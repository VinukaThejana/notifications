package models

// Notification contains the shape of the notification
type Notification struct {
	Message string `json:"message"`
	From    User   `json:"from"`
	To      User   `json:"to"`
}
