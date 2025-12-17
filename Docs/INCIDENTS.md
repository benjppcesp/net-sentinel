# Incidents Log

-----

## INCIDENT-001 — Docker Compose crash (ContainerConfig)
**Date:** 17-12-2025  
**Severity:** Medium  
**Environment:** Local (WSL + Docker Desktop)  

### Summary
Docker Compose failed to recreate containers with a `KeyError: 'ContainerConfig'`, blocking the local stack startup.

### Impact
- Stack could not be started
- Development blocked
- Time lost debugging non-application issues

### Root Cause
Use of deprecated `docker-compose v1 (1.29.2, Python)` with a modern Docker Engine, causing metadata incompatibility during container recreation.

### Trigger
Changes to container configuration (logging, volumes, resources) followed by `docker-compose up -d`.

### Resolution
1. Full cleanup of containers, volumes, and networks:
   ```bash
   docker-compose down -v
