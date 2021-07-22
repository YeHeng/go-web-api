package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/YeHeng/go-web-api/internal/api/repository/db_repo"
	"github.com/YeHeng/go-web-api/internal/api/repository/db_repo/user_repo"
	"github.com/YeHeng/go-web-api/internal/pkg/logger"
	"github.com/YeHeng/go-web-api/pkg/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const identityKey = "id"

type Credential struct {
	UserName string
	Roles    []*Role
}

type Role struct {
	gorm.Model
	Name string
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var jwtAuthMiddleware *jwt.GinJWTMiddleware

func init() {
	AddMiddleware(&jwtMiddleware{})
}

type jwtMiddleware struct {
}

func (m *jwtMiddleware) Apply(r *gin.Engine) {

	log := logger.Get()

	r.POST("/login", jwtAuthMiddleware.LoginHandler)
	r.GET("/refresh_token", jwtAuthMiddleware.RefreshHandler)
	r.GET("/logout", jwtAuthMiddleware.LogoutHandler)

	r.NoRoute(jwtAuthMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Info(fmt.Sprintf("NoRoute claims: %#v\n", claims))
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}

func (m *jwtMiddleware) Destroy() {
}

func (m *jwtMiddleware) Get() gin.HandlerFunc {
	return jwtAuthMiddleware.MiddlewareFunc()
}

func (m *jwtMiddleware) Init() {

	log := logger.Get()
	var err error

	jwtAuthMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
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
		log.Fatal("JWT Error:" + err.Error())
		panic(err)
	}

	// When you use auth.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := jwtAuthMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		panic(err)
	}

}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	md5 := util.GeneratePassword(password)

	qb := user_repo.NewQueryBuilder()
	qb.WhereIsDeleted(db_repo.EqualPredicate, 0)
	qb.WhereUsername(db_repo.EqualPredicate, userID)
	qb.WherePassword(db_repo.EqualPredicate, md5)
	u, err := qb.QueryOne()
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &Credential{
		UserName: u.Username,
		Roles:    nil,
	}, nil

}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*Credential); ok {
		return jwt.MapClaims{
			identityKey: v.UserName,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &Credential{
		UserName: claims[identityKey].(string),
	}
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*Credential); ok && v.UserName == "admin" {
		return true
	}

	return false
}

func Unauthorized(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"_link":   "/login",
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
		"code":    http.StatusOK,
		"message": "logout success",
	})
}

func RefreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}
