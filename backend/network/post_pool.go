package network

import (
	"b2n3/backend/model"
	"errors"
	"math"
	"net/http"
)

type PosterPool struct {
	client    *http.Client
	capacity  int
	remainder int
	datas     []*model.Data
}

func NewPosterPool(datas []*model.Data) *PosterPool {
	l := len(datas)
	return &PosterPool{
		client:    &http.Client{},
		datas:     datas,
		capacity:  l,
		remainder: int(math.Mod(float64(l), 3)),
	}
}

func (p *PosterPool) postSingleData(req *http.Request) error {
	res, err := p.client.Do(req)
	if res.StatusCode != 200 || res.StatusCode == 429 || err != nil {
		p.postRetry(req)
	}
	return nil
}

// 考虑递归
func (p *PosterPool) postRetry(req *http.Request) error {
	return nil
}

func (p *PosterPool) Start() error {
	reqs := make([]*http.Request, 0)
	for i := 0; i < p.capacity; i++ {
		req, err := generateSinglePostRequest(p.datas[i])
		if err != nil {
			return errors.New("generate post request failed")
		}
		reqs = append(reqs, req)
	}
	for i := p.remainder; i < p.capacity; i = i + 3 {
		go p.postSingleData(reqs[i])
	}
}
