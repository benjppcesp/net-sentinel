package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Definicion de Objeto (POO)
type Sentinel struct {
	TargetURL      string        `json:"target_url"`
	CheckInterval  time.Duration `json:"check_interval"`
	RequestTimeout time.Duration `json:"request_timeout"`
	Status         string        `json:"status"`
	mu             sync.RWMutex  // Protege el objeto de accesos simultáneos
}

func (s *Sentinel) SetStatus(newStatus string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = newStatus
}

var (
	httpDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "net_sentinel_http_duration_seconds",
		Help:    "Tiempo de respuesta de la petición HTTP",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
	})

	httpSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "net_sentinel_http_success",
		Help: "Estado de la conexión: 1 si fue exitosa, 0 si falló",
	})

	apiRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "net_sentinel_api_calls_total",
		Help: "Total de consultas a la API de Control",
	})
)

func init() {
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(httpSuccess)

}

func (s *Sentinel) RunProbe() {
	go func() {
		for {
			s.mu.RLock()
			url := s.TargetURL
			timeout := s.RequestTimeout
			interval := s.CheckInterval
			s.mu.RUnlock()

			client := &http.Client{Timeout: timeout}
			start := time.Now()

			resp, err := client.Get(url)
			duration := time.Since(start).Seconds()

			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				httpSuccess.Set(0)
				s.mu.Lock()
				s.Status = "Error de Conexion"
				s.mu.Unlock()
			} else {
				httpSuccess.Set(1)
				httpDuration.Observe(duration)
				resp.Body.Close()
				s.mu.Lock()
				s.Status = "Online"
				s.mu.Unlock()
			}
			time.Sleep(interval)

		}
	}()
}

func main() {

	sonda := &Sentinel{
		TargetURL:      getENV("TARGET_URL", "https://www.bing.com"),
		CheckInterval:  5 * time.Second,
		RequestTimeout: 5 * time.Second,
		Status:         "Iniciando",
	}
	sonda.RunProbe()

	muxMetrics := http.NewServerMux()
	muxMetrics.Handle("/metrics", promhttp.Handler())

	go func() {
		fmt.Println("📊 Métricas en :2112/metrics")
		if err := http.ListenAndServe(":2112", muxMetrics); err
			fmt.Printf("error de Metrcias: %v\n", err)
		
	}()
	// Endopoint de API (Puerto 8080)
	muxAPI := http.NewServerMux()
	muxAPI.HandleFunc("/status", func (w http.ResponseWriter, r *http.Reques)  {
		apiRequest.Inc()
		w.Header().Set("content-type", "application/json")
		sonda.mu.RLock()
		json.NewEncoder(w).Encode(sonda)
		sonda.mu.RUnlock()
	})


	fmt.Println("🚀 API de Control en :8080/status")
	if err := http.ListenAndServe(":8080", muxAPI); err != nil {
		fmt.Printf("Error de API: %v\n", err)
	}
}

func getENV(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
