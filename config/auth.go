package config

import (
	"errors"
	"net/http"
	"strings"

	"github.com/calculeat/main_rest_api/logger"
	"github.com/calculeat/main_rest_api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	FIRE_AUTH_CLIENT *auth.Client
	FIRE_USER        *auth.UserRecord
	CREDENTIAL_FILE  string
	CURRENT_USER     models.User
)

type UserInfoProvider struct {
	Email      string `json:"email"`
	ProviderId string `json:"providerId"`
	RawId      string `json:"rawId"`
}

// Extract the requested part of the JWT
func ExtractJWT(jwt_token string, part_to_extract string) (requested_data interface{}) {
	// Parse the JWT token
	token, err := jwt.Parse(jwt_token, func(token *jwt.Token) (interface{}, error) {
		return nil, nil // No key is required for parsing, as we're only interested in claims
	})
	if err != nil {
		logger.Log.Fatalf("Error parsing token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		logger.Log.Fatal("Invalid token")
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Log.Fatal("Unable to extract claims")
	}

	// Extract the "sign_in_provider" claim
	requested_data, exists := claims[part_to_extract]
	if !exists {
		logger.Log.Fatalf("%s claim not found", part_to_extract)
	}
	return requested_data

}

func init() {
	// Create a Firebase app instance
	CREDENTIAL_FILE = logger.GoDotEnvVariable("CREDENTIAL_FILE")
	opt := option.WithCredentialsFile(CREDENTIAL_FILE)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Log.Fatalf("Failed to create Firebase instance: %v", err)
	}

	// Create a Firebase auth client instance
	authClient, err := app.Auth(context.Background())
	if err != nil {
		logger.Log.Fatalf("Failed to create Firebase auth client: %v", err)
	}

	FIRE_AUTH_CLIENT = authClient

}

func AuthMiddleware() gin.HandlerFunc {
	logger.Log.Debugln("AuthMiddleware called.")

	return func(ctx *gin.Context) {

		// Get JWT value from request headers
		authorization_header := ctx.Request.Header.Get("Authorization")

		// Check if the header is present
		if authorization_header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Authorization header missing"})
			logger.Log.Error("Missing authorization token.")
			return
		}

		// Extract the JWT token from the header (assuming it's in the "Bearer" scheme)
		authParts := strings.Split(authorization_header, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Invalid authorization header!"})
			logger.Log.Errorf("Invalid authorization header!. Provided Authorization header: %v", authorization_header)
			return
		}

		jwt_token := authParts[1]
		logger.Log.Debugf("JWT token that will be verified: %v", jwt_token)
		verified_jwt_token, err := FIRE_AUTH_CLIENT.VerifyIDToken(ctx, jwt_token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Invalid authorization header!"})
			logger.Log.Errorf("error verifying ID token: %v", err)
		}

		logger.Log.Debugf("JWT token verified. Verified Token: %v", verified_jwt_token)

		uid := verified_jwt_token.UID
		logger.Log.Debugf("Extracted JWT UID: %v", uid)

		user, err := FIRE_AUTH_CLIENT.GetUser(context.Background(), uid)

		if err != nil {
			logger.Log.Errorf("Cannot find the firebase user: %v", err)
		}
		FIRE_USER = user
		CURRENT_USER, err = findUserByEmail(FIRE_USER.Email)

	}
}

func findUserByEmail(email string) (models.User, error) {
	var user models.User

	err := DB.Where("email=?", email).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	if user.Email != email {
		logger.Log.Debugf("Email not found! Input email: %v, Database model email: %v", email, user.Email)
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}
