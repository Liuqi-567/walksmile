package models

type Hotwords struct {
	ID         int
	Content    string
	UpdateTime string
}

func (this *Hotwords) Update() error {
	_, err := db.Exec("update hotwords set content =?  where id=?", this.Content, this.ID)
	if err != nil {
		return err
	}
	return nil
}

func (this *Hotwords) SelectByID() (*Hotwords, error) {
	var word *Hotwords = new(Hotwords)
	err := db.QueryRow("select * from hotwords where id  = ?", this.ID).Scan(&word.ID, &word.Content, &word.UpdateTime)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (this *Hotwords) SelectByLastID() (*Hotwords, error) {
	var word *Hotwords = new(Hotwords)
	err := db.QueryRow("select * from  hotwords order by id desc limit 1 ").Scan(&word.ID, &word.Content, &word.UpdateTime)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (this *Hotwords) Insert() error {
	stmt, err := db.Prepare("insert into hotwords (content) values (?) ")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(&this.Content)
	if err != nil {
		return err
	}
	return nil
}

func (this *Hotwords) FindList() ([]*Hotwords, error) {
	rows, err := db.Query("select * from  hotwords order by id desc")
	if err != nil {
		return nil, err
	}
	var pics []*Hotwords
	for rows.Next() {
		pic := &Hotwords{}
		rows.Scan(&pic.ID, &pic.Content, &pic.UpdateTime)
		if pic.ID == 1 || pic.ID == 2 {
			continue
		}
		pics = append(pics, pic)
	}
	return pics, nil
}
