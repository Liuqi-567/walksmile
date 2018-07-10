package models

import (
	"encoding/json"
	"fmt"
	"html/template"

	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

type SharedFileForFinder struct {
	ID         int           `json:"id"`
	Title      template.HTML `json:"title"`
	Hot        int           `json:"hot"`
	Category   int           `json:"category"`
	LinkCount  int           `json:"link_count"`
	UpdateTime string        `json:"update_time"`
	CreateTime string        `json:"create_time"`
}

type FinderResponse struct {
	SharedFileList []*SharedFileForFinder
	Total          int64
	TookTime       int64
}

var client *elastic.Client

func SearchSharedfile(q string, from int, size int, cat int) (*FinderResponse, error) {
	ctx := context.Background()
	if client == nil {
		client, _ = elastic.NewClient()
	}

	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("title"))
	hl = hl.PreTags("<font color=\"red\">").PostTags("</font>")

	termQuery := elastic.NewMatchQuery("title", q)
	query := elastic.NewBoolQuery()
	query = query.Must(termQuery)
	if cat != 0 {
		query = query.Filter(elastic.NewTermQuery("category", cat))
	}
	query = query.Filter(elastic.NewRangeQuery("link_count").Gt(0))

	searchResult, err := client.Search().
		Index("sharedfiles").
		Query(query).
		From(from).Size(size).
		Pretty(true).
		Highlight(hl).
		Do(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var finderResponse *FinderResponse = new(FinderResponse)
	finderResponse.Total = searchResult.TotalHits()
	finderResponse.TookTime = searchResult.TookInMillis
	var sharedfiles []*SharedFileForFinder
	for _, hit := range searchResult.Hits.Hits {
		sharedfile := new(SharedFileForFinder)
		if err := json.Unmarshal(*hit.Source, sharedfile); err != nil {
			return nil, err
		}
		if hl, found := hit.Highlight["title"]; found {
			sharedfile.Title = template.HTML(hl[0])
		}
		sharedfiles = append(sharedfiles, sharedfile)
	}
	finderResponse.SharedFileList = sharedfiles
	return finderResponse, nil

}
