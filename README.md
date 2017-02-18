# DataTweet

Provee datos y estadisticas sobre tus tweets mas recientes.

Luego de seguir el proceso de autenticacion, la API retorna los datos(token, secret) que las aplicaciones deberan guardar para hacer consultas posteriores y poder utilizar Datatweet.

####Con DataTweet podes obtener informacion sobre el top ten de porcentaje de likes que realiza un usuario en sus ultimos 200 tweets, haciendo un GET a `/favorites` y pasando como parametros token, secret y el scren_name

####Ademas, podes obtener un diccionario con el top ten de porcentaje de las palabras mas usadas de los ultimos 200 tweets en orden descendente, haciendo un GET a `/dictionary` y pasando como parametros token, secret y el scren_name

Uso de la API:

-Obtener access token y consumer Secret:

1) http://localhost:8888/permission 
2) se redirige a la pagina de twitter, donde debera loguearse y aceptar los permisos para poder continuar.
3) Luego se redirige automaticamente a /maketoken donde se retorna un json con el access token, consumer secret, que se usaran para las consultas.

-para obtener los favoritos: 

`http://localhost:8888/favorites?screen_name=(nombre de usuario de twitter)&token=(access token obtenido)&secret=(consumer secret obtenido)`

-para obtener tu diccionario: 

`http://localhost:8888/dictionary?screen_name=(nombre de usuario de twitter)&token=(access token obtenido)&secret=(consumer secret obtenido)`

Ejemplos de uso: 

Es posible averiguar sobre la cuenta de otro usuario distinto al nuestro

1) `http://localhost:8888/dictionary?screen_name=BarackObama&token=811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16&secret=p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv` 
[link ejemplo 1](http://localhost:8888/dictionary?screen_name=BarackObama&token=811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16&secret=p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv)

2) `http://localhost:8888/favorites?screen_name=realDonaldTrump&token=811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16&secret=p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv` 
[link ejemplo 2](http://localhost:8888/favorites?screen_name=realDonaldTrump&token=811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16&secret=p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv)
