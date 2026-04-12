# 📦 Order Service (service.go)

## 📌 ¿Qué hace este archivo?

Define la **lógica de negocio de las órdenes**.

👉 En concreto:

* Crea órdenes
* Las guarda en base de datos
* Dispara eventos (usando Outbox)
* Actualiza el estado según respuestas

---

## 🧠 Rol dentro del sistema

```text id="nq8k2g"
Cliente → Order Service → DB + Evento → Inventory → Respuesta → Order Service
```

👉 Es el **orquestador inicial del flujo SAGA**

---

## 🧩 Estructura

```go id="6q7bcs"
type Service struct {
	repo         Repository
	rabbitClient *rabbitmq.Client
}
```

### ¿Qué significa?

* `repo` → acceso a base de datos
* `rabbitClient` → comunicación con RabbitMQ

👉 Importante:

> Este service NO habla directo con RabbitMQ en CreateOrder
> usa Outbox (mucho mejor)

---

## 🚀 CreateOrder() (LO MÁS IMPORTANTE)

```go id="q7i3dz"
func (s *Service) CreateOrder(ctx context.Context, o *Order) error
```

---

## 🔄 Flujo completo

```text id="l4u1iy"
1. Genera ID si no existe
2. Setea estado = PENDING
3. Crea evento OrderCreated
4. Guarda en DB + Outbox (transacción)
```

---

## 🆔 Generación de ID

```go id="d5k2l8"
o.ID = uuid.New().String()
```

👉 Cada orden tiene un identificador único global

✔ Evita colisiones
✔ Funciona en sistemas distribuidos

---

## 📊 Estado inicial

```go id="z2p8xm"
o.Status = StatusPending
```

👉 Traducción:

> “todavía no sé si hay stock”

---

## 📦 Creación del evento

```go id="v9c3ws"
event := events.OrderCreated{...}
```

👉 Esto representa:

```text id="vkn9b3"
“Se creó una orden y necesita validación de stock”
```

Incluye:

* OrderID
* UserID
* Items
* Cantidad

---

## 💣 Outbox Pattern (MUY IMPORTANTE)

```go id="r3x7jp"
SaveTransactional(ctx, o, "order.created", event)
```

👉 Esto es clave. Traducción real:

> “guardo la orden Y el evento en la misma transacción”

---

## 🧠 ¿Por qué es importante?

Problema clásico:

```text id="g7m8xa"
Guardar en DB ✔
Publicar evento ❌ (falló)
→ sistema inconsistente
```

---

### ✔ Con Outbox:

```text id="p4z2ym"
Guardar orden + evento juntos ✔
Luego otro proceso (relay) publica
```

👉 Resultado:

> nunca perdés eventos

---

## 🔄 Flujo con Outbox

```text id="8z2xkm"
Order Service
   ↓
[ DB Transaction ]
   ├── guarda orden
   └── guarda evento (outbox)
   ↓
Relay
   ↓
RabbitMQ (order.created)
```

---

## 🔁 UpdateOrderStatus()

```go id="c1k9mz"
func (s *Service) UpdateOrderStatus(...)
```

👉 Sirve para:

```text id="xq3t9p"
Inventory responde → Order actualiza estado
```

Ejemplo:

* `PENDING` → `CONFIRMED`
* `PENDING` → `FAILED`

---

## 🔄 Flujo completo SAGA (con este service)

```text id="kq2b7z"
1. Usuario compra
2. CreateOrder → guarda + evento
3. Relay publica order.created
4. Inventory procesa
5. Inventory responde
6. Order actualiza estado
```

---

## 🧪 Ejemplo real

```text id="w9f3yt"
Usuario compra un celular
↓
Order Service:
   → crea orden (PENDING)
   → guarda evento
↓
Inventory:
   → verifica stock
↓
Si hay:
   → CONFIRMED
Si no hay:
   → FAILED
```

---

## ⚠️ Cosas importantes (nivel entrevista)

* Usar Outbox para evitar inconsistencias
* No publicar eventos directo desde lógica de negocio
* Manejar estados explícitos (PENDING, CONFIRMED, FAILED)
* Usar UUID en sistemas distribuidos

---

## 🧠 Resumen final

Este archivo:

✔ Maneja la creación de órdenes
✔ Usa transacciones para consistencia
✔ Implementa Outbox Pattern
✔ Inicia el flujo SAGA
✔ Actualiza estados según eventos

---

## 🚀 Qué aprendiste

* Cómo iniciar una SAGA correctamente
* Cómo evitar pérdida de eventos
* Cómo diseñar lógica desacoplada
* Cómo manejar estados en microservicios

---
