package middleware

import (
	"github.com/YeHeng/go-web-api/pkg/logger"
	"net/http"
	"time"

	"github.com/YeHeng/go-web-api/common/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const identityKey = "id"

var JwtAuthMiddleware *jwt.GinJWTMiddleware

func InitJwt(r *gin.Engine) {
	JwtAuthMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Golang Web Tools",
		Key:             []byte("OTM4QzgzMDktODRDNi00RDcyLUI5ODctQzEzMEU0ODQwNThECg=="),
		SecureCookie:    true,
		CookieName:      "auth",
		SendCookie:      true,
		CookieHTTPOnly:  true,
		CookieMaxAge:    time.Hour,
		CookieSameSite:  http.SameSiteDefaultMode,
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     PayloadFunc,
		IdentityHandler: IdentityHandler,
		Authenticator:   Authenticator,
		Authorizator:    Authorizator,
		Unauthorized:    Unauthorized,
		LogoutResponse:  LogoutResponse,
		LoginResponse:   LoginResponse,
		RefreshResponse: RefreshResponse,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: auth",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		logger.Logger.Fatal("JWT Error:" + err.Error())
		panic(err)
	}

	// When you use auth.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := JwtAuthMiddleware.MiddlewareInit()

	if errInit != nil {
		logger.Logger.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		panic(err)
	}

	r.POST("/login", JwtAuthMiddleware.LoginHandler)
	r.GET("/refresh_token", JwtAuthMiddleware.RefreshHandler)
	r.GET("/logout", JwtAuthMiddleware.LogoutHandler)

	r.NoRoute(JwtAuthMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logger.Logger.Infof("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals model.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	var user model.User
	if err := Db.Where("Username = ? AND Password = ?", userID, password).First(&user).Error; err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &model.Credential{
		UserName: user.Username,
		Roles:    nil,
	}, nil

}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*model.Credential); ok {
		return jwt.MapClaims{
			identityKey: v.UserName,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &model.Credential{
		UserName: claims[identityKey].(string),
	}
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*model.Credential); ok && v.UserName == "admin" {
		return true
	}

	return false
}

func Unauthorized(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"_href":   c.Request.Proto,
		"code":    code,
		"message": msg,
	})
}

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

func LogoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func RefreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}
