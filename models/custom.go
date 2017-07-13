package models

type Text struct {
	Content string `json:"content"`
}

type TextCustom struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	Text Text `json:"text"`
}
