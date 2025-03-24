package network

import (
	"b2n3/backend/model"

	"context"
)

func SubmitVideo(datas []*model.Data, ctx *context.Context) {
	pool := NewPosterPool(datas, ctx)
	pool.Start()
}
