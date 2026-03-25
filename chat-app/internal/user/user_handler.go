package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}


func (h *UserHandler) CreateUser(c *gin.Context){

	var req CreateUserRequest

	if err:=c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user,err := h.userService.CreateUser(req.Username,req.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)

}

func (h *UserHandler) GetAllUsers(c *gin.Context){

	users,err:= h.userService.GetAllUsers()

	if  err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,users)
}

func (h *UserHandler) GetById(c *gin.Context){
	
	id:= c.Param("id")

	user,err:= h.userService.GetUserByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RegisterRoutes (rg *gin.RouterGroup) {
	baseUrl := rg.Group("/users")
	{
		baseUrl.GET("",h.GetAllUsers)
		baseUrl.GET("/:id",h.GetById)
		baseUrl.POST("",h.CreateUser)
	}
}