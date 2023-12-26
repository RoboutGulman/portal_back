package role

import (
	"context"
	"portal_back/role/impl/domain"
)

type RoleRepository interface {
	GetAllRoles(context context.Context) ([]domain.Role, error)
	GetUserRoles(context context.Context, userId int) ([]domain.Role, error)
	AssignRoleToUser(context context.Context, roleId, userId int) error
	RemoveRoleFromUser(context context.Context, roleId, userId int) error
}
