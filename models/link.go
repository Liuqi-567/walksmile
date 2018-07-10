package models

import (
	"errors"
	"strconv"
	"strings"
)

type Link struct {
	ShareID         int
	FkUK            int
	IsValid         int
	UpdateTimes     int
	FkSharedFilesId int
	ShortUrl        string
}

type LinkAndSharer struct {
	Link
	Sharer
}

func (l *Link) Insert() error {
	_, err := db.Exec("insert into links (shareid,fk_uk,fk_sharedfiles_id,shorturl) values (?,?,?,?)", &l.ShareID, &l.FkUK, &l.FkSharedFilesId, l.ShortUrl)
	if err != nil {
		return err
	}
	return nil

}
func (l *Link) FindByUpdateTimes() ([]*Link, error) {
	rows, err := db.Query("select * from links order by update_times limit 5")
	defer rows.Close()
	if err != nil {
		//panic(err.Error())
		return nil, err
	}

	var links []*Link
	for rows.Next() {
		var p *Link = new(Link)
		rows.Scan(&p.ShareID, &p.FkUK, &p.IsValid, &p.UpdateTimes, &p.FkSharedFilesId, &p.ShortUrl)
		links = append(links, p)
		p.AddUpdateTimes()
	}
	return links, nil

}

func (l *Link) AddUpdateTimes() error {
	_, err := db.Exec("update links set update_times=update_times+1 where shareid=? and fk_uk=?", &l.ShareID, &l.FkUK)
	if err != nil {
		return err
	}
	return nil
}

func (l *Link) Delete() error {
	_, err := db.Exec("delete from links  where shareid=? and fk_uk=? ", &l.ShareID, &l.FkUK)
	if err != nil {
		return err
	}
	return nil
}

func (l *Link) FindByFkSharedFileId() ([]*Link, error) {
	rows, err := db.Query("select * from links where fk_sharedfiles_id=? ", l.FkSharedFilesId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var links []*Link
	for rows.Next() {
		var p *Link = new(Link)
		rows.Scan(&p.ShareID, &p.FkUK, &p.IsValid, &p.UpdateTimes, &p.FkSharedFilesId, &p.ShortUrl)
		links = append(links, p)

	}
	return links, nil
}

func (l *Link) BatchFindByFkSharedFileId(id []int) ([]*Link, error) {
	if len(id) <= 0 {
		return nil, errors.New("len of id is zero")
	}
	var param string
	for _, v := range id {
		param = param + strconv.Itoa(v) + ","
	}
	param = strings.TrimSuffix(param, ",")
	rows, err := db.Query("select * from links where fk_sharedfiles_id in (?) ", param)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var links []*Link
	for rows.Next() {
		var p *Link = new(Link)
		rows.Scan(&p.ShareID, &p.FkUK, &p.IsValid, &p.UpdateTimes, &p.FkSharedFilesId, &p.ShortUrl)
		links = append(links, p)

	}
	return links, nil
}

func SelectLinksAndSharerBySharedfileId(id int) ([]*LinkAndSharer, error) {
	rows, err := db.Query("SELECT * FROM walksmile.links left join walksmile.sharers on links.fk_uk=sharers.uk where links.fk_sharedfiles_id= ? limit 50;", id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var linksAndSharer []*LinkAndSharer
	for rows.Next() {
		var p *LinkAndSharer = new(LinkAndSharer)
		rows.Scan(&p.Link.ShareID, &p.Link.FkUK, &p.Link.IsValid, &p.Link.UpdateTimes, &p.Link.FkSharedFilesId, &p.Link.ShortUrl,
			&p.Sharer.UK, &p.Sharer.UName, &p.Sharer.AvatarUrl, &p.Sharer.PubshareCount, &p.Sharer.FansCount, &p.Sharer.FollowCount, &p.Sharer.UpdateTime, &p.Sharer.UpdateTimes)
		linksAndSharer = append(linksAndSharer, p)

	}
	return linksAndSharer, nil

}

func SelectLinkByKey(l *Link) (*Link, error) {
	link := &Link{}
	err := db.QueryRow("select * from links where shareid = ? and fk_uk= ?", l.ShareID, l.FkUK).Scan(
		&link.ShareID, &link.FkUK, &link.IsValid, &link.UpdateTimes, &link.FkSharedFilesId, &link.ShortUrl)
	if err != nil {
		return nil, err
	}

	return link, nil
}
