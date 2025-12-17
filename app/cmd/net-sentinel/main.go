package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	targetURL      string
	checkInterval  time.Duration
	requestTimeout time.Duration
)

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

	loadConfig()
}

func loadConfig() {
	targetURL = getENV("TARGET_URL", "https://www.bing.com")
	checkInterval = getDurationEnv("CHECK_INTERVAL", 5*time.Second)
	requestTimeout = getDurationEnv("REQUEST_TIMEOUT", 5*time.Second)

}

func getENV(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("variable %s invalida: %s"))
	}
	return duration
}

func probeNetwork() {
	go func() {
		client := &http.Client{
			Timeout: requestTimeout,
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

			time.Sleep(checkInterval)
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
