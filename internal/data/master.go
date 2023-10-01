package data

type MasterQ interface {
	New() MasterQ
	ClientQ() ClientQ
	TemplateQ() TemplateQ
	Transaction(func(data interface{}) error, interface{}) error
}
