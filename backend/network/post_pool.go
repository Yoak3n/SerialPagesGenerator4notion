package network

import (
	"b2n3/backend/model"
	"b2n3/package/logger"
	"errors"
	"log"
	"math"
	"net/http"
	"time"
)

type PosterPool struct {
	client    *http.Client
	capacity  int
	remainder int
	datas     []*model.Data
	count     chan int
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
	nreq := req
	res, err := p.client.Do(req)
	if res.StatusCode != 200 {
		logger.ERROR.Println("Too many requests, retry after 1 second", res.Body)
		log.Println(res.StatusCode, req.Header, req.Body)
		time.Sleep(time.Second)
		// 估计会有BUG

		p.postRetry(*nreq)
	}
	if err != nil {
		logger.ERROR.Println("Posting data failed: ", err)
	}
	p.count <- 1
	return nil
}

// 考虑递归
func (p *PosterPool) postRetry(req http.Request) error {

	client := &http.Client{}
	res, _ := client.Do(&req)
	if res.StatusCode != 200 {
		return p.postRetry(req)
	}
	p.count <- 1
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

func (p *PosterPool) Watch() {
	select {
	case <-p.count:
		logger.INFO.Println("Posting data add")
	case p.count <- 1:
		logger.INFO.Println("Posting data added")
	}
}
