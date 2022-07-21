package model

type Video struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	View      int    `json:"view"`
	Category  string `json:"category"`
	Series    string `json:"series"`
	EndTime   int    `json:"end_time"`
	Start     int    `json:"start" gorm:"default:0"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
}
