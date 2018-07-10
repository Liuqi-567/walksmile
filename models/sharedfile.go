package models

type SharedFile struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Hot        int    `json:"hot"`
	Category   int    `json:"category"`
	LinkCount  int    `json:"link_count"`
	UpdateTime string `json:"update_time"`
	CreateTime string `json:"create_time"`
}

func (s *SharedFile) Insert() error {
	res, err := db.Exec("insert into sharedfiles (title,category) values (?,?)", s.Title, s.Category)
	if err != nil {
		//panic(err.Error())
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		//panic(err.Error())
		return err
	}

	s.ID = int(id)

	return nil

}

func (s *SharedFile) FindByName() (*SharedFile, error) {

	err := db.QueryRow("select * from sharedfiles where title = ?", s.Title).Scan(
		&s.ID, &s.Title, &s.Hot, &s.Category, &s.LinkCount, &s.UpdateTime, &s.CreateTime)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SharedFile) UpdateHot() error {
	_, err := db.Exec("update sharedfiles set hot=hot+1 where id=?", s.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SharedFile) AddLinkCount() error {
	_, err := db.Exec("update sharedfiles set link_count=link_count+1 where id=?", s.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SharedFile) FindByID() (*SharedFile, error) {
	err := db.QueryRow("select * from sharedfiles where id  = ?", s.ID).Scan(
		&s.ID, &s.Title, &s.Hot, &s.Category, &s.LinkCount, &s.UpdateTime, &s.CreateTime)

	if err != nil {
		return nil, err
	}
	return s, nil
}

func FindLast24() ([]*SharedFile, error) {
	rows, err := db.Query("select * from sharedfiles where create_time  >=  NOW() - interval 24 hour order by create_time desc;")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var sflist []*SharedFile
	for rows.Next() {
		sharedfile := &SharedFile{}
		rows.Scan(&sharedfile.ID, &sharedfile.Title, &sharedfile.Hot, &sharedfile.Category, &sharedfile.LinkCount, &sharedfile.UpdateTime, &sharedfile.CreateTime)
		sflist = append(sflist, sharedfile)

	}
	return sflist, nil

}

func FindSFNearID(id int) ([]*SharedFile, error) {
	rows, err := db.Query("select * from sharedfiles where id  > ? limit 20;", id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	sflist := make([]*SharedFile, 0, 20)
	for rows.Next() {
		sharedfile := &SharedFile{}
		rows.Scan(&sharedfile.ID, &sharedfile.Title, &sharedfile.Hot, &sharedfile.Category, &sharedfile.LinkCount, &sharedfile.UpdateTime, &sharedfile.CreateTime)
		sflist = append(sflist, sharedfile)

	}
	return sflist, nil

}
