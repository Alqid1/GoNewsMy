package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

type Interface interface {
	NewPost(Post) error
	Posts() ([]Post, error)
	/* UpdatePost(Post) error
	DeletePost(int) error */
}

type DatBase struct {
	db *pgxpool.Pool
}

// потом добавить constr
func New(constr string) (*DatBase, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	s := DatBase{
		db: db,
	}
	return &s, nil
}

// Новая запись
func (d *DatBase) NewPost(p Post) error {
	_, err := d.db.Exec(context.Background(), `
	INSERT INTO posts(title,content,pubtime,link) VALUES ($1,$2,$3,$4);
`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	)
	if err != nil {
		return err
	}
	return nil
}

// Показать n-ое кол-во записей
func (d *DatBase) Posts(n int) ([]Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := d.db.Query(context.Background(), `
	SELECT * FROM posts
	ORDER BY pubtime DESC
	LIMIT $1;
	`, n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

/* func (d *DatBase) UpdatePost(p Post) error {
	_, err := d.db.Exec(context.Background(), `
	UPDATE posts
	SET title='$1',content = '$2',pubtime=$3,link='$4'
	WHERE id = $5;
	`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
		p.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *DatBase) DeletePost(id int) error {
	_, err := d.db.Exec(context.Background(), `
	DELETE FROM posts
	WHERE id = $1;
	`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
} */
