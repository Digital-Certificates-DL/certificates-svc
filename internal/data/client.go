package data

type ClientQ interface {
	New() ClientQ
	Get() (*Client, error)
	Insert(data *Client) error
	Update(data *Client) error
	FilterByName(name string) ClientQ
	FilterByID(id int64) ClientQ
}

type Client struct {
	ID    int64  `db:"id" structs:"-"`
	Name  string `db:"name" structs:"name"`
	Token []byte `db:"token" structs:"token"`
	Code  string `db:"code" structs:"code"`
}
