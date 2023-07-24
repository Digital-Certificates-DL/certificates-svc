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

const clientTableName = "users"

const (
	idField    = "id"
	nameField  = "name"
	tokenField = "token"
	codeField  = "code"
)

func NewClientQ(db *pgdb.DB) data.ClientQ {
	return &ClientQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", clientTableName)),
		upd: sq.Update("b.*"),
	}
}

type ClientQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func (q *ClientQ) New() data.ClientQ {
	return NewClientQ(q.db.Clone())
}

func (q *ClientQ) Get() (*data.Client, error) {
	var result data.Client
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client")
	}
	return &result, nil
}

func (q *ClientQ) Update(client *data.Client) error {
	clauses := structs.Map(client)
	if err := q.db.Exec(q.upd.SetMap(clauses)); err != nil {
		return errors.Wrap(err, "failed to update client")
	}
	return nil
}

func (q *ClientQ) Insert(value *data.Client) error {
	clauses := structs.Map(value)
	var id int64

	if err := q.db.Get(&id, sq.Insert(clientTableName).SetMap(clauses).Suffix("returning id")); err != nil {
		return errors.Wrap(err, "failed to insert client")
	}

	return nil
}

func (q *ClientQ) FilterByID(id int64) data.ClientQ {
	q.sql = q.sql.Where(sq.Eq{idField: id})
	q.upd = q.upd.Where(sq.Eq{idField: id})
	return q
}

func (q *ClientQ) FilterByName(name string) data.ClientQ {
	q.sql = q.sql.Where(sq.Eq{nameField: name})
	q.upd = q.upd.Where(sq.Eq{nameField: name})
	return q
}

func (q *ClientQ) Page(pageParams pgdb.OffsetPageParams) data.ClientQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}
