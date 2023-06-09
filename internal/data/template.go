package data

type TemplateQ interface {
	New() TemplateQ
	Get() (*Template, error)
	Insert(data *Template) (int64, error)
	GetByUserID(hash string) (*Template, error)
	GetByName(name string) (*Template, error)
	Update(data *Template) error
}

type Template struct {
	ID       int64  `db:"id" structs:"-"`
	UserID   int64  `db:"user_id" structs:"user_id"`
	Name     string `db:"name" structs:"name"`
	Template []byte `db:"template" structs:"template"`
	ImgBytes []byte `db:"img_bytes" structs:"img_bytes"`
}
