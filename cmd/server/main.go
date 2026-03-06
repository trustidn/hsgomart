package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/api"
	"github.com/trustidn/hsmart-saas/pkg/config"
	"github.com/trustidn/hsmart-saas/pkg/database"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	r := gin.Default()

	// Allow frontend (Vue dev server) to call API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", healthHandler(db))

	api.RegisterRoutes(r, db, cfg)

	// SPA fallback: serve frontend so /dashboard and other routes work when opening the app from backend URL
	const frontendDist = "frontend/dist"
	if dir, err := filepath.Abs(frontendDist); err == nil {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			r.Static("/assets", filepath.Join(dir, "assets"))
			r.NoRoute(spaFallback(filepath.Join(dir, "index.html")))
		} else {
			r.NoRoute(devSPAMessage)
		}
	} else {
		r.NoRoute(devSPAMessage)
	}

	addr := ":" + cfg.AppPort
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"database": "unavailable",
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"database": "disconnected",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	}
}

// spaFallback serves index.html for non-API paths so client-side routing (e.g. /dashboard) works.
func spaFallback(indexPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := os.Stat(indexPath); err != nil {
			devSPAMessage(c)
			return
		}
		c.File(indexPath)
	}
}

// devSPAMessage tells the user to use the frontend dev server when frontend/dist is not present.
func devSPAMessage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusNotFound, `<!DOCTYPE html><html><head><meta charset="utf-8"><title>HSMart</title></head><body style="font-family:sans-serif;padding:2rem;max-width:480px;"><h1>Frontend not served</h1><p>To open the app (e.g. <strong>/dashboard</strong>), either:</p><ol><li><strong>Development:</strong> Run <code>cd frontend && npm run dev</code>, then open <a href="http://localhost:5173">http://localhost:5173</a></li><li><strong>Production:</strong> Run <code>cd frontend && npm run build</code>, then restart this server</li></ol></body></html>`)
}
