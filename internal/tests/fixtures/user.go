package fixtures

import "homework-7/internal/pkg/repository"

type UserBuilder struct {
	instance *repository.User
}

func User() *UserBuilder {
	return &UserBuilder{instance: &repository.User{}}
}

func (b *UserBuilder) Name(v string) *UserBuilder {
	b.instance.Name = v
	return b
}

func (b *UserBuilder) Id(v int64) *UserBuilder {
	b.instance.ID = v
	return b
}

func (b *UserBuilder) Pointer() *repository.User {
	return b.instance
}

func (b *UserBuilder) Value() repository.User {
	return *b.instance
}
