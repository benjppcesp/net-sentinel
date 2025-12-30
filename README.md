# 🛡️ Net-Sentinel: Network Monitoring Probe

![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Grafana](https://img.shields.io/badge/Grafana-F46800?style=for-the-badge&logo=Grafana&logoColor=white)

**Net-Sentinel** es una sonda de monitoreo de red desarrollada en **Go** diseñada para la observabilidad de servicios críticos. Realiza chequeos de disponibilidad y latencia, exponiendo métricas nativas para **Prometheus** y visualización profesional en **Grafana**.



---

## 🚀 Características Principales
* **Monitoreo en Tiempo Real:** Seguimiento de latencia (ms) y estado de disponibilidad (UP/DOWN).
* **Arquitectura Concurrente:** Uso de Goroutines para evitar bloqueos entre el monitoreo y la exposición de métricas.
* **Cloud Native:** Totalmente contenedorizado con Docker y orquestado con Docker Compose.
* **Métricas Estándar:** Exportación de métricas compatibles con el ecosistema de Prometheus.

---

## 🏗️ Arquitectura Técnica

### Concurrencia y Seguridad
El núcleo utiliza un modelo de **Goroutines independientes**:
1.  **Sonda (Probe):** Ejecuta chequeos asíncronos mediante un `time.Ticker`.
2.  **Servidor de Métricas:** Un servidor HTTP dedicado en el puerto `:2112` expone los datos para el scraping de Prometheus.

### Stack Tecnológico
* **Lenguaje:** Go 1.23+
* **Métricas:** Prometheus Client Golang
* **Infraestructura:** Docker, Docker Compose
* **Visualización:** Grafana 10.x

---

## 📊 Observabilidad

El proyecto incluye un dashboard preconfigurado. Para usarlo:
1.  Importa el archivo JSON ubicado en `/grafana/dashboards/net-sentinel.json`.
2.  Conecta con el Data Source de Prometheus (`http://prometheus:9090`).

### Métricas Clave Expuestas:
* `net_sentinel_http_success`: `1` si el objetivo es alcanzable, `0` si falla.
* `net_sentinel_http_duration_seconds`: Latencia de la petición HTTP en segundos.
* `go_goroutines`: Cantidad de hilos lógicos en ejecución.

---

## 🛠️ Instalación y Despliegue

### Requisitos Previos
* Docker y Docker Compose instalados.
* Archivo `.env` configurado (ver `.env.example`).

### Pasos para iniciar:
```bash
# 1. Clonar el repositorio
git clone [https://github.com/benjppcesp/net-sentinel.git](https://github.com/benjppcesp/net-sentinel.git)
cd net-sentinel

# 2. Configurar variables de entorno
cp .env.example .env

# 3. Desplegar el stack completo (App + Prometheus + Grafana)
docker-compose up -d

<p align="center">
  <img src="./assets/dashboard-preview.png" width="900" alt="Grafana Dashboard Preview">
</p>
