package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const targetURL = "https://www.bing.com"

// VARIABLES DE MÉTRICAS (PROMETHEUS)
// Usamos "Gauge" (Medidor) porque el valor puede subir y bajar (como un velocímetro).
var (
	// Métrica 1: Latencia (Cuánto tarda en responder la web)
	httpDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "net_sentinel_http_duration_seconds",
		Help: "Tiempo de respuesta de la petición HTTP en segundos",
	})

	// Métrica 2: Éxito (1 = Arriba, 0 = Caído)
	httpSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "net_sentinel_http_success",
		Help: "Estado de la conexión: 1 si fue exitosa, 0 si falló",
	})
)


// Registrar variables para que prometheus sepa que eisten
func init() {
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(httpSuccess)
}

// LÓGICA DE LA SONDA (EL "ROBOT")
// -------------------------------
func probeNetwork() {
	go func() {
		for {
			// 1. Iniciamos el cronómetro
			start := time.Now()

			// 2. Hacemos la petición REAL a Internet
			// http.Get intenta descargar la página principal de Google
			resp, err := http.Get(targetURL)

			// 3. Detenemos el cronómetro y calculamos la duración
			duration := time.Since(start).Seconds()

			// 4. Analizamos el resultado
			if err != nil {
				// SI FALLA (ej. no hay internet):
				fmt.Printf("❌ Error conectando a %s: %v\n", targetURL, err)
				httpSuccess.Set(0) // Reportamos "0" (Fallo) a Prometheus
			} else {
				// SI FUNCIONA:
				fmt.Printf("✅ Conexión exitosa a %s en %.4f segundos\n", targetURL, duration)
				httpSuccess.Set(1)          // Reportamos "1" (Éxito)
				httpDuration.Set(duration)  // Reportamos cuánto tardó
				resp.Body.Close()           // Importante: cerramos la conexión para no saturar memoria
			}

			// 5. Esperamos 5 segundos antes del siguiente test
			time.Sleep(5 * time.Second)
		}
	}()
}

// FUNCIÓN PRINCIPAL
// -----------------
func main() {
	// Arrancamos nuestra sonda en segundo plano
	fmt.Println("📡 Net-Sentinel iniciado. Monitoreando:", targetURL)
	probeNetwork()

	// Exponemos el servidor web para que Prometheus pueda leer los datos
	// Esto estará disponible en http://localhost:2112/metrics
	fmt.Println("🎧 Servidor de métricas escuchando en puerto :2112")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}