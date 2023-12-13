package di

import (
	"portal_back/roles/api/internalapi"
	"portal_back/roles/impl/app/rolesrequest"
)

func InitRolesModule() internalapi.RolesRequestService {
	return rolesrequest.NewService()
}
