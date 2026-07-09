package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/reshap0318/api-gateway/internal/di"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/middleware"
	"github.com/reshap0318/api-gateway/internal/proxy"
	"github.com/reshap0318/api-gateway/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	host := helpers.GetEnv("APP_HOST", "0.0.0.0")
	port := helpers.GetEnv("APP_PORT", "8080")
	trustedProxies := helpers.GetEnv("TRUSTED_PROXIES", "")
	allowedOrigins := helpers.GetEnv("ALLOWED_ORIGINS", "*")

	gin.SetMode(helpers.GetEnv("GIN_MODE", "release"))

	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	r := gin.Default()

	if trustedProxies != "" {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			log.Printf("Warning: failed to set trusted proxies: %v", err)
		}
	}

	r.Use(middleware.RateLimit(container.RateLimiter))
	r.Use(middleware.CORS(allowedOrigins))

	r.Static("/storage", "./storage")

	apiGroup := r.Group("/api")
	protected := apiGroup.Group("")
	protected.Use(middleware.JWTAuth(container.Services))

	routes.RegisterAll(r, apiGroup, protected, container.Handlers, container.Access)

	// Dynamic Proxy Engine catch-all: any request that doesn't match a Management API
	// route above falls through here and is resolved against gateway_routes (§4.6).
	// Rate limit is already resolved per-route by RouteManager (§2.15 chain); the handler
	// just enforces it via the shared RateLimiter.
	r.NoRoute(proxy.Handler(container.RouteManager, container.Services, container.Access, container.RateLimiter))

	addr := host + ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
