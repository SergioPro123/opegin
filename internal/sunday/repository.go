package sunday

import "devopegin/internal/domain"

type IRepository interface {
	AddExtraHour(extraHours domain.ExtraHour)
	GetExtraHour(id int) *domain.ExtraHour

	AddPosition(position domain.Position)
	GetPosition(name string) *domain.Position
}

type repository struct {
	ExtraHours []domain.ExtraHour
	Positions  []domain.Position
}

func NewRepository() IRepository {
	return &repository{}
}

func (r *repository) AddExtraHour(extraHours domain.ExtraHour) {
	r.ExtraHours = append(r.ExtraHours, extraHours)
}
func (r *repository) GetExtraHour(id int) *domain.ExtraHour {
	for i := 0; i < len(r.ExtraHours); i++ {
		if r.ExtraHours[i].ID == id {
			return &r.ExtraHours[i]
		}
	}
	return nil
}

func (r *repository) AddPosition(position domain.Position) {
	if r.GetPosition(position.Name) != nil {
		return
	}
	r.Positions = append(r.Positions, position)
}
func (r *repository) GetPosition(name string) *domain.Position {
	for i := 0; i < len(r.Positions); i++ {
		if r.Positions[i].Name == name {
			return &r.Positions[i]
		}
	}
	return nil
}
