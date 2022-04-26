Microservicios con Gin
======================

Ejemplo de CRUD utilizando Mongo como base de datos y Redis como cache.

Para poder ejecutar este proyecto es necesario tener una instancia de mongo y  
de redis localmente en el puerto default, y se deben establecer las siguientes
variables de entorno:

- MONGO_URI="mongodb://username:password@localhost:27017"
- MONGO_DATABASE=demo