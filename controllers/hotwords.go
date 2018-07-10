package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"walksmile/models"
)

var Last24hourSharedFiles []*models.SharedFile
var Hotwords map[string][]string
var BackgroundPic string
var BackgroundPicArray []string

func InitDATA() {

	h := &models.Hotwords{ID: 1}
	h, err := h.SelectByID()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(h.Content), &Hotwords)
	if err != nil {
		panic(err)
	}

	pic := &models.Hotwords{}
	pic, err = pic.SelectByLastID()
	if err != nil {
		panic(err)
	}
	BackgroundPic = pic.Content
	lsfs, err := models.FindLast24()
	if err != nil {
		panic(err)
	}
	Last24hourSharedFiles = lsfs

	pics, err := WallPaperList()
	if err != nil {
		panic(err)
	}
	BackgroundPicArray = pics

}

const apiurl = "http://top.baidu.com/detail/list?boardid="

type BaiDuKeyWordsResp struct {
	List []struct {
		Keyword string `json:"keyword"`
	} `json:"list"`
}

func FetchHotword() error {
	typeList := make(map[string]string)
	typeList["电影"] = "26"
	typeList["电视"] = "4"
	typeList["美剧"] = "452"
	typeList["韩剧"] = "453"
	typeList["综艺"] = "19"
	typeList["动漫"] = "23"
	typeList["小说"] = "7"
	content := make(map[string][]string)
	for _, v := range typeList {
		resp, err := http.Get(apiurl + v)

		if err != nil {
			println(err.Error())
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
			continue
		}
		resp.Body.Close()
		baidu := &BaiDuKeyWordsResp{}
		err = json.Unmarshal(body, baidu)
		if err != nil {
			println(err.Error())
			continue
		}
		var keys []string
		for _, b := range baidu.List {
			keys = append(keys, b.Keyword)
		}
		if len(keys) > 0 {
			content[v] = keys
		}
	}
	con, err := json.Marshal(content)
	if err != nil {
		return err
	}
	Hotwords = content
	h := &models.Hotwords{ID: 1, Content: string(con)}
	err = h.Update()
	if err != nil {
		return err
	}

	resBing, err := http.Get("http://area.sinaapp.com/bingImg")
	resBing.Body.Close()
	if err != nil {
		return err
	}
	BackgroundPic = resBing.Request.URL.String()
	pic := &models.Hotwords{Content: BackgroundPic}
	err = pic.Insert()
	if err != nil {
		if !strings.Contains(err.Error(), "Duplicate entry") {
			return err
		}
	}

	pics, err := WallPaperList()
	if err != nil {
		return err
	}
	BackgroundPicArray = pics

	return nil
}

func Last24Update() error {
	lsfs, err := models.FindLast24()
	if err != nil {
		return err
	}
	Last24hourSharedFiles = lsfs
	return nil
}

func WallPaperList() ([]string, error) {
	wallpaper := &models.Hotwords{}
	wps, err := wallpaper.FindList()
	if err != nil {
		return nil, err
	}
	var pics []string
	for _, v := range wps {
		pics = append(pics, v.Content)
	}
	return pics, nil
}
