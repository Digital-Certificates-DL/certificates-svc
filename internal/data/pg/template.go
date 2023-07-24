package pg

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
)

const templateTableName = "template"

func NewTemplateQ(db *pgdb.DB) data.TemplateQ {
	return &TemplateQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", clientTableName)),
	}
}

type TemplateQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *TemplateQ) New() data.TemplateQ {
	return NewTemplateQ(q.db.Clone())
}

func (q *TemplateQ) Get() (*data.Template, error) {
	var result data.Template
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get template")
	}

	return &result, nil
}

func (q *TemplateQ) Update(client *data.Template) error {
	clauses := structs.Map(client)
	if err := q.db.Exec(sq.Update(templateTableName).SetMap(clauses).Where(sq.Eq{"id": client.ID})); err != nil {
		return errors.Wrap(err, "failed to update template")
	}

	return nil
}

func (q *TemplateQ) Select(id int64) ([]data.Template, error) {
	var result []data.Template

	err := q.db.Select(&result, sq.Select("*").From(templateTableName).Where(sq.Eq{"user_id": id}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to select templates")
	}

	return result, nil
}

func (q *TemplateQ) Insert(value *data.Template) error {
	clauses := structs.Map(value)
	var id int64

	if err := q.db.Get(&id, sq.Insert(templateTableName).SetMap(clauses).Suffix("returning id")); err != nil {
		return errors.Wrap(err, "failed to insert template ")
	}

	return nil
}

func (q *TemplateQ) GetByUserID(id string) (*data.Template, error) {
	var result data.Template
	err := q.db.Get(&result, sq.Select("*").From(templateTableName).Where(sq.Eq{"user_id": id}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get template")
	}

	return &result, nil
}

func (q *TemplateQ) GetByName(name string, clientID int64) (*data.Template, error) {
	var result data.Template
	err := q.db.Get(&result, sq.Select("*").From(templateTableName).Where(sq.Eq{"name": name, "user_id": clientID}))
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get template")
	}
	return &result, nil
}

func (q *TemplateQ) Page(pageParams pgdb.OffsetPageParams) data.TemplateQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *TemplateQ) FilterByUser(id int64) data.TemplateQ {
	q.sql = q.sql.Where(sq.Eq{"user_id": id})
	return q
}
