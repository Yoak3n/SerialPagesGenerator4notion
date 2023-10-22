package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func NewVideoInfo(in string) *VideoInfo {
	idBytes, typeID := checkVideoID(in)
	if idBytes == nil {
		log.Fatal("请检查输入是否正确")
		return nil
	}
	video = &VideoInfo{}
	switch typeID {
	case 0:
		// base 参数的意思是进制
		video.AID, _ = strconv.ParseInt(string(idBytes), 10, 64)

		video.BvID = aid2bvid(video.AID)
		log.Println("bvid", video.BvID)
	case 1:
		video.BvID = string(idBytes)
		video.AID = bvid2aid(video.BvID)
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

func GetVideoInfo()*VideoInfo{
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
	if seasonDisplay{
		episodes := jsonResult.Get("data.ugc_season.sections").Array()[0].Get("episodes").Array()
		for i := 0; i < len(episodes); i++ {
			v.Titles = append(v.Titles, episodes[i].Get("title").String())
		}
	}else{
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

func extractID(input string) ([]byte, bool) {
	/*
		返回值为获取到的真实视频ID及是否为av号
	*/
	vid := []byte(input)
	prefix := [][]byte{[]byte("BV"), []byte("av")}
	var id []byte
	ok := false

	// 如果是BV号
	if bytes.HasPrefix(vid, prefix[0]) {
		id, ok = bytes.CutPrefix(vid, prefix[0])
		if ok {
			fmt.Println("有BV前缀")
			return id, false
		}
		// 如果是av号
	} else if bytes.HasPrefix(vid, prefix[1]) {
		id, ok = bytes.CutPrefix(vid, prefix[1])
		if ok {
			log.Println("有av前缀", id)
			return id, true
		}
	}
	// 如果没有前缀标识，那么直接粗略检查是否是纯数字
	if _, err := strconv.Atoi(string(vid)); err == nil {
		return vid, true
	} else {
		return vid, false
	}
}

func checkVideoID(target string) ([]byte, int) {
	id := ""
	if strings.HasPrefix(target, "https://www.bilibili.com/video") {
		uri, err := url.Parse(target)
		if err != nil {
			return nil, -1
		}
		id = strings.Split(uri.Path, "/")[2]
	} else {
		id = target
	}
	tid, isAv := extractID(id)
	log.Println(string(tid), isAv)
	if isAv {
		return tid, 0
	} else {
		return tid, 1
	}
}

const (
	table string = "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
	xor   int64  = 177451812
	add   int64  = 8728348608
)

func bvid2aid(bvid string) int64 {
	tr := make(map[byte]int)
	for i := 0; i < 58; i++ {
		tr[table[i]] = i
	}
	s := [6]int{11, 10, 3, 8, 4, 6}
	return func(x string) int64 {
		var r int64
		for i := 0; i < 6; i++ {
			result := math.Pow(float64(58), float64(i))
			st := tr[x[s[i]-2]]
			r += int64(result) * int64(st)
		}
		return (r - add) ^ xor
	}(bvid)
}

func aid2bvid(aid int64) string {
	tr := make(map[byte]int)
	for i := 0; i < 58; i++ {
		tr[table[i]] = i
	}
	s := [6]int{11, 10, 3, 8, 4, 6}

	return func(x int64) string {
		x = (x ^ xor) + add
		r := []byte("BV1  4 1 7  ")
		for i := 0; i < 6; i++ {
			r[s[i]] = table[int(math.Mod(math.Floor(float64(x)/math.Pow(58, float64(i))), 58))]
		}
		return string(r)
	}(aid)
}
