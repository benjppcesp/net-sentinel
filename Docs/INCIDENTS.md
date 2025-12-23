# Incidents Log

-----

## INCIDENT-001 — Docker Compose image pull failure (invalid reference / not found)
**Date:** 17-12-2025  
**Severity:** Medium  
**Environment:** Local (WSL + Docker Desktop)  

### Summary
Docker Compose failed to pull and start the `net-sentinel` service due to an invalid image reference and later a non-existent image tag, preventing the stack from starting.

### Impact
- Docker Compose stack could not be started
- Local development environment blocked
- Time lost debugging Docker and Compose configuration
- Prometheus and Grafana unavailable during incident

### Root Cause
The incident had **two chained root causes**:

1. **Undefined environment variable (`GITHUB_USER`)**  
   Docker Compose attempted to resolve the image name using an empty variable, resulting in an invalid image reference:


2. **Non-existent image tag (`v1.0.0`)**  
After fixing the variable resolution, Docker Compose attempted to pull an image tag that was not published in GitHub Container Registry (GHCR).

### Trigger
- Introduction of variable-based image definition in `docker-compose.yml`
- Execution of `docker compose pull` and `docker compose up -d` without verifying:
- Variable availability
- Existing image tags in the registry

### Detection
- Docker Compose warning
- Docker daemon error
- Later error


### Resolution
1. Defined required environment variables using a `.env` file:
 ```env
 GITHUB_USER=benjppcesp
 IMAGE_TAG=main-032fdb9
```
2. Updated docker-compose.yml to validate required variables:
image: ghcr.io/${GITHUB_USER:?GITHUB_USER no definido}/net-sentinel:${IMAGE_TAG:-latest}

3. Verified existing image tags locally and in GHCR:
```
docker images | grep net-sentinel
```
4. Pulled the correct image tag and restarted the stack:
```
docker compose pull net-sentinel
docker compose up -d
```
5. Verified that the stack started successfully:
```
docker ps --filter name=net-sentinel
```

