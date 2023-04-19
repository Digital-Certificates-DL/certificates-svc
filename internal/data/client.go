package data

type ClientQ interface {
	New() ClientQ
	Get() (*Client, error)
	Insert(data *Client) (int64, error)
	GetByID(hash string) (*Client, error)
	GetByName(name string) (*Client, error)
	Update(data *Client) error
}

type Client struct {
	ID    int64  `db:"id" structs:"-"`
	Name  string `db:"name" structs:"name"`
	Token []byte `db:"token" structs:"token"`
	Code  string `db:"code" structs:"code"`
}
