package user

import (
	"chat-app/internal/db"
	"testing"
)

func setupTestService(t *testing.T) *UserService {

	if err := db.Connect("mongodb://localhost:27017"); err != nil{
		t.Fatal("No se pudo conectar a MongoDB:",err)
	}

	userRepo := NewUserRepository()
	return NewUserService(userRepo)
}

func TestCreaterUser(t *testing.T) {
	s := setupTestService(t)

	user, err := s.CreateUser("raul", "raul@email.com")

	if err !=nil {
		t.Fatalf("Se esperaba un valor nulo pero se ha obtenido: %v\n",err)
	}

	if user.Username != "raul" {
		t.Errorf("Se ha obtenido lo esperado: %v\n",err)
	}

	if user.Email != "raul@email.com" {
		t.Errorf("Se ha obtenido lo esperado: %v\n",err)
	}

	if user.ID.IsZero() {
		t.Error("el ID no debería estar vacío")
	}
}

func TestCreateUser_UsernameVacio(t *testing.T){
	s := setupTestService(t)

	_,err := s.CreateUser("","raul@email.com")

	if err == nil{
		t.Fatalf("Hay un error al crear el usuario distinto al esperado")
	}
	
	if err.Error() != "El username no puede estar vacío" {
		t.Errorf("Mensaje de error inesperado %v\n",err)
	}

}

func TestCreateUser_EmailInvalido(t *testing.T) {
	s := setupTestService(t)

	_, err := s.CreateUser("ana_go", "esto-no-es-un-email")

	if err == nil {
		t.Fatal("se esperaba error por email inválido")
	}
}

func TestCreateUser_UsernameRepetido(t *testing.T) {
	s := setupTestService(t)


	_, err := s.CreateUser("ana_go", "ana@email.com")
	if err != nil {
		t.Fatalf("primer usuario no debería fallar: %v", err)
	}

	_, err = s.CreateUser("ana_go", "otro@email.com")
	if err == nil {
		t.Fatal("se esperaba error por username duplicado")
	}
}

func TestGetAllUsers(t *testing.T) {
	s := setupTestService(t)

	s.CreateUser("ana_go", "ana@email.com")
	s.CreateUser("luis_dev", "luis@email.com")


	users, err := s.GetAllUsers()
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(*users) != 2 {
		t.Errorf("se esperaban 2 usuarios, obtenidos %d", len(*users))
	}
}