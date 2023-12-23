package di

import (
	"github.com/jackc/pgx/v5"
	"net/http"
	"portal_back/authentication/api/internalapi"
	"portal_back/company/api/frontend"
	"portal_back/company/impl/app/department"
	"portal_back/company/impl/app/employeeaccount"
	"portal_back/company/impl/infrastructure/sql"
	"portal_back/company/impl/infrastructure/transport"
	rolesapi "portal_back/roles/api/internalapi"
)

func InitCompanyModule(authApi internalapi.AuthRequestService, userApi internalapi.UserRequestService, rolesApi rolesapi.RolesRequestService, conn *pgx.Conn) {

	accountRepo := sql.NewEmployeeAccountRepository(conn)
	accountService := employeeaccount.NewService(accountRepo, userApi)

	departmentRepo := sql.NewDepartmentRepository(conn)
	departmentService := department.NewService(departmentRepo, accountService)

	server := transport.NewServer(accountService, departmentService, rolesApi, authApi)

	http.Handle("/", frontendapi.Handler(server))
}
