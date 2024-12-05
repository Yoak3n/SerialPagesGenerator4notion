package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

var bangumi *Bangumi

type Bangumi struct {
	SeasonID     string    `json:"season_id"`
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
	id := checkBangumiID(input)
	if id == "" {
		return nil
	}
	bangumi = &Bangumi{}
	bangumi.SeasonID = id
	err := bangumi.getInfo()
	if err != nil {
		log.Println(err, "info")
		return nil
	}
	err = bangumi.getDetail()
	if err != nil {
		log.Println(err, "detail")
		return nil
	}
	return bangumi
}

func GetBangumiInfo() *Bangumi {
	return bangumi
}

func (b *Bangumi) getDetail() error {
	uri := fmt.Sprintf("https://api.bilibili.com/pgc/view/web/season?season_id=%s", b.SeasonID)
	res, err := http.Get(uri)
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
	return nil
}

func (b *Bangumi) getInfo() error {
	uri := fmt.Sprintf("https://api.bilibili.com/pgc/web/season/section?season_id=%s", b.SeasonID)
	res, err := http.Get(uri)
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
	subSection := result.Get("result.sub_section.episodes").Array()
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

func checkBangumiID(target string) string {
	if strings.HasPrefix(target, "https://www.bilibili.com/bangumi/play/") {
		uri, err := url.Parse(target)
		if err != nil {
			return ""
		}
		return extractBangumiID(strings.Split(uri.Path, "/")[3])
	} else {
		id := extractBangumiID(target)
		if id != "" {
			return id
		}
		log.Println("无效的值")
		return ""
	}
}
func extractBangumiID(target string) string {
	if strings.HasPrefix(target, "ss") {
		id, ok := strings.CutPrefix(target, "ss")
		if ok {
			return id
		}
		return ""
	} else {
		return target
	}
}
