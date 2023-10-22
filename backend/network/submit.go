package network

import (
	"b2n3/backend/model"
	"math"
	"net/http"
)

func SubmitVideo(datas []*model.Data) {
	l := len(datas)
	remainder := int(math.Mod(float64(l), 3))

	for i := 0; i < l-remainder; i = i + 3 {
		go postSingleData(datas[i])
	}
}

func postSingleData(data *model.Data) {

	client := &http.Client{}
	// TODO
}
