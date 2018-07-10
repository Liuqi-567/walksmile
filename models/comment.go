package models

type Comment struct {
	Comment_id  int
	Name        string
	Email       string
	Content     string
	Kind        int
	Fatherid    int
	Create_time string
}

func (c *Comment) Insert() error {
	stmt, err := db.Prepare("insert into comment (name,email,content,kind, fatherid) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(&c.Name, &c.Email, &c.Content, &c.Kind, &c.Fatherid)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) FindByFatherID() ([]*Comment, error) {
	rows, err := db.Query("select * from comment where fatherid = ? and kind = ?  order by create_time desc ", c.Fatherid, c.Kind)
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		rows.Scan(&comment.Comment_id, &comment.Name, &comment.Email, &comment.Content, &comment.Kind, &comment.Fatherid, &comment.Create_time)
		comments = append(comments, comment)
	}
	return comments, nil
}
