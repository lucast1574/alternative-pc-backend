# Backend Alternative

Backend modular desarrollado en Go con MongoDB.

## Arquitectura
- `cmd/api`: Punto de entrada de la aplicación.
- `internal/modules`: Contiene la lógica de negocio dividida por módulos funcionales.
- `internal/platform`: Infraestructura y servicios compartidos.

## Cómo ejecutar con Docker
```bash
docker build -t backend-alternative .
docker run -p 8080:8080 backend-alternative
```

## Ejecución local
```bash
go run cmd/api/main.go
```
