package network

import (
	"b2n3/backend/model"
	"b2n3/package/logger"
	"context"
	"errors"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PosterPool struct {
	client    *http.Client
	capacity  int
	remainder int
	datas     []*model.Data
	ctx       *context.Context
	wg        *sync.WaitGroup
}

func NewPosterPool(datas []*model.Data, ctx *context.Context) *PosterPool {
	l := len(datas)
	wg := &sync.WaitGroup{}
	wg.Add(l)
	return &PosterPool{
		client:    &http.Client{},
		datas:     datas,
		capacity:  l,
		remainder: int(math.Mod(float64(l), 3)),
		ctx:       ctx,
		wg:        wg,
	}
}

func (p *PosterPool) postSingleData(req *http.Request, index int) error {
	nreq := req
	res, err := p.client.Do(req)
	if res.StatusCode != 200 || err != nil {
		logger.INFO.Println("Too many requests, retry after 1 second")
		time.Sleep(time.Second)
		p.postRetry(*nreq, index)
		return errors.New("too many requests")
	}
	// log.Println("Posting data success", index)
	p.wg.Done()
	runtime.EventsEmit(*p.ctx, "postProgress", index)
	return nil
}

// 考虑递归
func (p *PosterPool) postRetry(req http.Request, index int) error {
	logger.INFO.Println("Posting data failed, retrying...", index)
	client := &http.Client{}
	res, _ := client.Do(&req)
	if res.StatusCode != 200 {
		time.Sleep(time.Second)
		return p.postRetry(req, index)
	}
	p.wg.Done()
	return nil
}

func (p *PosterPool) Start() error {
	reqs := make([]*http.Request, 0)
	for i := 0; i < p.capacity; i++ {
		req, err := generateSinglePostRequest(p.datas[i])
		if err != nil {
			log.Println(err)
		}
		reqs = append(reqs, req)
	}
	// 并发请求
	for i := p.remainder; i < p.capacity; i = i + 3 {
		go p.postSingleData(reqs[i], i+1)
		go p.postSingleData(reqs[i+1], i+2)
		go p.postSingleData(reqs[i+2], i+3)
		time.Sleep(time.Second)
	}
	for i := 0; i < p.remainder; i++ {
		go p.postSingleData(reqs[i], i+1)
	}
	p.wg.Wait()
	return nil
}
