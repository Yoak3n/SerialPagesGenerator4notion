package api

import (
	"b2n3/backend/model"
	"b2n3/backend/network"
	"b2n3/config"
	"b2n3/package/util"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

var video *VideoInfo

type VideoInfo struct {
	Name      string   `json:"name"`
	AID       int64    `json:"aid"`
	BvID      string   `json:"bvid"`
	Titles    []string `json:"titles"`
	Cover     string   `json:"cover"`
	VideoView int64    `json:"video_view"`
	VideoTime int64    `json:"video_time"`
}

// GetVideoInfo
func NewVideoInfo(in string) *VideoInfo {
	idBytes, typeID := util.CheckVideoID(in)
	if idBytes == nil {
		log.Fatal("请检查输入是否正确")
		return nil
	}
	video = &VideoInfo{}
	switch typeID {
	case 0:
		// base 参数的意思是进制
		video.AID, _ = strconv.ParseInt(string(idBytes), 10, 64)

		video.BvID = util.Aid2bvid(video.AID)
		log.Println("bvid", video.BvID)
	case 1:
		video.BvID = string(idBytes)
		video.AID = util.Bvid2aid(video.BvID)
		log.Println("aid", video.AID)
	default:
		log.Fatal("视频数据类型错误")
		return nil
	}
	err := video.getVideoTitle()
	if err != nil {
		return nil
	}
	err = video.getVideoData()
	if err != nil {
		return video
		// 接口失效，继续观望
	}
	return video
}

func GetVideoInfo() *VideoInfo {
	return video
}

func (v *VideoInfo) getVideoTitle() error {
	uri := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", v.BvID)
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	jsonResult := gjson.ParseBytes(all)

	code := jsonResult.Get("code").Int()

	seasonDisplay := jsonResult.Get("data.is_season_display").Bool()
	v.Name = jsonResult.Get("data.title").String()
	v.Cover = jsonResult.Get("data.pic").String()
	if code != 0 {
		log.Fatalln("请求错误")
		return err
	}
	if seasonDisplay {
		episodes := jsonResult.Get("data.ugc_season.sections").Array()[0].Get("episodes").Array()
		for i := 0; i < len(episodes); i++ {
			v.Titles = append(v.Titles, episodes[i].Get("title").String())
		}
	} else {
		pages := jsonResult.Get("data.pages").Array()
		for i := 0; i < len(pages); i++ {
			v.Titles = append(v.Titles, pages[i].Get("part").String())
		}
	}
	return nil
}

func (v *VideoInfo) getVideoData() error {
	uri := fmt.Sprintf("https://api.bilibili.com/archive_stat/stat?aid=%d", v.AID)
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	jsonResult := gjson.ParseBytes(all)
	code := jsonResult.Get("code").Int()
	if code != 0 {
		return errors.New("has no data")
	}
	v.VideoTime = jsonResult.Get("data.vt").Int()
	v.VideoView = jsonResult.Get("data.vv").Int()
	return nil
}

// SubmitVideoInfo

func SumbitVideo(ctx *context.Context) []*model.Data {
	datas := initVideoBody()
	network.SubmitVideo(datas, ctx)
	return datas
}

func initVideoBody() (datas []*model.Data) {
	parent := &model.Parent{
		Type:       "database_id",
		DatabaseID: config.Conf.DatabaseID,
	}

	for episode, episodeName := range video.Titles {
		properties := &model.Properties{
			Episode: model.Episode{
				Number: episode + 1,
			},
			EpisodeName: *genEpisodeName(&episodeName),
			Name: model.Name{
				Select: struct {
					Name string "json:\"name\""
				}{
					Name: video.Name,
				},
			},
		}

		data := &model.Data{
			Parent:     *parent,
			Properties: *properties,
		}

		datas = append(datas, data)
	}

	return

}

func genEpisodeName(name *string) *model.EpisodeName {
	titles := make([]model.Title, 0)
	title := &model.Title{
		Type: "text",
		Text: struct {
			Content string "json:\"content\""
		}{
			Content: *name,
		},
	}
	titles = append(titles, *title)

	episodeName := &model.EpisodeName{
		Title: titles,
	}

	return episodeName
}
