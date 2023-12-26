package transport

import (
	"encoding/json"
	"net/http"
	"portal_back/core/utils"
	frontendapi "portal_back/role/api/frontend"
	"portal_back/role/impl/app/role"
	"portal_back/role/impl/domain"
)

func NewServer(service role.RoleService) frontendapi.ServerInterface {
	return &roleServer{roleService: service}
}

type roleServer struct {
	roleService role.RoleService
}

// GetRoles implements frontendapi.ServerInterface.
func (server *roleServer) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := server.roleService.GetAllRoles(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rolesJson := utils.Map(roles, func(r domain.Role) frontendapi.Role {
		return frontendapi.Role{
			Description: r.Description,
			Id:          r.Id,
			RoleType:    r.RoleType,
			Title:       r.Title,
		}
	})
	response, err := json.Marshal(frontendapi.RolesResponse{Roles: rolesJson})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
