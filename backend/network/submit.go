package network

import (
	"b2n3/backend/model"
	"b2n3/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func SubmitVideo(datas []*model.Data, ctx *context.Context) {
	pool := NewPosterPool(datas, ctx)
	pool.Start()
}

func generateSinglePostRequest(data *model.Data) (*http.Request, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
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
