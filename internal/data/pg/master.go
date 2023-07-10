package pg

import (
	"database/sql"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
)

type masterQ struct {
	db *pgdb.DB
}

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{
		db: db,
	}
}

func (q *masterQ) New() data.MasterQ {
	return NewMasterQ(q.db.Clone())
}

func (q *masterQ) ClientQ() data.ClientQ {
	return NewClientQ(q.db.Clone())
}

func (q *masterQ) TemplateQ() data.TemplateQ {
	return NewTemplateQ(q.db.Clone())
}

func (q *masterQ) Transaction(fn func(data interface{}) error, data interface{}) error {
	return q.db.TransactionWithOptions(&sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}, func() error {
		return fn(data)
	})
}
