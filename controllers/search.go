package controllers

import (
	"encoding/json"
	"fmt"
	"walksmile/conf"

	"math"
	"net/http"
	"strconv"
	"strings"
	"walksmile/helpers/xunlei/yike"
	"walksmile/models"

	"github.com/gin-gonic/gin"
)

func FileList(c *gin.Context) {
	q := c.Query("q")
	from := c.DefaultQuery("from", "0")
	size := c.DefaultQuery("size", "10")
	cat := c.DefaultQuery("cat", "0")
	f, _ := strconv.Atoi(from)
	s, _ := strconv.Atoi(size)
	ca, _ := strconv.Atoi(cat)
	finderResponse, err := models.SearchSharedfile(q, f, s, ca)
	if err != nil {
		fmt.Println(err.Error())
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	if conf.Conf.OpenCheck {
		go func() {
			for _, v := range finderResponse.SharedFileList {
				link := &models.Link{FkSharedFilesId: v.ID}
				links, err := link.FindByFkSharedFileId()
				if err != nil {
					continue
				}
				for _, vv := range links {
					linkqueue := &models.LinkQueue{
						UK:       vv.FkUK,
						Shareid:  vv.ShareID,
						ShortUrl: vv.ShortUrl,
					}
					err := linkqueue.Insert()
					if err != nil {
						continue
					}
				}
			}
		}()
	}

	pag := float64(f+1) / float64(s)
	pa := int(math.Ceil(pag))

	paginator := paginator(pa, s, finderResponse.Total)
	c.HTML(http.StatusOK, "list.html", gin.H{
		"title":     q + "-行笑网-网盘搜索-百度云搜索",
		"list":      finderResponse.SharedFileList,
		"total":     finderResponse.Total,
		"tooktime":  finderResponse.TookTime,
		"q":         q,
		"paginator": paginator,
		"cat":       ca,
	})

}

func FileDetail(c *gin.Context) {
	strid := c.Param("id")
	id, _ := strconv.Atoi(strid)

	//百度云资源
	sf := &models.SharedFile{ID: id}
	file, err := sf.FindByID()
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	} else {
		//热度+1
		go sf.UpdateHot()
	}

	linksandsharer, err := models.SelectLinksAndSharerBySharedfileId(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	//迅雷资源
	resources, _ := models.ResourcesBySharedfileID(models.GetDB(), uint(id))

	download := &yike.Resource{}

	if len(resources) > 0 {
		download = &yike.Resource{}
		entity := resources[0].LinkEntity
		err := json.Unmarshal([]byte(entity), download)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	// more
	more, _ := models.FindSFNearID(id)

	//news
	newfile := make([]*models.SharedFile, 0, 20)
	if len(Last24hourSharedFiles) > 20 {
		newfile = Last24hourSharedFiles[:20]
	}

	//评论
	// comment := &models.Comment{Fatherid: id, Kind: 0}
	// comments, err := comment.FindByFatherID()
	// if err != nil {
	// 	println(err)
	// }
	c.HTML(http.StatusOK, "detail.html", gin.H{
		//"comments": comments,
		//"recommend": recommend,
		"title":    file.Title + " 行笑网-百度云搜索-网盘搜索",
		"file":     file,
		"list":     linksandsharer,
		"more":     more,
		"newfile":  newfile,
		"download": download})

}

func Last24hour(c *gin.Context) {

	if Last24hourSharedFiles == nil {
		lsfs, err := models.FindLast24()
		if err != nil {
			c.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		Last24hourSharedFiles = lsfs
	}

	from := c.DefaultQuery("from", "0")
	size := c.DefaultQuery("size", "10")
	f, _ := strconv.Atoi(from)
	s, _ := strconv.Atoi(size)
	pag := float64(f+1) / float64(s)
	pa := int(math.Ceil(pag))
	total := int64(len(Last24hourSharedFiles))
	paginator := paginator(pa, s, total)
	if f >= int(total) {
		if int(total)%s > 0 {
			f = int(total)/s + 1
		} else {
			f = int(total) / s
		}
	}
	endindex := f + s
	if endindex >= int(total) {
		endindex = int(total - 1)
	}
	var cutList []*models.SharedFile
	if total > 0 {
		cutList = Last24hourSharedFiles[f:endindex]
	}

	c.HTML(http.StatusOK, "last24.html", gin.H{
		"title":     "最近24小时更新-行笑网-网盘搜索-百度云搜索 ",
		"list":      cutList,
		"total":     total,
		"paginator": paginator,
	})

}

func Comment(c *gin.Context) {
	referer := c.Request.Referer()
	if !strings.Contains(referer, "/") {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	referer2 := strings.Split(referer, "?")
	referer3 := strings.Split(referer2[0], "/")
	strID := referer3[len(referer3)-1]
	fileID, err := strconv.Atoi(strID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	content := c.PostForm("content")
	comment := &models.Comment{
		Name:     name,
		Email:    email,
		Content:  content,
		Fatherid: fileID,
	}
	err = comment.Insert()
	if err != nil {
		c.String(http.StatusOK, "评论失败"+err.Error())
		return
	}
	c.Redirect(303, "http://"+c.Request.Host+"/sharedfile/"+strID)

}
