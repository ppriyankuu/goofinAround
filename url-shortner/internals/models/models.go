package models

type URL struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Original string `json:"original"`
	Slug     string `json:"slug"`
}

type Analytics struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	URLID       uint   `json:"url_id"`
	IP          string `json:"ip"`
	DeviceType  string `json:"device_type"`
	Geolocation string `json:"geolocation"`
}
