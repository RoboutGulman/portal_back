package di

import (
	"context"
	"fmt"
	"net/http"
	"portal_back/authentication/cmd"
	frontendapi "portal_back/role/api/frontend"
	"portal_back/role/impl/app/role"
	"portal_back/role/impl/infrasructure/sql"
	"portal_back/role/impl/infrasructure/transport"

	"github.com/jackc/pgx/v5"
)

func InitRolesModule(config cmd.Config) (role.RoleService, *pgx.Conn, error) {

	dbStringConnection := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	conn, _ := pgx.Connect(context.Background(), dbStringConnection)

	roleRepository := sql.NewRepository(conn)
	roleService := role.NewService(roleRepository)
	roleServer := transport.NewServer(roleService)
	http.Handle("/role/", frontendapi.Handler(roleServer))
	return roleService, conn, nil
}
