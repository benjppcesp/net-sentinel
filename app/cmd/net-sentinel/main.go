package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const targetURL = "https://www.bing.com"

var (
	httpDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "net_sentinel_http_duration_seconds",
		Help:    "Tiempo de respuesta de la petición HTTP",
		Buckets: prometheus.DefBuckets,
	})
	
	httpSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "net_sentinel_http_success",
		Help: "Estado de la conexión: 1 si fue exitosa, 0 si falló",
	})
)

func init() {
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(httpSuccess)
}

func probeNetwork() {
	go func() {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		for {
			start := time.Now()

			resp, err := client.Get(targetURL)
			duration := time.Since(start).Seconds()

			if err != nil {
				fmt.Printf("❌ Error conectando a %s: %v\n", targetURL, err)
				httpSuccess.Set(0)
			} else {
				fmt.Printf("✅ Conexión exitosa a %s en %.4f segundos\n", targetURL, duration)
				httpSuccess.Set(1)
				httpDuration.Observe(duration)
				resp.Body.Close()
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func main() {
	fmt.Println("📡 Net-Sentinel iniciado. Monitoreando:", targetURL)
	probeNetwork()

	fmt.Println("🎧 Servidor de métricas escuchando en puerto :2112")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
