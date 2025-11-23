package main

import (
	"fmt"
	"log"
	"net/http"

	"job-portal-api/internal/config"
	"job-portal-api/internal/db"
	httpServer "job-portal-api/internal/http"
	"job-portal-api/internal/modules/healthcheck"
)

func main(){
	cfg := config.Load()

	pg, err := db.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect PostgreSQL : %v", err)
	}
	defer pg.Close()

	fmt.Println("🚀 PostgreSQL connected")

	// Init modules
	hcService := healthcheck.NewService()
	hcHandler := healthcheck.NewHandler(hcService)

	// Router
	router := httpServer.NewRouter(hcHandler)

	fmt.Printf("🚀 Server running on : %s\n", cfg.AppPort)
	http.ListenAndServe(":"+cfg.AppPort, router)
}