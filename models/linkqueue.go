package models

type LinkQueue struct {
	QueueID  int
	UK       int
	Shareid  int
	ShortUrl string
}

func (l *LinkQueue) Insert() error {
	stmt, err := db.Prepare("insert into linkqueue (uk,shareid,shorturl) values (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(&l.UK, &l.Shareid, &l.ShortUrl)
	if err != nil {
		return err
	}
	return nil
}

func (l *LinkQueue) FindWithLimit(limit int) ([]*LinkQueue, error) {
	rows, err := db.Query("select * from linkqueue order by queue_id desc limit ? ", limit)
	if err != nil {
		return nil, err
	}
	var links []*LinkQueue
	for rows.Next() {
		link := &LinkQueue{}
		rows.Scan(&link.QueueID, &link.UK, &link.Shareid, &link.ShortUrl)
		links = append(links, link)
	}
	return links, nil
}

func (l *LinkQueue) Delete() error {
	_, err := db.Exec("delete from linkqueue  where uk=?  and shareid=?", &l.UK, &l.Shareid)
	if err != nil {
		return err
	}
	return nil
}
