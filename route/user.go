package route

import (
	"mafia/services"
	"mafia/storage"

	"net/http"

	"github.com/gin-gonic/gin"
)

type userInfo struct {
	FullName string `json:"name" binding:"required"`
	AvatarUrl string `json:"avatar" binding:"required"`
}


func RunUserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(authMiddleware)
	user.POST("/info", createInfo)
	
}

func authMiddleware(c *gin.Context) {
	
	token := c.GetHeader("Authorization")



	// Проверка токена
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный JWT токен", "token": token})
		return
	}
	token = services.ExtractToken(token)
	


	// Парсинг токена
	claims , err := services.ParseJWTToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "token": token})
		return
	}

	c.Set("email", claims.Email)

	// Продолжите выполнение запроса, если токен валиден
	c.Next()
}

func createInfo(c *gin.Context) {

	//get email from context
	email, _ := c.Get("email")
	emails, _ := email.(string)
	

	//get access to db
	db := storage.GetDatabase()

	//get id by email 
	id, err := storage.GetIDByEmail(db, emails)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}


	//bind json
	var req userInfo
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//save data
	err = storage.CreateUserInfo(db, id, req.FullName, req.AvatarUrl)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}


	
	c.JSON(200, gin.H{"message": "User info created", "id": id, "name": req.FullName, "avatar": req.AvatarUrl})
}