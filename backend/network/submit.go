package network

import (
	"b2n3/backend/model"
	"b2n3/config"
	"bytes"
	"encoding/json"
	"net/http"
)

func SubmitVideo(datas []*model.Data) {

	pool := NewPosterPool(datas)

}

func generateSinglePostRequest(data *model.Data) (*http.Request, error) {

	b, _ := json.Marshal(data)
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages/", buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Conf.Token)

	return req, nil
	// TODO
}
