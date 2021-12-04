# transmission

Este microservicio se encarga de realizar la transmision de audio y video a los clientes conectados a una reunion, usa una api desarrollada en Go para esto 
y una base de datos en Cassandra para almacenar los clientes conectados a una reunion.
---

## API
La api en Go expone 5 endpoints

| Metodo | URL | Descripcion |
| --- | --- | --- |
| GET | /api/v1/session/{idMeeting} | Retorna todos los clientes conectados a una reunion |
| POST | /api/v1/session | Añade un nuevo registro |
| DELETE | /api/v1/session/{idMeeting}/{idSession} | Elimina un registro correspondiente a una reunion y una sesion especifica |
| POST | /api/v1/image | Registra una nueva imagen para ser transmitida a los demas clientes de la reunion |
| POST | /api/v1/sound | Registra una nueva seccion de audio para ser transmitida a los demas clientes de la reunion |

---

## Cassandra
La base de datos en cassandra tiene una unica tabla donde se almacenan los clientes conectados y la ultima imagen transmitida

| Columna | Tipo | Descripcion |
| --- | --- | --- |
| idmeeting (PK) | text | Codigo de la reunion |
| idsession (PK) | int | Codigo del cliente conectado |
| media | byte[] | Ultima imagen enviada al servidor por un cliente |
