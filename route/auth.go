package route

import (
	"mafia/services"
	"mafia/storage"
	"time"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}



func RunAuthRoutes (r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/login", login)
	auth.POST("/register", register)
}

func login(c *gin.Context) {
	//bind json
	var req userRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get access to db
	db := storage.GetDatabase()

	//get user
	userID, email, password, isEmailVerified, verificationCode, refreshToken, err := storage.GetUserByEmail(db, req.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//compare password
	err = services.ComparePassword(req.Password, password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get jwt refresh token
	token, err := services.GenerateJWTToken(req.Email, time.Now().Add(time.Hour * 30 * 30))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	//save refresh token
	userAccess, err := storage.UpdateRefresh(db, userID, token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	
	//create access token
	accessToken , err := services.GenerateJWTToken(req.Email, time.Now().Add(time.Minute * 20))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	//set headers and cookies
	c.Header("access-token", accessToken)
	c.SetCookie("refresh-token", userAccess.RefreshToken, 3600, "/", "localhost", false, true)



	c.JSON(200, gin.H{"userID": userID, "email": email, "isEmailVerified": isEmailVerified,  "verificationCode": verificationCode, "refreshToken": refreshToken})
}

func register(c *gin.Context) {
	//bind json
	var req userRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get access to db
	db := storage.GetDatabase()

	//hash password
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//create user
	userID, err := storage.CreateUserAccount(db, req.Email, hashedPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get jwt refresh token
	token, err := services.GenerateJWTToken(req.Email, time.Now().Add(time.Hour * 30 * 30))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	link, err := services.GenerateRandomString(20)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//save refresh token
	err = storage.CreateUserAccessData(db, false, link, token, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}


	//create access token
	accessToken , err := services.GenerateJWTToken(req.Email, time.Now().Add(time.Minute * 20))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//set headers and cookies
	c.Header("access-token", accessToken)
	c.SetCookie("refresh-token", token, 3600, "/", "localhost", false, true)

	
	c.JSON(200, gin.H{"message": "success", "id": userID, "email": req.Email, })
}