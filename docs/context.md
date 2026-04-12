# 🧠 Project Context

## 📌 Project Overview
Este proyecto es un sistema de software enfocado en [DESCRIBIR: e-commerce, microservicios, API, etc.].

Objetivo principal:
- [Ej: gestionar pedidos, usuarios e inventario]
- [Ej: arquitectura escalable para CV/proyecto real]

---

## 🏗️ Architecture

Tipo:
- [ ] Monolito
- [ ] Microservicios
- [ ] Modular

Descripción:
- Servicios principales:
  - auth-service → autenticación y autorización
  - product-service → gestión de productos
  - order-service → gestión de pedidos
  - inventory-service → stock

Comunicación:
- REST / gRPC / eventos (RabbitMQ, Kafka, etc.)

---

## 🧪 Development Workflow (SDD)

Este proyecto sigue **Spec-Driven Development**:

1. Definir SPEC antes de código
2. Escribir TESTS primero (TDD si aplica)
3. Implementar solución mínima
4. Refactorizar

Regla:
> Nunca escribir código sin especificación previa.

---

## 📏 Coding Standards

- Código claro y mantenible
- Evitar lógica duplicada
- Funciones pequeñas y con responsabilidad única
- Nombres descriptivos (no abreviaciones innecesarias)

---

## 🔒 Security Rules

- Validar todos los inputs
- No exponer datos sensibles
- Manejar errores correctamente
- Usar variables de entorno (.env)

---

## ⚙️ Tech Stack

- Lenguaje: [Node.js / Python / C++ / etc.]
- Framework: [Express / FastAPI / etc.]
- Base de datos: [PostgreSQL / MongoDB]
- Mensajería: [RabbitMQ / Kafka]
- Contenedores: Docker (si aplica)

---

## 📂 Important Files

- agents.md → reglas principales del agente
- context.md → memoria persistente del proyecto
- docker-compose.yml → servicios
- .env → configuración

---

## 🧠 AI Behavior Rules

El agente debe:

- Leer SIEMPRE `agents.md` y `context.md`
- Seguir el flujo SDD
- Explicar decisiones técnicas cuando sea relevante
- Priorizar código limpio y escalable
- Evitar soluciones rápidas sin justificación

---

## 🚫 Anti-Patterns

Evitar:
- Código espagueti
- Lógica en controladores
- Falta de validaciones
- Hardcoding de datos
- Ignorar errores

---

## 🎯 Current Goals

- [Ej: implementar sistema de pedidos]
- [Ej: integrar RabbitMQ]
- [Ej: mejorar arquitectura]

---

## 📌 Notes

- Este proyecto está pensado para nivel profesional (CV)
- Priorizar buenas prácticas sobre velocidad
## Perfil y Objetivos

-Quién sos: Estudiante de 3er año de Lic. en Sistemas de Información.
-Tu meta: Aprender arquitectura a fondo, entender el razonamiento detrás del código (no solo copiar y pegar) y construir un portfolio sólido para conseguir tu primer laburo como Junior.
## Cómo te gusta que laburemos

-Metodología: Siempre tengo que proponer el camino antes de ejecutar nada. Nunca, pero nunca, avanzo sin tu consentimiento.
-Pedagogía: Mi prioridad es explicarte el "por qué", las ventajas y desventajas (tradeoffs) de cada decisión, en vez de solo tirar código por tirar.
-Proactividad: Tengo la orden grabada de revisar siempre si tenés instaladas las herramientas/dependencias necesarias antes de arrancar algo nuevo.
-Resúmenes: Adoptamos el patrón de darte "bajadas a tierra" (explicaciones simples y claras) cada vez que cerramos un hito importante del proyecto.