package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/internal/auth"
	"github.com/trustidn/hsmart-saas/internal/user"
	"github.com/trustidn/hsmart-saas/pkg/config"
	"github.com/trustidn/hsmart-saas/pkg/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	authSvc := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	// Auth routes (public except profile). No tenant middleware; profile uses Auth only.
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/profile", middleware.Auth(authSvc), authHandler.Profile)
	}

	// Protected API: require Auth then Tenant. All handlers use tenant_id from context for isolation.
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.Auth(authSvc), middleware.Tenant())
	{
		userSvc := user.NewService(db)
		userHandler := user.NewHandler(userSvc)
		apiGroup.GET("/users", userHandler.List)
		apiGroup.POST("/users", userHandler.Create)
		apiGroup.PUT("/users/:id", userHandler.Update)
		apiGroup.DELETE("/users/:id", userHandler.Delete)
	}
}
