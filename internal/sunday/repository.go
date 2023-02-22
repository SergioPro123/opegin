package sunday

import "devopegin/internal/domain"

type IRepository interface {
	AddExtraHour(extraHours domain.ExtraHour)
	GetExtraHour(id int) *domain.ExtraHour

	AddLocation(extraHours domain.Location)
	GetLocation(name string) *domain.Location
}

type repository struct {
	ExtraHours []domain.ExtraHour
	Locations  []domain.Location
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

func (r *repository) AddLocation(location domain.Location) {
	if r.GetLocation(location.Name) != nil {
		return
	}
	r.Locations = append(r.Locations, location)
}
func (r *repository) GetLocation(name string) *domain.Location {
	for i := 0; i < len(r.Locations); i++ {
		if r.Locations[i].Name == name {
			return &r.Locations[i]
		}
	}
	return nil
}
