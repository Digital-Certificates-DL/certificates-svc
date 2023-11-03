package data

type TemplateQ interface {
	New() TemplateQ
	Get() (*Template, error)
	Insert(data *Template) error
	Update(data *Template) error
	Select() ([]Template, error)
	FilterByUser(ids int64) TemplateQ
	FilterByName(name string) TemplateQ
	FilterByShortName(name string) TemplateQ
	FilterByID(ids int64) TemplateQ
}

type Template struct {
	ID        int64  `db:"id" structs:"-"`
	UserID    int64  `db:"user_id" structs:"user_id"`
	Name      string `db:"name" structs:"name"`
	ShortName string `db:"short_name" structs:"short_name"`
	Template  []byte `db:"template" structs:"template"`
	ImgBytes  []byte `db:"img_bytes" structs:"img_bytes"` //todo  make better
}
