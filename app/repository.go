package app

import (
	"context"
	"github.com/EmreZURNACI/apistack.git/domain"
)

type Repository interface {
	GetActors(ctx context.Context, search string, offset, limit int, orderBy bool) ([]domain.Actor, error)
	CreateActor(ctx context.Context, firstName, lastName string) (int64, error)
	DeleteActor(ctx context.Context, id string) error
	GetActor(ctx context.Context, id string) (*domain.Actor, error)
	UpdateActor(ctx context.Context, id, firstname, lastname string) error
}
