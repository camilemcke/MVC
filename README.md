# Beauty Source Inteligencia de inventario

Un sistema distribuido para consultar productos, inventario y métricas de Beauty Source, la empresa en la que trabajo actualmente. Ya se que no es un tema de ciberseguridad pero prefiero hacerlo con algo que si me interese aplicando lo que vimos

## Ejecutar el proyecto

```bash
docker compose up --build
```

- Frontend: http://localhost:8081
- Middleware: http://localhost:8082
- PostgreSQL: localhost:5432

## Arquitectura

Cada servicio en su propio contenedor de Docker

- *Frontend:* ver métricas y consultar productos
- *Middleware:* tomas las peticiones y las redirige al servicio que le toca
- *Products:* información de los productos
- *Inventory:* existencias de cada producto
- *Analytics:* indicadores de inventario, acá esta lo realmente útil para la empresa, ahorita es muy sencillo
- *Load Balancer* reparte las peticiones de productos entre dos instancias (2 backends), con esta aplicación tan sencilla no aplicaría pero en producción si
- *PostgreSQL:* base de datos de productos e inventario


```text
Frontend
   |
   v
Middleware
   |
   +---- /product ----> Load Balancer ----> backend-1
   |                         |
   |                         +------------> backend-2
   |
   +---- /inventory --> Inventory
   |
   +---- /analytics --> Analytics
                              |
                              v
                          PostgreSQL
```

## Service Discovery

No quise meter algo complicado sino que usa el DNS interno de Docker Compose como para Service Discovery.

Cada servicio se comunica usando sus dentro de la red de Docker. Así:

```text
http://inventory:8080
http://analytics:8080
http://load-balancer:8080
```

No usa direcciones IP fijas sino que Docker resuelve automáticamente el nombre de cada servicio a la dirección que le corresponde en la red

## Load Balancing

Load balancing aplicado a consultas de productos con dos instancias de backend en contenedores separados. Alterna entre 1 y 2

```text
Petición 1 -> backend-1
Petición 2 -> backend-2
Petición 3 -> backend-1
Petición 4 -> backend-2
```

