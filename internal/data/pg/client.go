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

func NewClientQ(db *pgdb.DB) data.ClientQ {
	return &ClientQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", clientTableName)),
		//upd: sq.Update(),
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

//func (q *ClientQ) WhereId(id int64) *ClientQ {
//	q.sql = q.sql.Where()
//	//q.upd = q.upd.Where()
//	return q
//}

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
	err := q.db.Exec(sq.Update(clientTableName).SetMap(clauses).Where(sq.Eq{"id": client.ID}))
	if err != nil {
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

func (q *ClientQ) GetByID(id string) (*data.Client, error) {
	var result data.Client
	err := q.db.Get(&result, sq.Select("*").From(clientTableName).Where(sq.Eq{"id": id}))
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get client")
	}

	return &result, nil
}

func (q *ClientQ) GetByName(name string) (*data.Client, error) {
	var result data.Client
	err := q.db.Get(&result, sq.Select("*").From(clientTableName).Where(sq.Eq{"name": name}))
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get client")
	}

	return &result, nil
}

func (q *ClientQ) Page(pageParams pgdb.OffsetPageParams) data.ClientQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}
