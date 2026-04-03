# ARQUITECTURA_RESUMIDA.md

## Objetivo del Proyecto
E-commerce robusto basado en microservicios, diseñado para ser escalable, resiliente y consistente.

## Arquitectura
- **Patrón:** Event-Driven Microservices (Arquitectura basada en eventos).
- **Comunicación:** RabbitMQ (Direct Exchange) para desacoplamiento total.
- **Consistencia:** Patrón SAGA con orquestación (gestionado por el servicio de `Order`).
- **Stack:** Go (Backend), RabbitMQ (Docker).

## Estado del Arte (Fase 1 completada)
1.  **Infraestructura:** `docker-compose.yaml` (RabbitMQ + management plugin) listo.
2.  **Scaffolding:** Estructura de carpetas `/services/*` y `/pkg/events` creada.
3.  **Contratos:** Definición de esquemas compartidos (`OrderCreated`, `InventoryReserved`, `StockInsufficient`, `OrderCancelled`) en `/pkg/events/events.go`.

## Roadmap (Fases pendientes)
- **Fase 2:** Implementación de servicios Core (`Order`, `Inventory`, `User`, `Product`).
- **Fase 3:** Implementación de productores/consumidores de eventos y lógica de SAGA.
- **Fase 4:** Testing e integración (Testcontainers).

## Decisiones Técnicas Clave
- RabbitMQ sobre Kafka (simplicidad inicial).
- Saga Orchestrator para simplificar debugueo.
- Go para eficiencia y concurrencia.

---
*Este documento es la "Biblia" de contexto. Si sos el agente nuevo y estás leyendo esto, cargalo en tu memoria de corto plazo.*