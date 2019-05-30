# ejemplos-go-gorilla-mux

Hoy vamos a ver un **enrutador o router de Go**, que permite definir rutas y métodos HTTP para responder a ellos, de una manera fácil.

El enrutador, llamado **Mux** (que es de las herramientas de [Gorilla](https://www.gorillatoolkit.org/pkg/mux)) permite:

-   Definir **middleware** en las rutas, es decir, aplicar funciones que se ejecutan antes de cada petición HTTP y que permiten detener la ejecución o loguear determinadas cosas
-   Definición de rutas con **verbos HTTP**
-   Lectura de **parámetros** GET
-   Lectura de **variables dentro de la url**. Por ejemplo si definimos algo como usuario/{id} y se consulta a usuario/1 podemos obtener el valor 1 accediendo a la variable
-   Variables dentro de la URL con **expresiones regulares**

Mira la explicación en mi blog: https://parzibyte.me/blog/2019/05/30/enrutador-middleware-go-gorilla-mux/