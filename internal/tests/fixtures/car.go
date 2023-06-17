package fixtures

import "homework-7/internal/pkg/repository"

type CarBuilder struct {
	instance *repository.Car
}

func Car() *CarBuilder {
	return &CarBuilder{instance: &repository.Car{}}
}

func (b *CarBuilder) Model(v string) *CarBuilder {
	b.instance.Model = v
	return b
}

func (b *CarBuilder) Id(v int64) *CarBuilder {
	b.instance.ID = v
	return b
}

func (b *CarBuilder) UserId(v int64) *CarBuilder {
	b.instance.UserId = v
	return b
}

func (b *CarBuilder) Pointer() *repository.Car {
	return b.instance
}

func (b *CarBuilder) Value() repository.Car {
	return *b.instance
}
