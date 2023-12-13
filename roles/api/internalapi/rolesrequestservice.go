package internalapi

import (
	"context"
	"portal_back/roles/api/internalapi/model"
)

type RolesRequestService interface {
	IsUserHasRole(ctx context.Context, accountId int, roleType model.RoleType) (bool, error)
}
