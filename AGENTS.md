# AGENTS.md - Normas del Proyecto

Este proyecto es un entorno de aprendizaje profesional. Todo el trabajo sigue estos principios:

## 1. Filosofía de Trabajo
- **Conceptos > Código:** Priorizar el entendimiento de patrones (Saga, Event-Driven, Microservicios) sobre la escritura rápida.
- **Transparencia:** Ante cualquier decisión técnica, presentar al menos dos alternativas (pros/contras) antes de implementar.
- **Aprendizaje:** Si el usuario no entiende un concepto, detener el flujo y explicarlo usando analogías o ejemplos sencillos.

## 2. Estándares Técnicos
- **Lenguaje:** Go.
- **Commits:** Convencional (ej: `feat: add order creation saga`, `fix: handle inventory timeout`).
- **Arquitectura:** Hexagonal / Clean Architecture simplificada (por carpetas de dominio: `internal/order`, `internal/inventory`, etc.).
- **Testing:** Todo componente crítico debe tener tests unitarios.

## 3. Flujo de Trabajo (Protocolo de Seguridad)
- **NO AVANZAR SIN CONSENTIMIENTO:** El agente propone pasos. El usuario debe confirmar.
- **Documentación Viva:** Mantener `ARQUITECTURA_RESUMIDA.md` actualizado ante cambios estructurales.
