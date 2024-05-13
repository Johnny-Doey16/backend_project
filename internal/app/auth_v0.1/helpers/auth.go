package helpers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// TODO: FIX ISSUE

func CheckUserType(ctx *gin.Context, role string) (err error) {
	userRole := ctx.GetString("role")
	err = nil
	if userRole != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(ctx *gin.Context, userId string) (err error) {
	userRole := ctx.GetString("role")
	uid := ctx.GetString("uid")
	err = nil

	// Only Users can view their profile. Or Admin
	if userRole == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(ctx, userRole)

	return err
}

func HashPassword(pwd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(pwd string, providedPwd string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPwd), []byte(pwd))
	check := true
	msg := ""
	if err != nil {
		check = false
		msg = "Email or password incorrect"
	}
	return check, msg
}
