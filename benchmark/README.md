Benchmark Cache
===============

Para estas pruebas se utilizó la herramienta [ab](https://httpd.apache.org/docs/2.4/programs/ab.html)
con el siguiente comando:

> ab -n 2000 -c 100 -g output.data http://localhost:8080/recipes

El cual lanza 2000 peticiones al endpoint, en bloques de 100 de manera 
concurrente y guarda los resultados en el archivo output.data.

Los resultados se pueden graficar utilizando gnuplot y se obtiene una 
comparación como la siguiente:

![Comparación](benchmark.png)