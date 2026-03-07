package api

import (
	"github.com/trustidn/hsmart-saas/internal/admin"
	"github.com/trustidn/hsmart-saas/internal/auth"
	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/opname"
	"github.com/trustidn/hsmart-saas/internal/order"
	"github.com/trustidn/hsmart-saas/internal/pos"
	"github.com/trustidn/hsmart-saas/internal/product"
	"github.com/trustidn/hsmart-saas/internal/purchase"
	"github.com/trustidn/hsmart-saas/internal/refund"
	"github.com/trustidn/hsmart-saas/internal/report"
	"github.com/trustidn/hsmart-saas/internal/shift"
	"github.com/trustidn/hsmart-saas/internal/subscription"
	"github.com/trustidn/hsmart-saas/internal/tenant"
	"github.com/trustidn/hsmart-saas/internal/user"
	"github.com/trustidn/hsmart-saas/pkg/config"
	"gorm.io/gorm"
)

type ServiceRegistry struct {
	Auth         *auth.Service
	Product      *product.Service
	Inventory    *inventory.Service
	POS          *pos.Service
	Report       *report.Service
	Shift        *shift.Service
	User         *user.Service
	Purchase     *purchase.Service
	Refund       *refund.Service
	Opname       *opname.Service
	Subscription *subscription.Service
	Tenant       *tenant.Service
	Order        *order.Service
}

type HandlerRegistry struct {
	Auth         *auth.Handler
	Product      *product.Handler
	Inventory    *inventory.Handler
	POS          *pos.Handler
	Report       *report.Handler
	Shift        *shift.Handler
	User         *user.Handler
	Purchase     *purchase.Handler
	Refund       *refund.Handler
	Opname       *opname.Handler
	Subscription *subscription.Handler
	Tenant       *tenant.Handler
	Order        *order.Handler
	Admin        *admin.Handler
}

func NewServiceRegistry(db *gorm.DB, cfg config.Config) *ServiceRegistry {
	subscriptionSvc := subscription.NewService(db)
	return &ServiceRegistry{
		Auth:         auth.NewService(db, cfg.JWTSecret),
		Product:      product.NewService(db, subscriptionSvc),
		Inventory:    inventory.NewService(db),
		POS:          pos.NewService(db),
		Report:       report.NewService(db),
		Shift:        shift.NewService(db),
		User:         user.NewService(db, subscriptionSvc),
		Purchase:     purchase.NewService(db),
		Refund:       refund.NewService(db),
		Opname:       opname.NewService(db),
		Subscription: subscriptionSvc,
		Tenant:       tenant.NewService(db),
		Order:        order.NewService(db),
	}
}

func NewHandlerRegistry(svc *ServiceRegistry, db *gorm.DB) *HandlerRegistry {
	return &HandlerRegistry{
		Auth:         auth.NewHandler(svc.Auth),
		Product:      product.NewHandler(svc.Product),
		Inventory:    inventory.NewHandler(svc.Inventory),
		POS:          pos.NewHandler(svc.POS),
		Report:       report.NewHandler(svc.Report),
		Shift:        shift.NewHandler(svc.Shift),
		User:         user.NewHandler(svc.User),
		Purchase:     purchase.NewHandler(svc.Purchase),
		Refund:       refund.NewHandler(svc.Refund),
		Opname:       opname.NewHandler(svc.Opname),
		Subscription: subscription.NewHandler(svc.Subscription),
		Tenant:       tenant.NewHandler(svc.Tenant),
		Order:        order.NewHandler(svc.Order),
		Admin:        admin.NewHandler(db),
	}
}
