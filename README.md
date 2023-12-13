Для запуска:
1. `go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest`
2. `oapi-codegen -generate gorilla,types authentication/api/frontend/frontendapi.yaml > authentication/api/frontend/frontendapi.gen.go`
3. `oapi-codegen -generate gorilla,types documentation/api/frontend/frontendapi.yaml > documentation/api/frontend/frontendapi.gen.go`
4. `oapi-codegen -generate gorilla,types company/api/frontend/frontendapi.yaml > company/api/frontend/frontendapi.gen.go`
4. `go mod tidy -v`
5. `go run main.go`

Configure GoLand run config with the following env vars: 
- DB_USER
- DB_PASSWORD
- DB_NAME
- DB_HOST
- BACKEND_PORT
- DB_DOCUMENTATION_NAME

![img.png](img/envVars.png)

Для работы с бд
1. установить postgres
2. создать две таблицы

```
CREATE TYPE roleTypes AS ENUM('VIEWER', 'EDITOR', 'ADMIN', 'OWNER');

CREATE TABLE role (
	id SERIAL PRIMARY KEY NOT NULL,
	title VARCHAR(256) NOT NULL,
	description VARCHAR(256) NOT NULL,
	roleType roleTypes NOT NULL
);

CREATE TABLE company (
	id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR(256) NOT NULL
);

CREATE TABLE auth_user (
	id SERIAL PRIMARY KEY NOT NULL,
	password VARCHAR(256) NOT NULL,
	salt VARCHAR(256) NOT NULL,
	email VARCHAR(256) NOT NULL
);

CREATE TABLE tokens (
	id SERIAL PRIMARY KEY NOT NULL,
	userId INT NOT NULL,
	token VARCHAR(256) NOT NULL,
	FOREIGN KEY(userId) REFERENCES auth_user(id)
);

CREATE TABLE employeeAccount (
	id SERIAL PRIMARY KEY NOT NULL,
	userId INT NOT NULL,
	companyId INT NOT NULL,
	firstName VARCHAR(256) NOT NULL,
	secondName VARCHAR(256) NOT NULL,
	surname VARCHAR(256) NOT NULL,
	telephoneNumber VARCHAR(256) NOT NULL,
	avatarUrl VARCHAR(256) NOT NULL,
	dateOfBirth DATE NOT NULL,
	job VARCHAR(256) NOT NULL,
	UNIQUE (userId, companyId),
	FOREIGN KEY(userId) REFERENCES auth_user(id),
	FOREIGN KEY(companyId) REFERENCES company(id)
);

CREATE TABLE department (
	id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR(256) NOT NULL,
	parentDepartmentId INT,
	companyId INT NOT NULL,
	supervisorId INT NOT NULL,
	FOREIGN KEY(parentDepartmentId) REFERENCES department(id),
	FOREIGN KEY(companyId) REFERENCES company(id),
	FOREIGN KEY(supervisorId) REFERENCES employeeAccount(id)
);

CREATE TABLE employee_roles (
	accountId INT NOT NULL,
	roleId INT NOT NULL,
	FOREIGN KEY(accountId) REFERENCES employeeAccount(id),
	FOREIGN KEY(roleId) REFERENCES role(id),
	PRIMARY KEY(accountId, roleId)
);

CREATE TABLE employee_department (
	accountId INT NOT NULL,
	departmentId INT NOT NULL,
	FOREIGN KEY(accountId) REFERENCES employeeAccount(id),
	FOREIGN KEY(departmentId) REFERENCES department(id),
	PRIMARY KEY (accountId, departmentId)
);
```

3. В отдельной бд для документов
```
CREATE TABLE sections
(
    id serial primary key,
    title character varying(256) NOT NULL,
    thumbnail_url character varying(256) NOT NULL,
    is_favorite boolean NOT NULL,
    company_id integer NOT NULL
);

create table user_sections_prefs(
	user_id integer NOT NULL,
	section_id integer NOT NULL,
	FOREIGN KEY (section_id) REFERENCES sections(id),
	PRIMARY KEY (user_id, section_id)
)
```