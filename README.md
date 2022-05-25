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

Para la autenticación se utiliza el servicio en su capa gratuita de [Auth0](https://auth0.com/es)
Se deben establecer las siguientes variables de entorno:

- MONGO_URI="mongodb://username:password@localhost:27017"
- MONGO_DATABASE="demo"
- JWT_SECRET="somesecret"
- REDIS_URI="localhost:6379"
- REDIS_SECRET="somesecret"
- AUTH0_DOMAIN="YOURDOMAIN.auth0.com"
- AUTH0_API_IDENTIFIER="https://api.recipes.io"

Para realizar pruebas de rendimiento del cache se utilizó la herramienta 
`apache-benchmark` y se obtuvo la siguiente gráfica comparativa:

![Comparación Cache](benchmark/benchmark.png)
