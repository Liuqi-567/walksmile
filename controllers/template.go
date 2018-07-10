package controllers

import (
	"fmt"
	"html/template"
	"time"
)

//IsActive checks uri against currently active (uri, or nil) and returns "active" if they are equal
func IsActive(active interface{}, uri string) string {
	if s, ok := active.(string); ok {
		if s == uri {
			return "active"
		}
	}
	return ""
}

//DateTime prints timestamp in human format
func DateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//StringInSlice returns true if value is in list slice
func StringInSlice(value string, list []string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//把页面转化成from
func PageToFrom(page int) int {
	from := page*10 - 10
	return from

}

//把页面转化成from
func GoodPageToFrom(page int) int {
	from := page*20 - 20
	return from
}

func ToTemplate(ii string) template.URL {
	return template.URL(ii)
}
