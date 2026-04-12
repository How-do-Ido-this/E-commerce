# 🐇 RabbitMQ Client (client.go)

## 📌 ¿Qué hace este archivo?

Este archivo es el **punto central de comunicación con RabbitMQ**.

👉 Se encarga de:

* Conectarse a RabbitMQ
* Crear exchange y colas
* Definir reglas de enrutamiento (bindings)
* Publicar mensajes
* Consumir mensajes

En pocas palabras:

> Es el “cerebro” de la mensajería del sistema

---

## 🧠 Conceptos clave (en criollo)

* **Exchange (`orders`)** → Recepcionista
* **Routing Key (`order.created`)** → Motivo del mensaje
* **Queue (`inventory.queue`)** → Bandeja donde cae el mensaje
* **Bind** → Regla: “si viene X → mandalo a Y”

---

## 🧩 Constantes importantes

```go
const (
	OrdersExchange       = "orders"
	OrderCreatedKey      = "order.created"
	InventoryReservedKey = "inventory.reserved"
	InventoryFailedKey   = "inventory.failed"
)
```

👉 Esto define el **lenguaje de comunicación entre servicios**.

Ejemplo:

* Order crea → `order.created`
* Inventory responde → `inventory.reserved` o `inventory.failed`

---

## 🔌 Conexión a RabbitMQ

```go
func NewClient(amqpURL string) (*Client, error)
```

### ¿Qué hace?

1. Se conecta a RabbitMQ
2. Abre un canal

### ¿Por qué es importante?

RabbitMQ funciona así:

* **Connection** = conexión TCP (cara)
* **Channel** = canal lógico (barato)

👉 Siempre trabajás con channels, no directo con la conexión.

---

## 🏗️ SetupTopology() (LO MÁS IMPORTANTE)

```go
func (c *Client) SetupTopology() error
```

### ¿Qué hace?

Configura TODO:

1. Crea el exchange `orders`
2. Crea `inventory.queue`
3. La conecta a `order.created`
4. Crea `order.queue`
5. La conecta a:

   * `inventory.reserved`
   * `inventory.failed`

---

## 🔄 Flujo real del sistema

```text
Order Service
   ↓ order.created
[ Exchange: orders ]
   ↓
[ inventory.queue ]
   ↓
Inventory Service
   ↓ inventory.reserved / inventory.failed
[ order.queue ]
   ↓
Order Service
```

👉 Esto es tu SAGA funcionando.

---

## 📦 Exchange

```go
ExchangeDeclare(...)
```

* Tipo: `direct`
* Durable: `true`

### Traducción:

> El exchange queda guardado aunque reinicies RabbitMQ
> y enruta por routing key exacta

---

## 📥 Inventory Queue

```go
QueueBind(InventoryQueue, OrderCreatedKey, OrdersExchange)
```

👉 Regla:

> Todo lo que sea `order.created` → va a `inventory.queue`

---

## 📤 Order Queue

```go
for _, routingKey := range []string{
	InventoryReservedKey,
	InventoryFailedKey,
}
```

👉 Regla:

> Todo lo que sea respuesta de inventory → vuelve a Order

---

## 📡 Publish (enviar mensajes)

```go
func (c *Client) Publish(...)
```

### ¿Qué hace?

Envía mensajes al exchange con:

* routing key
* body (JSON)
* headers

### Ejemplo real:

```text
Order crea pedido
→ Publish("orders", "order.created")
```

---

## 🔗 Headers (correlation_id)

```go
Headers: amqp.Table(headers)
```

👉 Sirve para:

* seguir un pedido entre servicios
* debuggear
* logging

Ejemplo:

```text
Pedido #123 → viaja por todo el sistema con el mismo ID
```

---

## 📥 Consume (recibir mensajes)

```go
func (c *Client) Consume(...)
```

### Clave:

```go
auto-ack = false
```

👉 Traducción:

> Vos decidís cuándo el mensaje está procesado

---

## ✅ Ack vs ❌ Nack

* `Ack` → procesado OK
* `Nack(false, false)` → error, descartar
* `Nack(false, true)` → error, reintentar

👉 Esto evita:

* perder mensajes
* duplicar procesos

---

## 🔚 Close()

```go
func (c *Client) Close()
```

👉 Cierra:

* canal
* conexión

Importante para liberar recursos.

---

## ⚠️ Cosas importantes (nivel entrevista)

* Si no declarás exchange → los mensajes se pierden
* Nunca usar auto-ack en producción
* Usar routing key, no inferir por JSON
* Usar correlation_id en sistemas distribuidos

---

## 🧪 Ejemplo real (e-commerce)

```text
Usuario compra zapatillas
↓
Order crea pedido
↓
Inventory verifica stock
↓
Si hay:
   → inventory.reserved
Si no hay:
   → inventory.failed
↓
Order actualiza estado
```

---

## 🧠 Resumen final

Este archivo:

✔ Define cómo se comunican los servicios
✔ Garantiza que los mensajes lleguen correctamente
✔ Permite manejar errores con Ack/Nack
✔ Implementa la base de una SAGA

---

## 🚀 Qué aprendiste

* Cómo funciona RabbitMQ en la práctica
* Cómo enrutar mensajes correctamente
* Cómo hacer un sistema desacoplado
* Cómo manejar consistencia en microservicios

---

                ┌───────────────────────────┐
                │       Client (Go)         │
                │---------------------------│
                │ conn    (*Connection)     │
                │ channel (*Channel)        │
                └────────────┬──────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
 NewClient()         SetupTopology()        Close()
 (conexión)          (infraestructura)      (cleanup)
        │                    │
        │                    ▼
        │          ┌───────────────────────┐
        │          │     RabbitMQ          │
        │          └─────────┬─────────────┘
        │                    │
        ▼                    ▼
 Publish()              Consume()
 (envía)                (recibe)
---


              📤 PUBLISH
 Order Service ───────────────────────────────┐
   (order.created)                            │
                                              ▼
                                  ┌────────────────────┐
                                  │ Exchange: orders   │
                                  │ (tipo: direct)     │
                                  └─────────┬──────────┘
                                            │
                 ┌──────────────────────────┴──────────────────────────┐
                 │                                                     │
                 ▼                                                     ▼
     inventory.queue                                      order.queue
 (escucha order.created)                (escucha inventory.reserved / failed)
                 │                                                     ▲
                 ▼                                                     │
        Inventory Service ──────────────── Publish ─────────────────────┘
                       (inventory.reserved / inventory.failed)