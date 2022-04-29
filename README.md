Microservicios con Gin
======================

Ejemplo de CRUD utilizando Mongo como base de datos y Redis como cache.

Para poder ejecutar este proyecto es necesario tener una instancia de MongoDB y  
de Redis localmente en el puerto default.
Esto se puede lograr con Docker facilmente ejecutando:

* MongoDB

> docker run -d --name mongodb -v /path/to/volume -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=somepassword -p 27017:27017 mongo

* Redis
  
> docker run -d --name redis -p 6379:6379 redis 

También se deben establecer las siguientes variables de entorno:

- MONGO_URI="mongodb://username:password@localhost:27017"
- MONGO_DATABASE=demo
- JWT_SECRET=somesecret

Para realizar pruebas de rendimiento del cache se utilizó la herramienta 
`apache-benchmark` y se obtuvo la siguiente gráfica comparativa:

![Comparación Cache](benchmark/benchmark.png)