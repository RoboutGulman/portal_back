package userrequest

import "portal_back/authentication/api/internalapi"

func NewService() internalapi.UserRequestService {
	return &service{}
}

type service struct {
}

func (s service) CreateNewUser(email string) error {
	//TODO implement me
	panic("implement me")

	//если такой юзер уже существует, падает ошибка UserAlreadyExists
}

func (s service) GetUserId(email string) (int, error) {
	//TODO implement me
	return 35, nil
}
