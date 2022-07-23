package model

type Video struct {
	VideoId   int    `json:"video_id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	View      int    `json:"view"`
	Category  string `json:"category"`
	Series    string `json:"series"`
	End       int    `json:"end"`
	Start     int    `json:"start" gorm:"default:0"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
}
