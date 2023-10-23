package network

import (
	"net/http"
	"time"
)

type PosterClients struct {
	clients [3]*http.Client
	status  [3]chan int32
}

func NewPosterClients() *PosterClients {
	// 创建三个客户端
	clients := [3]*http.Client{}
	for i := 0; i < 3; i++ {
		clients[i] = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 60 * time.Second,
			},
		}
	}
	return &PosterClients{
		clients: clients,
		status:  [3]chan int32{},
	}
}

func (c *PosterClients) GetClient(index int) {
	c.status[index] = make(chan int32)
	go func(index int) {
		c.status[index] <- 1
	}(index)
}
func (c *PosterClients) loadBlanced() *http.Client {

	select {
	case s := <-c.status[0]:
		if s == 1 {
			c.status[0] <- 0
			return c.clients[0]
		}

	case s := <-c.status[1]:
		if s == 1 {
			c.status[1] <- 0
			return c.clients[1]
		}

	case s := <-c.status[2]:
		if s == 1 {
			c.status[2] <- 0
			return c.clients[2]
		}
	}
	return nil
}
