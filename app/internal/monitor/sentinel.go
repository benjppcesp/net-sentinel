package monitor

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sentinel_checks_total",
		Help: "Numero total de chequeos realizados",
	})

	IsUpMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sentinel_site_up",
		Help: "Estado del sitio (1 para UP, 0 para DOWN)",
	}, []string{"url"})

	latencyMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "sentinel_latency_seconds",
		Help:       "Latencia de la respuesta HTTP",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"url"})
)

type Sentinel struct {
	TargetURL string
	Interval  time.Duration
	client    *http.Client
	mu        sync.RWMutex
}

func NewSentinel(url string, interval time.Duration) *Sentinel {
	return &Sentinel{
		TargetURL: url,
		Interval:  interval,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *Sentinel) Start(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Sonda para %s detenida", s.TargetURL)
			return
		case <-ticker.C:
			s.Check(ctx)

		}
	}
}

func (s *Sentinel) Check(ctx context.Context) {
	req, _ := http.NewRequestWithContext(ctx, "GET", s.TargetURL, nil)

	start := time.Now()
	resp, err := s.client.Do(req)
	duration := time.Since(start)

	opsProcessed.Inc() // Incrementamos contador total

	if err != nil {
		IsUpMetric.WithLabelValues(s.TargetURL).Set(0)
		log.Printf("❌ %s está DOWN", s.TargetURL)
	} else {
		defer resp.Body.Close()
		IsUpMetric.WithLabelValues(s.TargetURL).Set(1)
		latencyMetric.WithLabelValues(s.TargetURL).Observe(duration.Seconds())
		log.Printf("✅ %s está UP (%v)", s.TargetURL, duration)
	}

}
