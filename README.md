# Laboratorio 1 — Sistema Híbrido de Registro, Autenticación, Presencia y Mensajería

## Integrantes

| Nombre completo | Rol USM |
|---|---|
| Héctor Muñoz Vonk | 202373628-7 |
| Lucas Morrison | 202273547-3 |

**Grupo N°:** 11

## Puertos por defecto

| Servicio | Protocolo | Puerto |
|---|---|---|
| Registro / Historial | HTTP (sobre TCP crudo) | 8080  |
| Autenticación / Mensajería | TCP | 9000 |
| Heartbeat | UDP | 9001 |

## Instrucciones de ejecución

Requiere Python 3.10+ (lado servidor) y Go 1.20+ (lado cliente).

Desde la carpeta `Servidor/`:

```bash
# Terminal 1 — servidor HTTP (registro + historial)
python servidor.py

# Terminal 2 — servidor TCP + UDP (login, mensajería, heartbeat, watchdog)
python main.py
```

Ambos deben quedar corriendo simultáneamente. `servidor.py` escucha en el puerto 8888; `main.py` levanta el servidor TCP en 9000 y el servidor UDP en 9001 dentro del mismo proceso, para compartir el estado de sesiones en memoria.

Desde la carpeta `Cliente/`:

Para conectar al endpoint HTTP para crear un usuario:

go run registrar_usuario.go

Para conectar a los endpoint TCP y UDP de comunicación:

go run cliente_tcp.go

## Archivos generados

- `usuarios.csv` — credenciales registradas (`username,password,fecha_registro`)
- `sesiones.csv` — sesiones activas (`token,username,timestamp_creacion,timestamp_ultimo_heartbeat,estado`)
- `historial.csv` — historial de mensajes (`timestamp,username,mensaje`)

Todos los archivos se generan automáticamente al iniciar los servidores si no existen, y se protegen con locks para evitar condiciones de carrera ante accesos concurrentes.

## Protocolo HTTP

### `POST /register`

Body en JSON:
```json
{"username": "alice", "password": "hunter2"}
```

Respuestas:
- `201 Created` — registro exitoso
- `409 Conflict` — el username ya existe
- `400 Bad Request` — falta username o password, o el body no es JSON válido

### `GET /history`

Devuelve el contenido de `historial.csv` en texto plano, una línea por mensaje.

Respuestas:
- `200 OK`

## Protocolo TCP (puerto 9000)

Todos los comandos y respuestas terminan en `\n`. Codificación UTF-8.

### `LOGIN <username> <password>`

Respuestas:
- `OK <token> <puerto_udp>` — login exitoso, entrega token de sesión (TTL 10 min) y puerto UDP para heartbeat
- `ERROR INVALID_CREDENTIALS` — usuario o contraseña incorrectos

### `MSG <token> <contenido_del_mensaje>`

Respuestas:
- `ACK` — mensaje aceptado, registrado en `historial.csv` y retransmitido a los demás clientes activos
- `ERROR SESSION_INVALID_TOKEN` — el token no existe (nunca existió, o ya fue revocado)
- `ERROR SESSION_EXPIRED` — el token existe pero la sesión ya no es válida (TTL vencido o timeout de heartbeat)

### Retransmisión (broadcast)

Ante un `MSG` válido, el servidor envía a todos los demás clientes con sesión activa:

```
INCOMING <usuario_emisor> <contenido_del_mensaje>
```

### `LOGOUT`

Cierra y revoca la sesión asociada al socket actual. Sin respuesta explícita — el servidor cierra la conexión.

### Comando no reconocido

Respuesta: `ERROR UNKNOWN_COMMAND`

## Protocolo UDP (puerto 9001)

### `HEARTBEAT <token>`

Sin respuesta explícita. Si el token existe, actualiza `timestamp_ultimo_heartbeat` en la sesión correspondiente. Debe enviarse cada 3 segundos mientras la sesión esté activa.

## Reglas de expiración y revocación

- Si no llega el primer heartbeat dentro de los 30 segundos posteriores al `LOGIN`, la sesión se revoca.
- Si transcurren más de 60 segundos sin un heartbeat válido, la sesión se revoca.
- Si transcurren más de 10 minutos (600s) desde la creación del token, la sesión se revoca sin importar los heartbeats.

En cualquiera de estos casos, el servidor cierra el socket TCP del cliente afectado y elimina/invalida su entrada en `sesiones.csv`.

## Librerías utilizadas (lado servidor, Python)

`socket`, `threading`, `http.server`, `csv`, `time`, `os`, `uuid`, `json`*

## Librerias utilizadas (lado cliente, Go)

bufio, net, net/http, fmt, os, time, strings
