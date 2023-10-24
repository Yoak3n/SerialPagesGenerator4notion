package network

import (
	"b2n3/backend/model"
	"b2n3/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SubmitVideo(datas []*model.Data) {

	pool := NewPosterPool(datas)
	go pool.Watch()
	pool.Start()
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Conf.Token))

	return req, nil
	// TODO
}
