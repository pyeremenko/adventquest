package model

type InputModelFactory interface {
	Create() interface{}
}

type CreateTaskInput struct {
	Day  int    `json:"day"`
	Link string `json:"link"`
}

type CreateTaskInputFactory struct{}

func (f *CreateTaskInputFactory) Create() interface{} {
	return &CreateTaskInputFactory{}
}
