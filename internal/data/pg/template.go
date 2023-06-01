package pg

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"sync"
)

const tempalteTableName = "template"

func NewTemplateQ(db *pgdb.DB) data.TemplateQ {
	return &TemplateQ{
		db:  db.Clone(),
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", clientTableName)),
	}
}

type TemplateQ struct {
	mx  sync.Mutex
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *TemplateQ) New() data.TemplateQ {
	q.mx.Lock()
	defer q.mx.Unlock()
	return NewTemplateQ(q.db)
}

func (q *TemplateQ) Get() (*data.Template, error) {
	q.mx.Lock()
	defer q.mx.Unlock()
	var result data.Template
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *TemplateQ) Update(client *data.Template) error {
	clauses := structs.Map(client)
	stmt := sq.Update(tempalteTableName).SetMap(clauses).Where(sq.Eq{"id": client.ID})
	err := q.db.Exec(stmt)

	return err
}

func (q *TemplateQ) FilterByUser(id int64) data.TemplateQ {
	q.sql = q.sql.Where(sq.Eq{"user_id": id})
	return q
}

func (q *TemplateQ) Select(id int64) ([]data.Template, error) {
	var result []data.Template

	stmt := sq.Select("*").From(tempalteTableName).Where(sq.Eq{"user_id": id})
	err := q.db.Select(&result, stmt)

	return result, err
}

func (q *TemplateQ) Insert(value *data.Template) (int64, error) {
	clauses := structs.Map(value)
	var id int64

	stmt := sq.Insert(tempalteTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)

	return id, err
}

func (q *TemplateQ) GetByUserID(id string) (*data.Template, error) {
	var result data.Template
	err := q.db.Get(&result, sq.Select("*").From(tempalteTableName).Where(sq.Eq{"user_id": id}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *TemplateQ) GetByName(name string, clientID int64) (*data.Template, error) {
	var result data.Template
	err := q.db.Get(&result, sq.Select("*").From(tempalteTableName).Where(sq.Eq{"name": name, "user_id": clientID}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *TemplateQ) Page(pageParams pgdb.OffsetPageParams) data.TemplateQ {
	q.mx.Lock()
	defer q.mx.Unlock()
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}
