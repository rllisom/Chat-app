package user

import (
	"errors"
	"strings"
)

type UserService struct {
	userRespo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
	return &UserService{userRespo: userRepo}
}

func (s *UserService) CreateUser(username, email string) (*User, error) {

	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" {
		return nil, errors.New("El username no puede estar vacío") 
	}

	if len(username) < 3 {
		return nil,errors.New("El username debe contener más de 3 caracteres")
	}
	
	if email == "" {
		return nil, errors.New("El email no puede estar vacío")
	}

	if !strings.Contains(email,"@"){
		return nil, errors.New("El email no es válido")
	}

	_,err := s.userRespo.FindByUsername(username)

	if err != nil {
		return nil, errors.New("El username ya existe")
	}

	newUser := NewUser(username,email)
	return s.userRespo.Create(newUser)
}

func (s *UserService) GetUserByID(id string) (*User,error){

	if id == "" {
		return nil, errors.New("El id no puede estar vacío")
	}

	return s.userRespo.FindByID(id)
}

func (s *UserService) GetAllUsers() (*[]User,error){
	return s.userRespo.FindAll()
}