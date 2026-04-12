# 📥 Inventory Consumer (consumer.go)

## 📌 ¿Qué hace este archivo?

Este componente escucha eventos desde RabbitMQ y ejecuta lógica de negocio.

👉 En este caso:

> Escucha `order.created` y decide si hay stock o no

---

## 🧠 Rol dentro del sistema

```text
Order crea pedido → Inventory Consumer recibe → procesa → responde
```

👉 Es el **puente entre mensajería y lógica real**

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
* `service` → lógica de negocio (stock)

👉 Separación clave:

> Consumer = recibe mensajes
> Service = decide qué hacer

---

## 🚀 StartListening() (LO IMPORTANTE)

```go
func (c *Consumer) StartListening()
```

👉 Hace 4 cosas:

1. Se suscribe a la cola
2. Escucha mensajes
3. Los valida
4. Ejecuta lógica + Ack/Nack

---

## 🔄 Flujo completo

```text
1. Consume inventory.queue
2. Llega mensaje
3. Verifica routing key
4. Parsea JSON
5. Extrae correlation_id
6. Ejecuta lógica (ReserveStock)
7. Ack o Nack
```

---

## 📡 Consumo de mensajes

```go
deliveries, err := c.client.Consume(rabbitmq.InventoryQueue, "inventory-consumer")
```

👉 Traducción:

> “escuchá todo lo que llegue a inventory.queue”

---

## 🔍 Validación de Routing Key

```go
if d.RoutingKey != rabbitmq.OrderCreatedKey
```

👉 Muy importante:

> Este consumer SOLO acepta `order.created`

Si llega otra cosa:

```text
→ lo rechaza (Nack)
→ no lo procesa
```

✔ Esto evita bugs silenciosos

---

## 📦 Parseo del evento

```go
var event events.OrderCreated
json.Unmarshal(d.Body, &event)
```

👉 Convierte JSON → struct de Go

Si falla:

```text
→ mensaje inválido
→ Nack sin retry
```

---

## 🔗 Correlation ID (tracking)

```go
if correlationID, ok := d.Headers["correlation_id"]
```

👉 Se extrae del mensaje y se pasa adelante

Sirve para:

* tracking de pedidos
* logs
* debug

---

## 🧠 Lógica de negocio

```go
c.service.ReserveStock(...)
```

👉 Acá pasa lo importante:

* verifica stock
* decide si reservar o fallar
* probablemente publica respuesta

---

## ❌ Error en lógica

```go
_ = d.Nack(false, true)
```

👉 Traducción:

> “falló → reintentá”

✔ Se usa para errores temporales (ej: DB caída)

---

## ✅ Ack (procesado OK)

```go
d.Ack(false)
```

👉 Significa:

> “todo bien, no lo vuelvas a enviar”

---

## ⚠️ Estrategia de errores

```text
Error de parseo → Nack(false, false) → descartar
Error de lógica → Nack(false, true)  → reintentar
```

👉 Esto es diseño profesional

---

## 🔄 Esquema del flujo

```text
Exchange (orders)
        ↓
inventory.queue
        ↓
[ Consumer ]
   ↓ valida routing key
   ↓ parsea JSON
   ↓ extrae correlation_id
   ↓ ejecuta ReserveStock
        ↓
   ├── OK  → Ack
   └── FAIL → Nack (retry)
```

---

## 🧪 Ejemplo real

```text
Usuario compra una notebook
↓
Order publica order.created
↓
Inventory Consumer recibe
↓
Chequea stock:
   ✔ hay → reserva
   ❌ no hay → falla
↓
Responde al sistema
```

---

## ⚠️ Cosas importantes (nivel entrevista)

* No procesar eventos inesperados
* Separar consumer de lógica de negocio
* Usar Ack/Nack manual
* Manejar retries correctamente
* Propagar correlation_id

---

## 🧠 Resumen final

Este archivo:

✔ Escucha eventos desde RabbitMQ
✔ Filtra eventos válidos
✔ Convierte JSON a objetos
✔ Ejecuta lógica de negocio
✔ Controla errores con Ack/Nack

---

## 🚀 Qué aprendiste

* Cómo consumir eventos correctamente
* Cómo validar mensajes
* Cómo manejar errores en mensajería
* Cómo conectar RabbitMQ con lógica real

---
