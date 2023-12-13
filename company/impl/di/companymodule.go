package di

import (
	"github.com/jackc/pgx/v5"
	"net/http"
	"portal_back/authentication/api/internalapi"
	"portal_back/company/api/frontend"
	"portal_back/company/impl/app/employeeaccount"
	"portal_back/company/impl/infrastructure/sql"
	"portal_back/company/impl/infrastructure/transport"
	rolesapi "portal_back/roles/api/internalapi"
)

func InitCompanyModule(authApi internalapi.AuthRequestService, rolesApi rolesapi.RolesRequestService, conn *pgx.Conn) {

	repo := sql.NewEmployeeAccountRepository(conn)
	accountService := employeeaccount.NewService(repo, authApi)

	server := transport.NewServer(accountService, rolesApi)

	http.Handle("/", frontendapi.Handler(server))
}
