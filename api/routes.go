package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/internal/auth"
	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/pos"
	"github.com/trustidn/hsmart-saas/internal/product"
	"github.com/trustidn/hsmart-saas/internal/report"
	"github.com/trustidn/hsmart-saas/internal/subscription"
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

	// Protected API: require Auth, Tenant, and valid Subscription (active/trial).
	subscriptionSvc := subscription.NewService(db)
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.Auth(authSvc), middleware.Tenant(), middleware.Subscription(subscriptionSvc))
	{
		userSvc := user.NewService(db, subscriptionSvc)
		userHandler := user.NewHandler(userSvc)
		apiGroup.GET("/users", userHandler.List)
		apiGroup.POST("/users", userHandler.Create)
		apiGroup.PUT("/users/:id", userHandler.Update)
		apiGroup.DELETE("/users/:id", userHandler.Delete)

		productSvc := product.NewService(db, subscriptionSvc)
		productHandler := product.NewHandler(productSvc)
		apiGroup.GET("/categories", productHandler.ListCategories)
		apiGroup.POST("/categories", productHandler.CreateCategory)
		apiGroup.PUT("/categories/:id", productHandler.UpdateCategory)
		apiGroup.DELETE("/categories/:id", productHandler.DeleteCategory)
		apiGroup.GET("/products", productHandler.ListProducts)
		apiGroup.GET("/products/barcode/:barcode", productHandler.GetProductByBarcode)
		apiGroup.POST("/products", productHandler.CreateProduct)
		apiGroup.GET("/products/:id", productHandler.GetProduct)
		apiGroup.PUT("/products/:id", productHandler.UpdateProduct)
		apiGroup.DELETE("/products/:id", productHandler.DeleteProduct)
		apiGroup.POST("/products/:id/barcodes", productHandler.AddBarcode)

		inventorySvc := inventory.NewService(db)
		inventoryHandler := inventory.NewHandler(inventorySvc)
		apiGroup.GET("/inventory", inventoryHandler.List)
		apiGroup.GET("/inventory/movements", inventoryHandler.ListMovements)
		apiGroup.GET("/products/:id/stock", inventoryHandler.GetStock)
		apiGroup.POST("/products/:id/adjust-stock", inventoryHandler.AdjustStock)

		posSvc := pos.NewService(db)
		posHandler := pos.NewHandler(posSvc)
		apiGroup.POST("/pos/checkout", posHandler.Checkout)

		reportSvc := report.NewService(db)
		reportHandler := report.NewHandler(reportSvc)
		apiGroup.GET("/reports/sales", reportHandler.SalesSummary)
		apiGroup.GET("/reports/products", reportHandler.TopProducts)
		apiGroup.GET("/reports/inventory", reportHandler.InventorySummary)
	}
}
