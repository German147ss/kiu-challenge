
# Proyecto Kiu Challenge

## Descripción

Este proyecto en Go maneja una API simple de saludo. La API mantiene un registro de los nombres que ha saludado anteriormente y ofrece una respuesta personalizada si el nombre se encuentra en su registro.

## Estructura del Proyecto

El proyecto consta de dos archivos principales:

- `main.go`: Este archivo maneja todas las solicitudes a la API.
- `main_test.go`: Este archivo contiene las pruebas para el microservicio.

## Instalación y Uso

Para instalar y usar este proyecto, sigue estos pasos:

1. Clona el repositorio.
   ```
   git clone https://github.com/german147ss/kiu-challenge.git
   ```

2. Navega al directorio del proyecto.
   ```
   cd kiu-challenge
   ```

3. Ejecuta las pruebas.
   ```
   go test
   ```

4. Inicia el servidor.
   ```
   go run main.go
   ```

La API ahora estará disponible en `localhost:8080`.

## Endpoints de la API

- `POST /hello`: Saluda al nombre proporcionado en el cuerpo de la solicitud. Si el nombre ya ha sido saludado antes, ofrece una respuesta de "hola devuelta".
- `GET /hello`: Devuelve una lista de todos los nombres que han sido saludados.

## Contribuciones

Las contribuciones a este proyecto son bienvenidas. Para contribuir, por favor haz un fork del repositorio, crea una nueva rama para tus cambios y luego abre un Pull Request una vez que estés listo para fusionar tus cambios en la rama principal.

## Licencia

Este proyecto está licenciado bajo los términos de la Licencia MIT.