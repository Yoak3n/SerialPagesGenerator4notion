package util

import (
	"bytes"
	"math"
	"net/url"
	"strconv"
	"strings"
)

func ExtractID(input string) ([]byte, bool) {
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
			return id, false
		}
		// 如果是av号
	} else if bytes.HasPrefix(vid, prefix[1]) {
		id, ok = bytes.CutPrefix(vid, prefix[1])
		if ok {
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

func CheckVideoID(target string) ([]byte, int) {
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
	tid, isAv := ExtractID(id)
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

func Bvid2aid(bvid string) int64 {
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

func Aid2bvid(aid int64) string {
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
