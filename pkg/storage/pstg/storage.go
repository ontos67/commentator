package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// База данных.
type DB struct {
	pool  *pgxpool.Pool
	CChan chan Comment
}

// Комментарий
type Comment struct {
	ID         int    // номер комментария
	User_id    int    // ID автора комментария
	Text       string // содержание комментария
	PubTime    int64  // время публикации, Unixtime
	ParentType string // тип родителя (A - статья (комментарий на саму статью), С - комментарий (отчеточка на комментарий) )
	ParentID   int    // ID родителя (или статьи или комментария)
}

func New(adr string) (*DB, error) {
	dbuserpass := os.Getenv("agrigatordb") //"postgres://postgres:" + "password"
	connstr := dbuserpass + adr
	if connstr == "" {
		return nil, errors.New("не указано подключение к БД")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)

	if err != nil {
		return nil, err
	}
	db := DB{
		pool:  pool,
		CChan: make(chan Comment),
	}

	return &db, nil
}

// Сохранить комментарций в БД
func (db *DB) SaveComment(c Comment) (int, error) {
	//row, err:= db.pool.Exec(context.Background(), `
	row := db.pool.QueryRow(context.Background(), `
		INSERT INTO comments (user_id, comment_text, pub_time, parent_type, parent_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		c.User_id,
		c.Text,
		c.PubTime,
		c.ParentType,
		c.ParentID,
	)
	err := row.Scan(&c.ID)
	if err != nil {
		return -1, err
	}
	return c.ID, nil
}

// DeleteComment удаляет комментарций из БД (например, из-за нарушения правил по содержанию)
func (db *DB) DeleteComment(id int) error {
	r, err := db.pool.Exec(context.Background(), `
		DELETE FROM comments
		WHERE id =$1;
		`,
		id,
	)
	if r.RowsAffected() != 1 {
		return fmt.Errorf("запись с таким id не найдена")
	}
	if err != nil {
		return err
	}

	return nil
}

// CommentList возвращает список комментариев к объекту из БД (ID, ID автора, время публикации, тип родителя, ID родителя).
// предполагается древовидная структура комментариев, метод выдает результат на 1 уровень вниз
// для каждого родителя. если у стратья 3 комментария, на 2 из которых даны ответочки,
// для получения всех комментариев необходимбо будет сделать 3 запроса: 1 - получение комментариев
// на саму статью, и по одному запросу - на получение комментариев на комменатарий
func (db *DB) CommentList(pType string, pID int) ([]Comment, error) {

	rows, err := db.pool.Query(context.Background(), `
	SELECT id, user_id, comment_text, pub_time, parent_type, parent_id FROM comments
	WHERE parent_type=$1 AND parent_id=$2
	ORDER BY pub_time DESC
	`,
		pType,
		pID,
	)
	if err != nil {
		return nil, err
	}
	var clist []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(
			&c.ID,
			&c.User_id,
			&c.Text,
			&c.PubTime,
			&c.ParentType,
			&c.ParentID,
		)
		if err != nil {
			return nil, err
		}
		clist = append(clist, c)
	}
	return clist, rows.Err()
}

// LastComments возвращает n последовательнодобавленных комментариев, начиная
// с ID = startID. Если startID==-1, возвращается 2n последних комментариев
func (db *DB) CommentsListCont(startID, n int) ([]Comment, error) {
	if n == 0 {
		n = 100
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, user_id, comment_text, pub_time, parent_type, parent_id FROM comments
	WHERE id BETWEEN $1 AND $2
	ORDER BY id DESC
	`,
		startID,
		startID+n,
	)
	if err != nil {
		return nil, err
	}
	var clist []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(
			&c.ID,
			&c.User_id,
			&c.Text,
			&c.PubTime,
			&c.ParentType,
			&c.ParentID,
		)
		if err != nil {
			return nil, err
		}
		clist = append(clist, c)
	}
	return clist, rows.Err()
}
