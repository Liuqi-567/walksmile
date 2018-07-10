package models

import (
	_ "database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Sharer struct {
	UK            int
	UName         string
	AvatarUrl     string
	PubshareCount int
	FansCount     int
	FollowCount   int
	UpdateTime    string
	UpdateTimes   int
}

func (s *Sharer) Insert() error {
	_, err := db.Exec("insert into sharers (uk,uname,avatar_url,pubshare_count,fans_count,follow_count) values (?,?,?,?,?,?) ", &s.UK, &s.UName, &s.AvatarUrl, &s.PubshareCount, &s.FansCount, &s.FollowCount)
	if err != nil {
		return err
	}
	return nil

}

func (s *Sharer) AddUpdateTimes() error {
	_, err := db.Exec("update sharers set update_times=update_times+1 where uk=?", s.UK)
	if err != nil {
		return err
	}
	return nil

}

func (s *Sharer) Update() error {
	_, err := db.Exec("update sharers set uname=?,avatar_url=?,pubshare_count=?,fans_count=?,follow_count=? where uk=? ", &s.UName, &s.AvatarUrl, &s.PubshareCount, &s.FansCount, &s.FollowCount, &s.UK)
	if err != nil {
		return err
	}
	return nil

}
func (p *Sharer) FindByUK(s int) (*Sharer, error) {
	err := db.QueryRow("select * from sharers where uk=?", s).Scan(&p.UK, &p.UName, &p.AvatarUrl, &p.PubshareCount, &p.FollowCount, &p.FansCount, &p.UpdateTime, &p.UpdateTimes)
	if err != nil {
		return nil, err
	}
	return p, nil

}
