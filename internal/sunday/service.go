package sunday

type IService interface {
}

type service struct {
	IService
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &service{
		repository: repository,
	}
}
