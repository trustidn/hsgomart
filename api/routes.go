package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/internal/auth"
	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/pos"
	"github.com/trustidn/hsmart-saas/internal/product"
	"github.com/trustidn/hsmart-saas/internal/purchase"
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

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/profile", middleware.Auth(authSvc), authHandler.Profile)
	}

	subscriptionSvc := subscription.NewService(db)
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.Auth(authSvc), middleware.Tenant(), middleware.Subscription(subscriptionSvc))

	// Owner + Cashier: product read (for POS) and checkout. Register first so GET /products matches here.
	cashier := apiGroup.Group("")
	cashier.Use(middleware.Role("owner", "cashier"))
	{
		productSvc := product.NewService(db, subscriptionSvc)
		productHandler := product.NewHandler(productSvc)
		cashier.GET("/products", productHandler.ListProducts)
		cashier.GET("/products/barcode/:barcode", productHandler.GetProductByBarcode)
		inventorySvc := inventory.NewService(db)
		inventoryHandler := inventory.NewHandler(inventorySvc)
		cashier.GET("/products/:id/stock", inventoryHandler.GetStock)

		posSvc := pos.NewService(db)
		posHandler := pos.NewHandler(posSvc)
		cashier.POST("/pos/checkout", posHandler.Checkout)
	}

	// Owner-only: users, categories, product write, inventory write, purchases, reports
	owner := apiGroup.Group("")
	owner.Use(middleware.Role("owner"))
	{
		userSvc := user.NewService(db, subscriptionSvc)
		userHandler := user.NewHandler(userSvc)
		owner.GET("/users", userHandler.List)
		owner.POST("/users", userHandler.Create)
		owner.PUT("/users/:id", userHandler.Update)
		owner.DELETE("/users/:id", userHandler.Delete)

		productSvc := product.NewService(db, subscriptionSvc)
		productHandler := product.NewHandler(productSvc)
		owner.GET("/categories", productHandler.ListCategories)
		owner.POST("/categories", productHandler.CreateCategory)
		owner.PUT("/categories/:id", productHandler.UpdateCategory)
		owner.DELETE("/categories/:id", productHandler.DeleteCategory)
		owner.POST("/products", productHandler.CreateProduct)
		owner.GET("/products/:id", productHandler.GetProduct)
		owner.PUT("/products/:id", productHandler.UpdateProduct)
		owner.DELETE("/products/:id", productHandler.DeleteProduct)
		owner.POST("/products/:id/barcodes", productHandler.AddBarcode)

		inventorySvc := inventory.NewService(db)
		inventoryHandler := inventory.NewHandler(inventorySvc)
		owner.GET("/inventory", inventoryHandler.List)
		owner.GET("/inventory/movements", inventoryHandler.ListMovements)
		owner.POST("/products/:id/adjust-stock", inventoryHandler.AdjustStock)

		purchaseSvc := purchase.NewService(db)
		purchaseHandler := purchase.NewHandler(purchaseSvc)
		owner.GET("/purchases", purchaseHandler.List)
		owner.GET("/purchases/:id", purchaseHandler.GetByID)
		owner.POST("/purchases", purchaseHandler.Create)

		reportSvc := report.NewService(db)
		reportHandler := report.NewHandler(reportSvc)
		owner.GET("/reports/sales", reportHandler.SalesSummary)
		owner.GET("/reports/sales/daily", reportHandler.SalesDaily)
		owner.GET("/reports/sales/hourly", reportHandler.SalesHourly)
		owner.GET("/reports/sales/transactions", reportHandler.SalesTransactions)
		owner.GET("/reports/payments", reportHandler.PaymentsReport)
		owner.GET("/reports/profit", reportHandler.ProfitReport)
		owner.GET("/reports/products", reportHandler.TopProducts)
		owner.GET("/reports/inventory", reportHandler.InventorySummary)
		owner.GET("/reports/cashiers", reportHandler.CashiersReport)
	}
}
