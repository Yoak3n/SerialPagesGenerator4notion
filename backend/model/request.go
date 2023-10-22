package model

type Data struct {
	Properties `json:"properties"`
	Parent     `json:"parent"`
}

type Parent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
}

type Properties struct {
	Episode
	EpisodeName
	Name
}

type Episode struct {
	Number int `json:"number"`
}

type EpisodeName struct {
	Title []Title `json:"title"`
}

type Name struct {
	Select struct {
		Name string `json:"name"`
	} `json:"select"`
}
