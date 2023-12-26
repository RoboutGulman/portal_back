package internalapi

import (
	"context"
	"portal_back/role/impl/domain"
)

type RoleRequestService interface {
	GetAllRoles(context context.Context) ([]domain.Role, error)
	GetUserRoles(context context.Context, userId int) ([]domain.Role, error)
	AssignRoleToUser(context context.Context, userId int) error
	RemoveRoleFromUser(context context.Context, userId int) error
}
