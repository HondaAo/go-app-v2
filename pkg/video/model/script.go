package model

type Script struct {
	ScriptId  int    `json:"script_id"`
	VideoId   int    `json:"video_id"`
	Text      string `json:"text"`
	Ja        string `json:"ja"`
	TimeStamp int    `json:"timestamp"`
}
