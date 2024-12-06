package api

import (
	"b2n3/backend/model"
	"b2n3/backend/network"
	"b2n3/config"
	"b2n3/package/logger"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

var bangumi *Bangumi

type Bangumi struct {
	SeasonID     string    `json:"season_id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	BangumiCover string    `json:"bangumi_cover"`
	Main         []Section `json:"main"`
	Sub          []Section `json:"sub"`
}

type Section struct {
	AID        int64  `json:"aid"`
	LongTitle  string `json:"long_title"`
	ShortTitle string `json:"short_title"`
	Cover      string `json:"cover"`
}

func NewBangumiInfo(input string) *Bangumi {
	id, t := checkBangumiID(input)
	if id == "" {
		return nil
	}
	bangumi = &Bangumi{
		Main: make([]Section, 0),
	}
	bangumi.SeasonID = id
	if t == "" {
		t = "ep"
	}
	bangumi.Type = t
	logger.INFO.Println(bangumi.Main)
	err := bangumi.getDetail()
	if err != nil {
		logger.ERROR.Println(err)
		return nil
	}
	return bangumi
}

func GetBangumiInfo() *Bangumi {
	return bangumi
}

func SubmitBangumiInfo(ctx *context.Context, in string) []*model.Data {
	datas := initBangumiBody(in)
	network.SubmitVideo(datas, ctx)
	return nil
}

func initBangumiBody(in string) (datas []*model.Data) {
	parent := &model.Parent{
		Type:       "database_id",
		DatabaseID: config.Conf.DatabaseID,
	}
	if bangumi == nil {
		bangumi = NewBangumiInfo(in)
	}
	name := strings.ReplaceAll(bangumi.Title, ",", "，")
	for index, episode := range bangumi.Main {
		properties := &model.Properties{
			Episode: model.Episode{
				Number: index + 1,
			},
			EpisodeName: *genEpisodeName(&episode.LongTitle),
			Name: model.Name{
				Select: model.Select{
					Name: name,
				},
			},
		}
		data := &model.Data{
			Parent:     parent,
			Properties: properties,
		}
		datas = append(datas, data)
	}
	return
}

func (b *Bangumi) getDetail() error {
	client := &http.Client{}
	uri := ""
	if b.Type == "ep" {
		uri = fmt.Sprintf("https://api.bilibili.com/pgc/view/web/season/?ep_id=%s", b.SeasonID)
	} else if b.Type == "ss" {
		uri = fmt.Sprintf("https://api.bilibili.com/pgc/view/web/season?season_id=%s", b.SeasonID)
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", config.Conf.Cookie)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	result := gjson.ParseBytes(buf)
	code := result.Get("code").Int()
	if code != 0 {
		return errors.New(result.Get("message").String())
	}
	b.BangumiCover = result.Get("result.cover").String()
	b.Title = result.Get("result.season_title").String()
	episodes := result.Get("result.episodes").Array()
	for _, r := range episodes {
		badge := r.Get("badge").String()
		if badge != "预告" {
			single := &Section{
				r.Get("aid").Int(),
				r.Get("long_title").String(),
				r.Get("title").String(),
				r.Get("cover").String(),
			}
			b.Main = append(b.Main, *single)
		}
	}
	return nil
}

func (b *Bangumi) getInfo() error {
	uri := fmt.Sprintf("https://api.bilibili.com/pgc/web/season/section?season_id=%s", b.SeasonID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", config.Conf.Cookie)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	result := gjson.ParseBytes(buf)
	code := result.Get("code").Int()
	mainSection := result.Get("result.main_section.episodes").Array()
	subSection := result.Get("result.section.episodes").Array()
	if code != 0 {
		return errors.New(result.Get("message").String())
	}
	mainEpisodes := make([]Section, 0)
	subEpisodes := make([]Section, 0)
	for _, r := range mainSection {
		singleEpisode := Section{r.Get("aid").Int(), r.Get("long_title").String(), r.Get("title").String(), r.Get("cover").String()}
		mainEpisodes = append(mainEpisodes, singleEpisode)
	}
	for _, r := range subSection {
		singleEpisode := Section{r.Get("aid").Int(), r.Get("long_title").String(), r.Get("title").String(), r.Get("cover").String()}
		subEpisodes = append(subEpisodes, singleEpisode)
	}
	b.Main = mainEpisodes
	b.Sub = subEpisodes
	return nil
}

func checkBangumiID(target string) (string, string) {
	if strings.HasPrefix(target, "https://www.bilibili.com/bangumi/play/") {
		uri, err := url.Parse(target)
		if err != nil {
			return "", ""
		}
		return extractBangumiID(strings.Split(uri.Path, "/")[3])
	} else {
		return extractBangumiID(target)
	}
}
func extractBangumiID(target string) (string, string) {
	bangumiPrefix := []string{"ss", "ep"}
	for _, prefix := range bangumiPrefix {
		if strings.HasPrefix(target, prefix) {
			id, ok := strings.CutPrefix(target, prefix)
			if ok {
				return id, prefix
			}
			return "", ""
		} else {
			return target, ""
		}
	}
	return "", ""
}
