package model

type Title struct{
	Type string `json:"type"`
	Text struct{
		Content string `json:"content"`
	} `json:"text"`
}