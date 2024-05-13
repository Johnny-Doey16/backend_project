package helpers

import (
	"context"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/steve-mir/diivix_backend/internal/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	Role       string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = db.OpenCollection(db.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// func GenerateAllTokens(email string, firstName string, lastName string, role string, uid string) (signedToken string, signedRefresh string, err error) {
// 	claims := &SignedDetails{
// 		Email:      email,
// 		First_name: firstName,
// 		Last_name:  lastName,
// 		User_role:  role,
// 		Role:       role,
// 		// TODO: LOOK INTO
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
// 			Issuer:    "Diivix",
// 		},
// 	}

// 	refreshClaims := &SignedDetails{
// 		// Email:      email,
// 		// First_name: firstName,
// 		// Last_name:  lastName,
// 		// User_role:  role,
// 		// Role:       role,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
// 			Issuer:    "Diivix",
// 		},
// 	}

// 	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
// 	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}

// 	return token, refreshToken, err
// }

func GenerateAllTokens(email string, firstName string, lastName string, role string, uid string) (signedToken string, signedRefresh string, err error) {
	// Set token expiration times
	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()      // 24 hours
	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7).Unix() // 7 days

	// Create access token claims
	accessTokenClaims := jwt.MapClaims{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"Uid":        uid,
		"role":       role,
		"exp":        accessTokenExp,
		"iss":        "Diivix",
	}

	// Create refresh token claims
	refreshTokenClaims := jwt.MapClaims{
		"exp": refreshTokenExp,
		"iss": "Diivix",
	}

	// Create and sign the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	// Create and sign the refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

// func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var updatedObj primitive.D

// 	updatedObj = append(updatedObj, bson.E{"access_token": signedToken})
// 	updatedObj = append(updatedObj, bson.E{"refresh_token" signedRefreshToken})

// 	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	updatedObj = append(updatedObj, bson.E{"updated_at", updatedAt})

// 	upsert := true
// 	filter := bson.M{"user_id": userId}
// 	opt := options.UpdateOptions{
// 		Upsert: &upsert,
// 	}

// 	_, err := userCollection.UpdateOne(
// 		ctx,
// 		filter,
// 		bson.D{
// 			{"$set", updatedObj},
// 		},
// 		&opt,
// 	)

// 	defer cancel()

// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}
// 	return
// }

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updatedObj := bson.M{
		"access_token":  signedToken,
		"refresh_token": signedRefreshToken,
		"updated_at":    time.Now().Format(time.RFC3339),
	}

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.Update().SetUpsert(upsert)

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": updatedObj},
		opt,
	)

	if err != nil {
		return err
	}

	return nil
}

// func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var updatedObj primitive.D

// 	updatedObj = append(updatedObj, bson.E{Key: "access_token", Value: signedToken}) // Fix the typo here
// 	updatedObj = append(updatedObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

// 	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: updatedAt})

// 	upsert := true
// 	filter := bson.M{"user_id": userId}
// 	opt := options.UpdateOptions{
// 		Upsert: &upsert,
// 	}

// 	_, err := userCollection.UpdateOne(
// 		ctx,
// 		filter,
// 		bson.D{
// 			{Key: "$set", Value: updatedObj},
// 		},
// 		&opt,
// 	)

// 	defer cancel()

// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}
// 	return
// }

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("The token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
