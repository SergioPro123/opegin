package sunday

type IRepository interface {
}

type repository struct {
	IRepository
}

func NewRepository() IRepository {
	return &repository{}
}
