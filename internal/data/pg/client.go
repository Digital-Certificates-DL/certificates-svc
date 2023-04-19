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

const clientTableName = "users"

func NewClientQ(db *pgdb.DB) data.ClientQ {
	return &ClientQ{
		db:  db.Clone(),
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", clientTableName)),
	}
}

type ClientQ struct {
	mx  sync.Mutex
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *ClientQ) New() data.ClientQ {
	q.mx.Lock()
	defer q.mx.Unlock()
	return NewClientQ(q.db)
}

func (q *ClientQ) Get() (*data.Client, error) {
	q.mx.Lock()
	defer q.mx.Unlock()
	var result data.Client
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *ClientQ) Update(client *data.Client) error {
	clauses := structs.Map(client)
	stmt := sq.Update(clientTableName).SetMap(clauses).Where(sq.Eq{"id": client.ID})
	err := q.db.Exec(stmt)

	return err
}

func (q *ClientQ) Insert(value *data.Client) (int64, error) {
	clauses := structs.Map(value)
	var id int64

	stmt := sq.Insert(clientTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)

	return id, err
}

func (q *ClientQ) GetByID(id string) (*data.Client, error) {

	var result data.Client
	err := q.db.Get(&result, sq.Select("*").From(clientTableName).Where(sq.Eq{"id": id}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *ClientQ) GetByName(name string) (*data.Client, error) {
	var result data.Client
	err := q.db.Get(&result, sq.Select("*").From(clientTableName).Where(sq.Eq{"name": name}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *ClientQ) Page(pageParams pgdb.OffsetPageParams) data.ClientQ {
	q.mx.Lock()
	defer q.mx.Unlock()
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}
