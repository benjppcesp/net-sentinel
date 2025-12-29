package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benjppcesp/net-sentinel/internal/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	target := os.Getenv("TARGET_URL")
	if target == "" {
		target = "https://google.com"
	}

	sonda := monitor.NewSentinel(target, 5*time.Second)

	// Servidor de métricas para Prometheus en puerto 2112
	go func() {
		log.Println("📊 Servidor de métricas en :2112/metrics")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Printf("Error en servidor de métricas: %v", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("🚀 Net-Sentinel Iniciado")
	go sonda.Start(ctx)

	<-ctx.Done()
	log.Println("🛑 Apagando...")
}
