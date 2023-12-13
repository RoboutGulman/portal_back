package rolesrequest

import (
	"context"
	"portal_back/roles/api/internalapi"
	"portal_back/roles/api/internalapi/model"
)

func NewService() internalapi.RolesRequestService {
	return &service{}
}

type service struct {
}

func (s service) IsUserHasRole(ctx context.Context, accountId int, roleType model.RoleType) (bool, error) {
	//TODO implement me
	panic("implement me")
}
