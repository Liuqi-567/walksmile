package yike

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"walksmile/models"

	"github.com/PuerkitoBio/goquery"
)

const MOVIE = "http://www.onek.cc/movie/00001/"

const TV = "http://www.onek.cc/tv/0001/"
const Catoon = "http://www.onek.cc/cartoon/00001/"

const detail = "http://www.onek.cc/tv/19856.html"

type Resource struct {
	Tap [][]*Item `json:"Tap"`
}

type Item struct {
	Title string            `json:Title`
	Links map[string]string `json:Links`
}

func main() {

	// urlstrs := util.ReadFileToArray("movie")
	// for _, v := range urlstrs {
	// 	Detail(v)
	// }

	// for i := 1; i <= 41; i++ {
	// 	resp, err := http.Get(Catoon + strconv.Itoa(i))
	// 	if err != nil {
	// 		println(err.Error())
	// 	}
	// 	doc, _ := goquery.NewDocumentFromResponse(resp)
	// 	doc.Find("a.thumbnail-img").Each(func(i int, sel *goquery.Selection) {
	// 		href, ok := sel.Attr("href")
	// 		if ok {
	// 			fmt.Println(href)
	// 		}
	// 	})
	// }

}

func Detail(urlstr string) {
	time.Sleep(time.Second)
	resp, _ := http.Get(urlstr)
	doc, _ := goquery.NewDocumentFromResponse(resp)

	title := doc.Find("div.info").Find("h1").Text()
	otitle := doc.Find("p.info-p").First().Find("span").Text()
	if strings.Contains(otitle, "又名：") {
		otitle = strings.Replace(otitle, "又名：", "", -1)
	} else {
		otitle = ""
	}
	image, _ := doc.Find("div.conver-img").Find("img").Attr("src")
	intr := doc.Find("p.info-p").Last().Text()
	fmt.Println(title, otitle, image, intr)
	res := &Resource{}
	doc.Find("ul.download-ul").Each(func(i int, sel *goquery.Selection) {
		var items []*Item
		sel.Find("li").Each(func(i int, se *goquery.Selection) {
			title := se.Find(".download-title").Text()
			item := &Item{Title: title}
			links := make(map[string]string, 0)
			se.Find(".download-btn").Find("a").Each(func(i int, s *goquery.Selection) {
				linkname := s.Text()
				linkhref, ok := s.Attr("href")
				if ok {
					links[linkname] = linkhref
				}
			})
			item.Links = links
			items = append(items, item)
		})
		res.Tap = append(res.Tap, items)
	})
	b, _ := json.Marshal(res)
	fmt.Print(string(b))

	resouces := &models.Resource{
		Title:      title,
		Otitle:     otitle,
		Image:      image,
		Intro:      intr,
		Kind:       1,
		LinkEntity: string(b),
		Originsite: urlstr,
	}

	sf := &models.SharedFile{}
	sf.Category = 7
	sf.Title = resouces.Title + resouces.Otitle

	err := sendSharedfile(resouces, sf)
	if err != nil {
		println(err.Error())
	}

}

func sendSharedfile(resouces *models.Resource, sf *models.SharedFile) error {
	sf1, _ := sf.FindByName()
	if sf1 == nil {
		err := sf.Insert()
		if err != nil {
			return err
		}
		resouces.SharedfileID = uint(sf.ID)
		err = resouces.Insert(models.GetDB())
		if err != nil {
			return err
		}
		sf.AddLinkCount()
		println("insert sharefile and link")
		return nil

	}

	link2, _ := models.ResourceByOriginsite(models.GetDB(), resouces.Originsite)
	if link2 == nil {
		resouces.SharedfileID = uint(sf1.ID)
		err := resouces.Insert(models.GetDB())
		if err != nil {
			return err
		}
		sf1.AddLinkCount()
		println("insert link just")
		return nil
	}

	println("link alread exist")
	return nil

}
