package cmd

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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

func InitCompanyModule(config Config, authApi internalapi.AuthRequestService, userApi internalapi.UserRequestService, rolesApi rolesapi.RolesRequestService) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBName)

	conn, _ := pgx.Connect(context.Background(), connStr)

	accountRepo := sql.NewEmployeeAccountRepository(conn)
	accountService := employeeaccount.NewService(accountRepo, userApi)

	departmentRepo := sql.NewDepartmentRepository(conn)
	departmentService := department.NewService(departmentRepo, accountService)

	server := transport.NewServer(accountService, departmentService, rolesApi, authApi)

	router := mux.NewRouter()
	router.MethodNotAllowedHandler = methodNotAllowedHandler()

	options := frontendapi.GorillaServerOptions{
		BaseRouter: router,
		Middlewares: []frontendapi.MiddlewareFunc{func(handler http.Handler) http.Handler {
			return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				setCorsHeaders(w)
				handler.ServeHTTP(w, r)
			}))
		}},
	}
	r := frontendapi.HandlerWithOptions(server, options)

	http.Handle("/", r)
}

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			setCorsHeaders(w)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-user-id, X-organization-id")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
