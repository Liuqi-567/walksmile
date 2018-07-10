package controllers

import (
	"math"
	"regexp"
	"strings"
	"walksmile/models"
)

//分页方法，根据传递过来的页数，每页数，总数，返回分页的内容 7个页数 前 1，2，3，4，5 后 的格式返回,小于5页返回具体页数
func paginator(page, prepage int, nums int64) map[string]interface{} {
	var firstpage int //前一页地址
	var lastpage int  //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //最后5页
		start := totalpages - 5 + 1
		firstpage = page - 1
		lastpage = int(math.Min(float64(totalpages), float64(page+1)))
		pages = make([]int, 5)
		for i, _ := range pages {
			pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		pages = make([]int, 5)
		firstpage = page - 3
		for i, _ := range pages {
			pages[i] = start + i
		}
		firstpage = page - 1
		lastpage = page + 1
	default:
		pages = make([]int, int(math.Min(5, float64(totalpages))))
		for i, _ := range pages {
			pages[i] = i + 1
		}
		firstpage = int(math.Max(float64(1), float64(page-1)))
		lastpage = page + 1
		//fmt.Println(pages)
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["firstpage"] = firstpage
	paginatorMap["lastpage"] = lastpage
	paginatorMap["currpage"] = page
	return paginatorMap
}

//根据文件名称计算类别

var (
	c12, _ = regexp.Compile(`\.(ppt$|pptx)$`)
	c11, _ = regexp.Compile(`\.(xls$|xlsx)$`)
	c10, _ = regexp.Compile(`\.(doc$|docx)$`)
	c9, _  = regexp.Compile(`\.torrent$`)
	c8, _  = regexp.Compile(`\.txt$`)
	c7, _  = regexp.Compile(`\.(avi|mpg|rm|rmvb|mov|asf|wmv|mp4|divx|mkv|xvid|vob)$`)
	c6, _  = regexp.Compile(`\.exe$`)
	c5, _  = regexp.Compile(`\.(jpg|png|gif|tif|bmp)$`)
	c4, _  = regexp.Compile(`\.(wav|mp3|wma|ape)$`)
	c3, _  = regexp.Compile(`\.(rar|zip|7z|tar|gzip|bz2|jar|iso)$`)
	c2, _  = regexp.Compile(`\.pdf$`)
	c1, _  = regexp.Compile(`\.[a-z]+$`)
)

func Catagory_caculate(title string) int {

	if c12.MatchString(title) {
		return 12
	}
	if c11.MatchString(title) {
		return 11
	}
	if c10.MatchString(title) {
		return 10
	}
	if c9.MatchString(title) {
		return 9
	}
	if c8.MatchString(title) {
		return 8
	}
	if c7.MatchString(title) {
		return 7
	}
	if c6.MatchString(title) {
		return 6
	}
	if c5.MatchString(title) {
		return 5
	}
	if c4.MatchString(title) {
		return 4
	}
	if c3.MatchString(title) {
		return 3
	}
	if c2.MatchString(title) {
		return 2
	}
	if !c1.MatchString(title) || !strings.Contains(title, ".") {
		return 1
	}
	return 0

}

//向后台发送sharedfile和link
func SendSharedfile(link *models.Link, sf *models.SharedFile) error {
	sf1, _ := sf.FindByName()
	if sf1 == nil {
		err := sf.Insert()
		if err != nil {
			return err
		}
		link.FkSharedFilesId = sf.ID
		err = link.Insert()
		if err != nil {
			return err
		}
		println("insert sharefile and link")
		return nil

	}

	link2, _ := models.SelectLinkByKey(link)
	if link2 == nil {
		link.FkSharedFilesId = sf1.ID
		err := link.Insert()
		if err != nil {
			return err
		}
		println("insert link just")
		return nil
	}

	println("link alread exist")
	return nil

}
