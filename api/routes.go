package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/config"
	"github.com/trustidn/hsmart-saas/pkg/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	svc := NewServiceRegistry(db, cfg)
	h := NewHandlerRegistry(svc, db)

	registerAuthRoutes(r, h, svc)
	registerAPIRoutes(r, h, svc)
	registerAdminRoutes(r, h, svc)
}

func registerAuthRoutes(r *gin.Engine, h *HandlerRegistry, svc *ServiceRegistry) {
	loginLimiter := middleware.NewRateLimiter(5, 1*time.Minute)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", middleware.RateLimit(loginLimiter), h.Auth.Login)
		auth.POST("/refresh", h.Auth.Refresh)
		auth.GET("/profile", middleware.Auth(svc.Auth), h.Auth.Profile)
	}
}

func registerAPIRoutes(r *gin.Engine, h *HandlerRegistry, svc *ServiceRegistry) {
	// Routes that bypass subscription middleware (so expired tenants can still upgrade)
	bypass := r.Group("/api")
	bypass.Use(middleware.Auth(svc.Auth), middleware.Tenant())
	{
		bypass.GET("/subscription", h.Subscription.GetSubscription)
		bypass.GET("/subscription/plans", h.Subscription.ListPlans)
		bypass.POST("/subscription/order", h.Order.CreateOrder)
		bypass.POST("/subscription/order/:id/payment", h.Order.UploadPaymentProof)
		bypass.GET("/subscription/orders", h.Order.ListOrders)
		bypass.GET("/subscription/order/:id", h.Order.GetOrder)

		bypass.GET("/tenant/profile", h.Tenant.GetProfile)
		bypass.PUT("/tenant/profile", h.Tenant.UpdateProfile)
		bypass.POST("/tenant/logo", h.Tenant.UploadLogo)
	}

	api := r.Group("/api")
	api.Use(middleware.Auth(svc.Auth), middleware.Tenant(), middleware.Subscription(svc.Subscription))

	registerCashierRoutes(api, h)
	registerOwnerRoutes(api, h)
}

func registerCashierRoutes(api *gin.RouterGroup, h *HandlerRegistry) {
	g := api.Group("")
	g.Use(middleware.Role("owner", "cashier"))
	{
		g.GET("/products", h.Product.ListProducts)
		g.GET("/products/barcode/:barcode", h.Product.GetProductByBarcode)
		g.GET("/products/:id/stock", h.Inventory.GetStock)

		g.POST("/pos/checkout", h.POS.Checkout)
		g.GET("/pos/receipt/:id", h.Report.GetReceipt)

		g.GET("/reports/sales", h.Report.SalesSummary)
		g.GET("/reports/products", h.Report.TopProducts)
		g.GET("/reports/inventory", h.Report.InventorySummary)

		g.POST("/shifts/open", h.Shift.OpenShift)
		g.POST("/shifts/close", h.Shift.CloseShift)
		g.GET("/shifts/current", h.Shift.GetCurrentShift)

		g.GET("/inventory/low-stock", h.Inventory.LowStock)
		g.GET("/inventory/expiring", h.Inventory.Expiring)
	}
}

func registerOwnerRoutes(api *gin.RouterGroup, h *HandlerRegistry) {
	g := api.Group("")
	g.Use(middleware.Role("owner"))
	{
		// User management
		g.GET("/users", h.User.List)
		g.POST("/users", h.User.Create)
		g.PUT("/users/:id", h.User.Update)
		g.DELETE("/users/:id", h.User.Delete)

		// Categories
		g.GET("/categories", h.Product.ListCategories)
		g.POST("/categories", h.Product.CreateCategory)
		g.PUT("/categories/:id", h.Product.UpdateCategory)
		g.DELETE("/categories/:id", h.Product.DeleteCategory)

		// Products
		g.POST("/products", h.Product.CreateProduct)
		g.GET("/products/:id", h.Product.GetProduct)
		g.PUT("/products/:id", h.Product.UpdateProduct)
		g.DELETE("/products/:id", h.Product.DeleteProduct)
		g.POST("/products/:id/barcodes", h.Product.AddBarcode)
		g.GET("/products/:id/barcodes", h.Product.ListBarcodes)
		g.DELETE("/products/:id/barcodes/:barcode", h.Product.DeleteBarcode)

		// Inventory
		g.GET("/inventory", h.Inventory.List)
		g.GET("/inventory/movements", h.Inventory.ListMovements)
		g.POST("/products/:id/adjust-stock", h.Inventory.AdjustStock)

		// Purchases
		g.GET("/purchases", h.Purchase.List)
		g.GET("/purchases/:id", h.Purchase.GetByID)
		g.POST("/purchases", h.Purchase.Create)

		// Reports
		g.GET("/reports/sales/daily", h.Report.SalesDaily)
		g.GET("/reports/sales/hourly", h.Report.SalesHourly)
		g.GET("/reports/sales/transactions", h.Report.SalesTransactions)
		g.GET("/reports/sales/compare", h.Report.SalesCompare)
		g.GET("/reports/payments", h.Report.PaymentsReport)
		g.GET("/reports/profit", h.Report.ProfitReport)
		g.GET("/reports/cashiers", h.Report.CashiersReport)
		g.GET("/reports/shifts", h.Report.ShiftsReport)
		g.GET("/reports/products/margin", h.Report.ProductMargin)

		// Shifts
		g.GET("/shifts", h.Shift.ListShifts)

		// Subscription (upgrade plan directly without order flow)
		g.PUT("/subscription", h.Subscription.ChangePlan)

		// Refunds
		g.POST("/pos/refund", h.Refund.CreateRefund)
		g.GET("/refunds", h.Refund.ListRefunds)

		// Stock opname
		g.POST("/inventory/opname", h.Opname.Start)
		g.PUT("/inventory/opname/:id", h.Opname.SubmitItems)
		g.POST("/inventory/opname/:id/approve", h.Opname.Approve)
		g.GET("/inventory/opname/:id", h.Opname.Get)
		g.GET("/inventory/opnames", h.Opname.List)
	}
}

func registerAdminRoutes(r *gin.Engine, h *HandlerRegistry, svc *ServiceRegistry) {
	g := r.Group("/admin")
	g.Use(middleware.Auth(svc.Auth), middleware.Role("superadmin"))
	{
		g.GET("/tenants", h.Admin.ListTenants)
		g.GET("/tenants/:id", h.Admin.GetTenant)
		g.PUT("/tenants/:id", h.Admin.UpdateTenant)
		g.GET("/subscriptions", h.Admin.ListSubscriptions)
		g.PUT("/subscriptions/:id", h.Admin.UpdateSubscription)
		g.GET("/stats", h.Admin.Stats)

		g.GET("/orders", h.Admin.ListOrders)
		g.GET("/orders/:id", h.Admin.GetOrderDetail)
		g.PUT("/orders/:id/approve", h.Admin.ApproveOrder)
		g.PUT("/orders/:id/reject", h.Admin.RejectOrder)
	}
}
