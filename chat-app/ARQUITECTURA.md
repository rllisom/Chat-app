# Chat App — Guía de Arquitectura

## Índice
1. [Visión general](#1-visión-general)
2. [Arquitectura de microservicios](#2-arquitectura-de-microservicios)
3. [Flujo de una petición](#3-flujo-de-una-petición)
4. [Comunicación entre servicios: gRPC](#4-comunicación-entre-servicios-grpc)
5. [Estructura de carpetas](#5-estructura-de-carpetas)
6. [Archivos — qué hace cada uno y por qué existe](#6-archivos--qué-hace-cada-uno-y-por-qué-existe)
7. [Base de datos MongoDB](#7-base-de-datos-mongodb)
8. [Cómo arrancar el proyecto](#8-cómo-arrancar-el-proyecto)

---

## 1. Visión general

Es una aplicación de chat en tiempo real construida con **3 procesos independientes** que se comunican entre sí:

```
Frontend
   │  HTTP REST
   ▼
┌─────────────────────────────┐
│   Gateway  :8080            │  ← Único punto de entrada público
│   (Gin HTTP)                │
└────────────┬────────────────┘
             │  gRPC
      ┌──────┴──────┐
      ▼             ▼
┌──────────┐  ┌──────────┐
│  User    │  │  Chat    │
│ Service  │  │ Service  │
│  :9001   │  │  :9002   │
└────┬─────┘  └────┬─────┘
     ▼             ▼
  MongoDB       MongoDB
```

**Tecnologías usadas:**
| Tecnología | Para qué |
|---|---|
| Go | Lenguaje principal |
| Gin | Framework HTTP para el Gateway |
| gRPC + Protocol Buffers | Comunicación entre microservicios |
| MongoDB | Base de datos |
| WebSocket (Gorilla) | Chat en tiempo real (Hub) |

---

## 2. Arquitectura de microservicios

El proyecto tiene **3 ejecutables** (cada uno con su propio `main`):

### Gateway (`cmd/gateway/`)
- Escucha peticiones HTTP en el puerto **8080**
- Es el único servicio visible para el frontend
- No tiene lógica de negocio ni accede a MongoDB directamente
- Traduce peticiones HTTP → llamadas gRPC a los servicios internos

### User Service (`cmd/user-service/`)
- Servidor gRPC en el puerto **9001**
- Gestiona todo lo relacionado con usuarios: crear, buscar
- Tiene su propia conexión a MongoDB
- Solo habla gRPC, nunca recibe peticiones HTTP directamente

### Chat Service (`cmd/chat-service/`)
- Servidor gRPC en el puerto **9002**
- Gestiona salas y mensajes
- Tiene su propia conexión a MongoDB
- Solo habla gRPC

---

## 3. Flujo de una petición

### Ejemplo: Crear un usuario

```
1. Frontend → POST /api/v1/users  { username: "raul", email: "raul@email.com" }
                    │
                    ▼
2. Gateway recibe la petición HTTP
   user/grpc_handler.go → CreateUser()
                    │
                    ▼
3. Gateway llama al User Service vía gRPC
   pb.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{...})
                    │
         [viaje por red, protocolo gRPC]
                    │
                    ▼
4. User Service recibe la llamada gRPC
   userservice/server.go → CreateUser()
                    │
                    ▼
5. Ejecuta lógica de negocio
   user/service.go → CreateUser() → valida username, email, duplicados...
                    │
                    ▼
6. Persiste en MongoDB
   user/repository.go → Create() → InsertOne()
                    │
                    ▼
7. Respuesta viaja de vuelta: MongoDB → repository → service → server gRPC → Gateway → HTTP → Frontend
```

---

## 4. Comunicación entre servicios: gRPC

### ¿Qué es gRPC?
gRPC es un protocolo de comunicación entre procesos. Es más rápido que REST porque usa formato binario (Protocol Buffers) en lugar de JSON.

### El contrato: archivos `.proto`

Los archivos `.proto` definen **qué métodos existen y qué datos intercambian**. Son el "contrato" entre Gateway y servicios.

**`proto/user.proto`** define:
```
CreateUser  → recibe (username, email)    → devuelve UserResponse
GetUserById → recibe (id)                 → devuelve UserResponse
GetAllUsers → recibe ()                   → devuelve lista de UserResponse
```

**`proto/chat.proto`** define:
```
CreateRoom       → recibe (name, created_by) → devuelve RoomResponse
GetAllRooms      → recibe ()                 → devuelve lista de RoomResponse
SendMessage      → recibe (room_id, user_id, username, content) → devuelve MessageResponse
GetRoomMessages  → recibe (room_id)          → devuelve lista de MessageResponse
```

### Archivos generados automáticamente
Al ejecutar `protoc`, se generan 2 archivos Go por cada `.proto`:

| Archivo generado | Contenido |
|---|---|
| `user.pb.go` | Structs de los mensajes (CreateUserRequest, UserResponse...) |
| `user_grpc.pb.go` | Interfaz del servidor + código del cliente gRPC |

**Estos archivos nunca se editan a mano.**

### Roles del cliente y el servidor gRPC

```
Gateway                          User Service
──────────────────               ────────────────────────────
userservice/client.go            userservice/server.go
  └─ Client.GRPCClient()           └─ Server.CreateUser()
       │                                 │
       │   pb.UserServiceClient          │   implementa
       └──────────────────────────────── pb.UserServiceServer
```

- **`client.go`** — envuelve la conexión gRPC y expone métodos simples al Gateway
- **`server.go`** — recibe las llamadas gRPC y las delega al service de Go

---

## 5. Estructura de carpetas

```
chat-app/
│
├── proto/                        ← Contratos gRPC
│   ├── user.proto                  Definición del servicio de usuarios
│   ├── chat.proto                  Definición del servicio de chat
│   ├── user.pb.go                  Generado: structs de mensajes
│   ├── user_grpc.pb.go             Generado: cliente y servidor gRPC
│   ├── chat.pb.go                  Generado: structs de mensajes
│   └── chat_grpc.pb.go             Generado: cliente y servidor gRPC
│
├── internal/
│   ├── db/
│   │   └── mongo.go              ← Conexión global a MongoDB
│   │
│   ├── user/                     ← Dominio: lógica de usuarios
│   │   ├── user.go                 Struct User + constructor NewUser
│   │   ├── user_repository.go      Acceso a MongoDB (CRUD)
│   │   ├── user_service.go         Lógica de negocio y validaciones
│   │   ├── grpc_handler.go         Handler HTTP del Gateway para /users
│   │   └── service_test.go         Tests del servicio
│   │
│   ├── chat/                     ← Dominio: lógica de chat
│   │   ├── chat.go                 Structs Room, Message + constructores
│   │   ├── chat_repository.go      Acceso a MongoDB (salas y mensajes)
│   │   ├── chat_service.go         Lógica de negocio y validaciones
│   │   └── grpc_handler.go         Handler HTTP del Gateway para /rooms
│   │
│   ├── userservice/              ← Adaptador gRPC para usuarios
│   │   ├── server.go               Servidor gRPC: traduce proto → user domain
│   │   └── client.go               Cliente gRPC: usado por el Gateway
│   │
│   └── chatservice/              ← Adaptador gRPC para chat
│       ├── server.go               Servidor gRPC: traduce proto → chat domain
│       └── client.go               Cliente gRPC: usado por el Gateway
│
└── cmd/                          ← Ejecutables (un main por servicio)
    ├── gateway/
    │   └── main.go               ← Arranca el servidor HTTP :8080
    ├── user-service/
    │   └── main.go               ← Arranca el servidor gRPC :9001
    └── chat-service/
        └── main.go               ← Arranca el servidor gRPC :9002
```

---

## 6. Archivos — qué hace cada uno y por qué existe

### `internal/db/mongo.go`
**Qué hace:** Gestiona la conexión global a MongoDB.

**Por qué existe:** Centraliza la conexión para que todos los repositorios usen el mismo cliente en lugar de crear una conexión nueva cada vez.

```
db.Connect("mongodb://localhost:27017")  → abre la conexión
db.GetCollection("chat_app", "users")   → devuelve una colección para operar
db.Disconnect()                          → cierra la conexión al apagar
```

---

### `internal/user/user.go`
**Qué hace:** Define el struct `User` y su constructor `NewUser`.

**Por qué existe:** Es el modelo de datos central del dominio de usuarios. Todos los demás archivos del paquete `user` trabajan con este struct.

```go
User {
    ID        // generado automáticamente
    Username
    Email
    CreatedAt // asignado al crear
}
```

---

### `internal/user/user_repository.go`
**Qué hace:** Operaciones de base de datos para usuarios (Create, FindByID, FindByUsername, FindAll).

**Por qué existe:** Separa el acceso a datos de la lógica de negocio. Si mañana cambias MongoDB por PostgreSQL, solo tocas este archivo.

```
Create()         → InsertOne en MongoDB
FindByID()       → FindOne con filtro por _id
FindByUsername() → FindOne con filtro por username
FindAll()        → Find con cursor
```

---

### `internal/user/user_service.go`
**Qué hace:** Lógica de negocio para usuarios: validaciones, reglas.

**Por qué existe:** Es la capa que decide **si** una operación puede hacerse. El repositorio no valida nada, solo ejecuta. El service valida primero.

```
CreateUser() → valida username no vacío, mínimo 3 chars,
               email no vacío y con @,
               username no duplicado,
               luego llama al repositorio
```

---

### `internal/user/grpc_handler.go`
**Qué hace:** Handler HTTP de Gin para las rutas `/api/v1/users` en el Gateway.

**Por qué existe:** Es el punto de entrada HTTP para el dominio de usuarios en el Gateway. Recibe la petición HTTP, la transforma en una llamada gRPC al User Service, y devuelve la respuesta JSON.

```
POST /users     → CreateUser  → gRPC → User Service
GET  /users     → GetAllUsers → gRPC → User Service
GET  /users/:id → GetUserByID → gRPC → User Service
```

---

### `internal/chat/chat.go`
**Qué hace:** Define los structs `Room` y `Message` y sus constructores.

**Por qué existe:** Igual que `user.go` pero para el dominio de chat. Separa la definición del modelo del resto de la lógica.

---

### `internal/chat/chat_repository.go`
**Qué hace:** Operaciones de base de datos para salas y mensajes.

**Por qué existe:** Misma razón que `user_repository.go`. Tiene dos colecciones: `rooms` y `messages`.

```
CreateRoom()          → InsertOne en rooms
FindRoomByID()        → FindOne en rooms
FindAllRooms()        → Find en rooms con cursor
SaveMessage()         → InsertOne en messages
FindMessagesByRoom()  → Find en messages filtrado por room_id,
                        ordenado por sent_at DESC, limitado a 50
```

---

### `internal/chat/chat_service.go`
**Qué hace:** Lógica de negocio para salas y mensajes: validaciones.

**Por qué existe:** Igual que `user_service.go`. Valida antes de persistir.

```
CreateRoom()     → nombre no vacío, mínimo 3 chars, ID de usuario válido
SendMessage()    → contenido no vacío, máximo 500 chars, IDs válidos
GetRoomMessages() → últimos 50 mensajes de una sala
```

---

### `internal/chat/grpc_handler.go`
**Qué hace:** Handler HTTP de Gin para `/api/v1/rooms` en el Gateway.

**Por qué existe:** Igual que `user/grpc_handler.go` pero para el dominio de chat.

```
POST /rooms              → CreateRoom
GET  /rooms              → GetAllRooms
POST /rooms/:id/messages → SendMessage
GET  /rooms/:id/messages → GetRoomMessages
```

---

### `internal/userservice/server.go`
**Qué hace:** Implementa la interfaz gRPC `UserServiceServer` generada por protoc.

**Por qué existe:** Es el puente entre el mundo gRPC (proto) y el mundo Go (domain). Recibe llamadas gRPC con tipos proto (`*pb.CreateUserRequest`), los convierte a tipos Go, llama al service, y convierte la respuesta de vuelta a proto.

```
gRPC request (proto) → server.go → user.UserService → respuesta → proto response
```

La función `toProto()` convierte un `*user.User` en un `*pb.UserResponse`.

---

### `internal/userservice/client.go`
**Qué hace:** Cliente gRPC del User Service. Abre la conexión TCP y expone métodos para llamar al servidor.

**Por qué existe:** El Gateway necesita este cliente para comunicarse con el User Service. `GRPCClient()` expone el cliente proto interno para pasárselo al handler HTTP.

```
NewClient("localhost:9001") → abre conexión gRPC
GRPCClient()                → devuelve pb.UserServiceClient para el Gateway
CreateUser()                → wrapper que añade contexto con timeout
```

---

### `internal/chatservice/server.go` y `client.go`
**Qué hacen:** Mismo rol que `userservice/server.go` y `client.go` pero para el Chat Service.

Los helpers `roomToProto()` y `messageToProto()` convierten los structs Go a respuestas proto.

---

### `cmd/user-service/main.go`
**Qué hace:** Arranca el User Service como proceso independiente.

**Por qué existe:** Es el punto de entrada del microservicio de usuarios.

```
1. Conecta a MongoDB
2. Crea repo → service → grpcServer (capas en orden)
3. Crea servidor gRPC y registra el servicio
4. Abre puerto TCP :9001
5. Escucha conexiones gRPC indefinidamente
```

---

### `cmd/chat-service/main.go`
**Qué hace:** Arranca el Chat Service. Mismo patrón que user-service pero en el puerto **:9002**.

---

### `cmd/gateway/main.go`
**Qué hace:** Arranca el Gateway HTTP.

**Por qué existe:** Es el único proceso que el frontend conoce.

```
1. Conecta al User Service gRPC en :9001
2. Conecta al Chat Service gRPC en :9002
3. Crea los handlers HTTP pasándoles los clientes gRPC
4. Registra las rutas en Gin bajo /api/v1
5. Escucha peticiones HTTP en :8080
```

---

### `proto/user.proto` y `proto/chat.proto`
**Qué hacen:** Definen el contrato de comunicación entre Gateway y servicios.

**Por qué existen:** Son la fuente de verdad sobre qué métodos existen y qué tipos de datos se intercambian. A partir de ellos `protoc` genera todo el código de comunicación gRPC automáticamente.

---

## 7. Base de datos MongoDB

El proyecto usa **una sola instancia** de MongoDB con la base de datos `chat_app` y tres colecciones:

| Colección | Gestionada por | Documentos |
|---|---|---|
| `users` | User Service | `{ _id, username, email, created_at }` |
| `rooms` | Chat Service | `{ _id, name, created_by, created_at }` |
| `messages` | Chat Service | `{ _id, room_id, user_id, username, content, sent_at }` |

Cada servicio abre su propia conexión a MongoDB de forma independiente.

---

## 8. Cómo arrancar el proyecto

Necesitas **3 terminales**, una por servicio:

```bash
# Terminal 1 — User Service
cd chat-app
go run cmd/user-service/main.go

# Terminal 2 — Chat Service
go run cmd/chat-service/main.go

# Terminal 3 — Gateway
go run cmd/gateway/main.go
```

MongoDB debe estar corriendo antes de arrancar cualquier servicio:
```bash
# Windows (si está instalado como servicio)
net start MongoDB

# O manualmente
"C:/Program Files/MongoDB/Server/8.2/bin/mongod.exe" --dbpath "C:/data/db"
```

**Endpoints disponibles una vez todo arrancado:**
```
POST   http://localhost:8080/api/v1/users
GET    http://localhost:8080/api/v1/users
GET    http://localhost:8080/api/v1/users/:id

POST   http://localhost:8080/api/v1/rooms
GET    http://localhost:8080/api/v1/rooms
POST   http://localhost:8080/api/v1/rooms/:id/messages
GET    http://localhost:8080/api/v1/rooms/:id/messages
```
