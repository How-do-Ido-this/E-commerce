# 🔁 Outbox Relay (relay.go)

## 📌 ¿Qué hace este archivo?

Se encarga de:

> Leer eventos guardados en la base de datos (outbox) y publicarlos en RabbitMQ

👉 Es el **puente entre la base de datos y la mensajería**

---

## 🧠 Problema que resuelve

Sin relay:

```text id="1s9kz4"
Order guarda evento en DB ✔
Pero nunca llega a RabbitMQ ❌
```

👉 El sistema queda roto.

---

## ✅ Con relay

```text id="7p2xkd"
DB (outbox) → Relay → RabbitMQ → Consumers
```

👉 Garantiza que los eventos salgan sí o sí

---

## 🧩 Estructura

```go id="p9x2fj"
type Relay struct {
	db           *pgxpool.Pool
	rabbitClient *rabbitmq.Client
}
```

### ¿Qué significa?

* `db` → acceso a PostgreSQL
* `rabbitClient` → publicación en RabbitMQ

---

## 🚀 Start() (loop principal)

```go id="q7z1nt"
func (r *Relay) Start(ctx context.Context)
```

### ¿Qué hace?

```text id="8y3nvs"
Cada 1 segundo:
  → ejecuta processOutbox()
```

👉 Es un **worker en loop constante**

---

## ⏱️ Ticker

```go id="r3t6bm"
time.NewTicker(1 * time.Second)
```

👉 Traducción:

> “chequeá la outbox cada 1 segundo”

✔ Simple
✔ Funciona
❗ No es lo más eficiente (pero está bien para empezar)

---

## 🔄 Flujo completo

```text id="k5d9xp"
1. Leer eventos de la tabla outbox
2. Publicarlos en RabbitMQ
3. Si sale bien → borrarlos de la DB
```

---

## 📥 Leer eventos

```go id="n8x2kw"
SELECT id, event_type, payload FROM outbox LIMIT 10
```

👉 Trae eventos pendientes

* `id` → se usa como correlation_id
* `event_type` → routing key
* `payload` → mensaje JSON

---

## 📡 Publicar en RabbitMQ

```go id="c4y7zm"
Publish(OrdersExchange, eventType, payload, headers)
```

### Importante:

```go id="z2m1ds"
headers := map[string]interface{}{
	"correlation_id": id,
}
```

👉 Cada evento lleva su ID para tracking

---

## ❌ Error al publicar

```text id="v9k2tp"
Si falla:
→ NO se borra de la DB
→ se reintenta en el próximo ciclo
```

👉 Esto es clave:

> evita pérdida de eventos

---

## ✅ Borrar evento (solo si sale bien)

```go id="j6k2xf"
DELETE FROM outbox WHERE id = $1
```

👉 Traducción:

> “ya se publicó → no lo necesito más”

---

## 🔄 Esquema del flujo

```text id="z4m8qp"
[ DB: outbox ]
        ↓
     Relay
        ↓
 Publica en RabbitMQ
        ↓
[ Exchange: orders ]
        ↓
 Consumers (Inventory / Order)
```

---

## 🧪 Ejemplo real

```text id="x8t4pd"
Usuario compra una TV
↓
Order guarda evento en outbox
↓
Relay lo lee
↓
Lo publica como order.created
↓
Inventory lo recibe
```

---

## ⚠️ Cosas importantes (nivel entrevista)

* Nunca borrar evento si no se publicó
* El relay garantiza consistencia eventual
* Outbox + Relay evita pérdida de eventos
* Usar correlation_id para tracking
* Procesar en batches (LIMIT 10)

---

## 🧠 Resumen final

Este archivo:

✔ Lee eventos desde la DB
✔ Los publica en RabbitMQ
✔ Maneja errores sin perder datos
✔ Elimina eventos ya procesados
✔ Cierra el patrón Outbox

---

## 🚀 Qué aprendiste

* Cómo conectar DB con mensajería
* Cómo evitar inconsistencias
* Cómo implementar Outbox completo
* Cómo diseñar sistemas confiables

---
