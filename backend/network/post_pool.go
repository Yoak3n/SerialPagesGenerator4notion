package network

import (
	"b2n3/backend/model"
	"b2n3/package/logger"
	"context"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PosterPool struct {
	client    *http.Client
	capacity  int
	remainder int
	datas     []*model.Data
	ctx       *context.Context
}

func NewPosterPool(datas []*model.Data, ctx *context.Context) *PosterPool {
	l := len(datas)
	return &PosterPool{
		client:    &http.Client{},
		datas:     datas,
		capacity:  l,
		remainder: int(math.Mod(float64(l), 3)),
		ctx:       ctx,
	}
}

func (p *PosterPool) postSingleData(req *http.Request) error {
	nreq := req
	res, err := p.client.Do(req)
	if err != nil || res.StatusCode != 200 {
		logger.DEBUG.Println("Too many requests, retry after 1 second", res.Body)
		time.Sleep(time.Second)
		// 估计会有BUG
		p.postRetry(*nreq)
	}
	runtime.EventsEmit(*p.ctx, "postProgress")
	return nil
}

// 考虑递归
func (p *PosterPool) postRetry(req http.Request) error {
	client := &http.Client{}
	res, err := client.Do(&req)
	if err != nil || res.StatusCode != 200 {
		return p.postRetry(req)
	}
	runtime.EventsEmit(*p.ctx, "postProgress")
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
	// 并发请求
	for i := p.remainder; i < p.capacity; i = i + 3 {
		go p.postSingleData(reqs[i])
		go p.postSingleData(reqs[i+1])
		go p.postSingleData(reqs[i+2])
		time.Sleep(time.Second)
	}
	for i := 0; i < p.remainder; i++ {
		go p.postSingleData(reqs[i])
	}
	return nil
}
