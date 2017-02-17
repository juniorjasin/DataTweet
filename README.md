# DataTweet

Provee estadisticas sobre tus tweets y mejora tu experiencia en twitter.

Luego de seguir el proceso de autenticacion, la api te retorna los datos(token, secret) que las aplicaciones deberan guardar para hacer consultas posteriores y poder utilizar Datatweet.

Con DataTweet podes obtener informacion sobre el top ten de porcentaje de likes que realiza un usuario en sus ultimos 200 tweets, haciendo un GET a /favorites y pasando como parametros token, secret y el scren_name

Ademas, podes obtener un diccionario con el top ten de porcentaje de las palabras mas usadas de los ultimos 200 tweets en orden descendente, haciendo un GET a /dictionary y pasando como parametros token, secret y el scren_name

Ejemplo:

-para obtener los favoritos: 

http://localhost:8888/favorites?screen_name=juniorjasin&token=811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16&secret=p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv

-para obtener tu diccionario: 

http://localhost:8888/dictionary?screen_name=juniorjasin&token=811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16&secret=p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv
