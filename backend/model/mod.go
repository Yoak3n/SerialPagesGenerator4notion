package model

type Title struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type Select struct {
	Name string `json:"name"`
}
