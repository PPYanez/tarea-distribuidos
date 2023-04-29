# Tarea 1 de Sistemas Distribuidos

Por:
* José Quezada
* Pedro Yáñez

## ¿Cómo ejecutar la tarea?
Primero se debe levantar el servidor de gin usando go run server/server.go

Luego, se debe ejecutar el menú para interactuar con el servidor, utilizando go run menu/menu.go

MongoDB está instalada en un container de docker, puede que al probar la tarea este no esté ejecutandose, para ejecutar el container:

docker run 6d707605d714

En caso que esto no funcione, significa que por alguna razón se borró la imágen de mongodb, para volver a crearla:

docker run --name mongo -d -p 27017:27017 mongodb/mongodb-community-server:latest

## Consideraciones

No alcanzamos a hacer la funcionalidad de actualizar una reserva completamente, solo está funcionando la opción de modificar fecha de vuelo.

La búsqueda solo se puede hacer con el apellido igual al con el que está creada la reserva, o sea que si una reserva tiene un pasajero con apellido Rojas, buscarla con el apellido rojas o ROJAS no funciona.