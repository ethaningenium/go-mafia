package route

import (
	"fmt"
	"mafia/services"
	"mafia/storage"

	"github.com/gin-gonic/gin"
)

type emailRequest struct {
	Email string `json:"email" binding:"required"`

}

func RunVerifyRoutes(r *gin.Engine) {
	verify := r.Group("/verify")
	verify.GET("/:code", verifyToken)
	verify.POST("/send", sendEmail)
}

func sendEmail (ctx *gin.Context) {
	//bind json
	var req emailRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get access to db
	db := storage.GetDatabase()

	code, err := storage.GetVerificationCode(db, req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	msg := fmt.Sprintf("<h1>Your verification code is:</h1><p>http://localhost:8082/verify/%s</p>", code)

	//send email
	err = services.SendEmail("Verification", msg, req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//return code
	ctx.JSON(200, gin.H{"message": "Email sent"})
}
	



func verifyToken(ctx *gin.Context) {
	//get code parameter
	code := ctx.Param("code")

	//get access to db
	db := storage.GetDatabase()

	//set isEmailVerified to true
	userData , err := storage.UpdateIsEmailVerified(db, code, true)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//return userData
	ctx.JSON(200, gin.H{"data": userData})
}

