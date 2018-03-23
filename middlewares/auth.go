package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"

	"devin/crypto"
	"devin/database"
	"devin/helpers"
	"devin/models"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check exist token
		_, ok := r.Header["Authorization"]
		if !ok {
			cookie, e := r.Cookie("Authorization")
			if e != nil {
				models.User{}.ExpireAuthorizationCookie(w)

				err := helpers.ErrorResponse{}
				err.ErrorCode = http.StatusUnauthorized
				err.Message = "Authentication Token not found"
				log.Println("Authentication Token not found")
				helpers.NewErrorResponse(w, &err)

				return
			}
			r.Header.Set("Authorization", cookie.Value)
		}

		// validate the token
		token, err := jwt_request.ParseFromRequestWithClaims(r, jwt_request.HeaderExtractor{"Authorization"}, &models.Claim{}, func(token *jwt.Token) (interface{}, error) {
			return crypto.GetJWTVerifyKey()
		})

		if err != nil {
			switch err.(type) {
			case *jwt.ValidationError:
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					models.User{}.ExpireAuthorizationCookie(w)

					err := helpers.ErrorResponse{}
					err.ErrorCode = http.StatusUnauthorized
					err.Message = "Token Expired, get a new one."
					log.Println("Token Expired, get a new one.")
					helpers.NewErrorResponse(w, &err)

					return

				default:
					models.User{}.ExpireAuthorizationCookie(w)

					err := helpers.ErrorResponse{}
					err.ErrorCode = http.StatusUnauthorized
					err.Message = "Error while Parsing Token!"
					log.Println("Error while Parsing Token!", vErr)
					helpers.NewErrorResponse(w, &err)

					return
				}

			default:
				models.User{}.ExpireAuthorizationCookie(w)

				err := helpers.ErrorResponse{}
				err.ErrorCode = http.StatusUnauthorized
				err.Message = "Error while Parsing Token!"
				log.Println("Error while Parsing Token!")
				helpers.NewErrorResponse(w, &err)

				return
			}
		}

		if !token.Valid {
			models.User{}.ExpireAuthorizationCookie(w)

			err := helpers.ErrorResponse{}
			err.ErrorCode = http.StatusUnauthorized
			err.Message = "Auhtentication failed (Invalid token)."
			log.Println("Auhtentication failed (Invalid token).")
			helpers.NewErrorResponse(w, &err)

			return
		}
		var user models.User
		claim := token.Claims.(*models.Claim)
		authUser, e := user.ExtractUserFromClaimPayload(claim.Payload)

		if e != nil {
			models.User{}.ExpireAuthorizationCookie(w)

			err := helpers.ErrorResponse{}
			err.ErrorCode = http.StatusUnauthorized
			err.Message = "Auhtentication failed (Bad payload)."
			log.Println("Auhtentication failed,", e)
			helpers.NewErrorResponse(w, &err)

			return
		}

		db := database.NewPGInstance()
		defer db.Close()

		db.Model(&user).Where("id=?", authUser.ID).First()
		if user.ID == 0 {
			models.User{}.ExpireAuthorizationCookie(w)

			err := helpers.ErrorResponse{}
			err.ErrorCode = http.StatusUnauthorized
			err.Message = "Auhtentication failed (User not found)."
			log.Println("Auhtentication failed,", e)
			helpers.NewErrorResponse(w, &err)

			return
		}

		// Check token remaining lifetime
		if float64(claim.ExpiresAt-time.Now().Unix())/user.TokenLifetime().Seconds() < .25 {
			log.Println("Token life time less than 0.25 of total lifetime")

			claim := user.GenerateNewTokenClaim()
			tokenString, err := user.GenerateNewTokenString(claim)
			if err != nil {
				helpers.NewErrorResponse(w, err)
				return
			}
			user.SetAuthorizationCookieAndHeader(w, tokenString)

			ctx := context.WithValue(r.Context(), "Authorization", claim)

			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			user.SetAuthorizationCookieAndHeader(w, token.Raw)
			ctx := context.WithValue(r.Context(), "Authorization", claim)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
