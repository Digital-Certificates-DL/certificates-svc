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

const (
	userIDField = "user_id"
)

func NewTemplateQ(db *pgdb.DB) data.TemplateQ {
	return &TemplateQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", templateTableName)),
		upd: sq.Update(templateTableName),
	}
}

type TemplateQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
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
	if err := q.db.Exec(q.upd.SetMap(clauses)); err != nil {
		return errors.Wrap(err, "failed to update template")
	}

	return nil
}

func (q *TemplateQ) Select() ([]data.Template, error) {
	var result []data.Template

	err := q.db.Select(&result, q.sql)
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

func (q *TemplateQ) Page(pageParams pgdb.OffsetPageParams) data.TemplateQ {
	q.sql = pageParams.ApplyTo(q.sql, idField)

	return q
}

func (q *TemplateQ) FilterByUser(id int64) data.TemplateQ {
	q.sql = q.sql.Where(sq.Eq{userIDField: id})

	return q
}

func (q *TemplateQ) FilterByID(id int64) data.TemplateQ {
	q.sql = q.sql.Where(sq.Eq{idField: id})
	q.upd = q.upd.Where(sq.Eq{idField: id})

	return q
}

func (q *TemplateQ) FilterByName(name string) data.TemplateQ {
	q.sql = q.sql.Where(sq.Eq{nameField: name})
	q.upd = q.upd.Where(sq.Eq{nameField: name})

	return q
}
