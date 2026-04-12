# 📥 Order Consumer (consumer.go)

## 📌 ¿Qué hace este archivo?

Escucha las respuestas del servicio de Inventory y actualiza el estado de la orden.

👉 En concreto:

> Recibe `inventory.reserved` o `inventory.failed` y decide el estado final

---

## 🧠 Rol dentro del sistema

```text
Inventory → responde → Order Consumer → actualiza estado
```

👉 Es el **último paso de la SAGA**

---

## 🧩 Estructura

```go
type Consumer struct {
	client  *rabbitmq.Client
	service *Service
}
```

### ¿Qué significa?

* `client` → conexión a RabbitMQ
* `service` → lógica de negocio (actualizar orden)

---

## 🚀 StartListening()

```go
func (c *Consumer) StartListening()
```

👉 Hace:

```text
1. Se suscribe a order.queue
2. Escucha mensajes
3. Identifica tipo por RoutingKey
4. Ejecuta lógica
5. Ack o Nack
```

---

## 📡 Consumo

```go
Consume(rabbitmq.OrderQueue, "order-consumer")
```

👉 Traducción:

> “escuchá respuestas del inventory”

---

## 🔑 Uso correcto de RoutingKey (MUY IMPORTANTE)

```go
switch d.RoutingKey
```

👉 Este es el cambio clave respecto a malas prácticas:

✔ Usa metadata (routing key)
❌ No intenta adivinar por JSON

---

## 🟢 Caso 1: inventory.reserved

```go
case rabbitmq.InventoryReservedKey
```

### Flujo:

```text
1. Parsea evento InventoryReserved
2. Actualiza estado → CREATED
3. Ack
```

👉 Traducción:

> “hay stock → confirmo la orden”

---

## 🔴 Caso 2: inventory.failed

```go
case rabbitmq.InventoryFailedKey
```

### Flujo:

```text
1. Parsea evento StockInsufficient
2. Actualiza estado → FAILED
3. Ack
```

👉 Traducción:

> “no hay stock → cancelo la orden”

---

## ⚠️ Caso inesperado

```go
default:
```

```text
→ routing key desconocida
→ Nack sin retry
```

👉 Esto evita procesar basura

---

## 📦 Parseo JSON

```go
json.Unmarshal(d.Body, &event)
```

Si falla:

```text
→ mensaje inválido
→ Nack sin retry
```

---

## 🧠 Lógica de negocio

```go
UpdateOrderStatus(...)
```

👉 Cambia el estado en DB:

```text
PENDING → CREATED
PENDING → FAILED
```

---

## ❌ Error en lógica

```go
Nack(false, true)
```

👉 Traducción:

> “falló algo → reintentá”

---

## ✅ Ack final

```go
d.Ack(false)
```

👉 Solo se hace si TODO salió bien

---

## 🔄 Esquema completo del flujo

```text
Order Service
   ↓ order.created
RabbitMQ
   ↓
Inventory
   ↓ inventory.reserved / failed
RabbitMQ
   ↓
Order Consumer
   ↓
Actualiza estado en DB
```

---

## 🧪 Ejemplo real

```text
Usuario compra auriculares
↓
Order crea (PENDING)
↓
Inventory:
   ✔ hay stock → reserved
↓
Order Consumer:
   → CREATED (confirmada)
```

o

```text
Usuario compra notebook
↓
Inventory:
   ❌ no hay stock → failed
↓
Order Consumer:
   → FAILED (cancelada)
```

---

## ⚠️ Cosas importantes (nivel entrevista)

* Usar RoutingKey para identificar eventos
* No inferir tipos por JSON
* Separar consumer de lógica de negocio
* Manejar correctamente Ack/Nack
* Tener estados claros en la orden

---

## 🧠 Resumen final

Este archivo:

✔ Escucha respuestas de otros servicios
✔ Identifica eventos correctamente
✔ Actualiza el estado de la orden
✔ Cierra la SAGA
✔ Maneja errores con Ack/Nack

---

## 🚀 Qué aprendiste

* Cómo cerrar una transacción distribuida
* Cómo manejar estados finales
* Cómo reaccionar a eventos externos
* Cómo hacer un flujo completo con RabbitMQ

---
