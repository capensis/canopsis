# Reverse Proxy

**Remarque** : Ces informations sont données à titre indicatif.  
La commande `canoctl deploy` propose déjà la mise à dispoition d'un reverse proxy.

## Reverse proxy avec Apache2

Pour utiliser un autre port pour votre serveur web, vous pouvez utiliser [Apache "mod_proxy"](https://httpd.apache.org/docs/2.4/fr/mod/mod_proxy.html)

### Installation

Installer Apache2 si cela n'est pas déjà fait.

#### Ubuntu / Debian

```sh
apt install apache2
```

#### CentOS / RHEL

```sh
yum install httpd
```

### Configurer le VHost

```
<VirtualHost *:80>
    ProxyRequests Off
    ProxyPreserveHost On

    <Location />
        ProxyPass http://127.0.0.1:8082/
        ProxyPassReverse http://127.0.0.1:8082/
    
        SetOutputFilter DEFLATE
        SetEnvIfNoCase Request_URI \.(?:gif|jpe?g|png)$ no-gzip dont-vary
        SetEnvIfNoCase Request_URI \.(?:exe|t?gz|zip|gz2|sit|rar)$ no-gzip dont-vary
    </Location>
</VirtualHost>
```

### Restart Apache & démarrage du service automatique

```sh
service apache2 restart or service httpd restart
update-rc.d apache2 defaults or chkconfig httpd on 
```

## Reverse Proxy avec Ngnix

### Ubuntu / Debian

```sh
apt install nginx
```

### CentOS / REHL 

```sh
yum install nginx
```

**Don't forget to add epel repositories for RedHat/CentOS**

### Configurer le VHost

```
server {
	listen 80;
	server_name <DNS_NAME>; # À adapter

	gzip on;
	gzip_disable "msie6";
	
	gzip_comp_level 6;
	gzip_min_length 1100;
	gzip_buffers 16 8k;
	gzip_proxied any;
	gzip_types
	    text/plain
	    text/css
	    text/js
	    text/xml
	    text/javascript
	    application/javascript
	    application/x-javascript
	    application/json
	    application/xml
	    application/xml+rss;

	location / {
		proxy_set_header X-Forwarded-Host $host;
		proxy_set_header X-Forwarded-Server $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_pass http://127.0.0.1:8082;    
	}	
}
```

### Restart nginx & démarrage du service automatique

```
service nginx restart or service nginx restart
update-rc.d nginx defaults or chkconfig nginx on 
```

## Reverse Proxy avec HA Proxy

### Ubuntu / Debian

```
apt install haproxy
```

### Centos / RHEL

```
yum install haproxy
```

### Configurer un accès au Webserver de Canopsis sur Port `80`

Éditer le fichier `/etc/haproxy/haproxy.cfg` et créer :

* Un bloc `frontend` qui va écouter sur un port d'écoute ( ici `TCP/80` ) et rediriger vers un backend `canopsis-webserver`

```
frontend canopsis-webui
  bind *:80
  option forwardfor
  mode http
  compression algo gzip
  compression type text/html text/plain text/css text/js text/xml text/javascript application/javascript application/x-javascript application/json application/xml application/xml+rss
  default_backend canopsis-webserver
```

* Un bloc `backend` qui va définir la redirection vers le webserver

```
backend canopsis-webserver
  server canopsis 127.0.0.1:8082
```

### Gestion des CORS Policy pour le serveur Web Gunicorn

Pour accéder à des ressources fournies par le Webserver Gunicorn de Canopsis depuis des origines/domaines multiples, il faut mettre en place une configuration particulière avec HA Proxy

Il se trouve que lorsque les cors policy sont contrôlées, un appel à l'URL avec la méthode HTTP `OPTIONS` est effectué vers le serveur web par le navigateur ou client HTTP.

Le serveur web gunicorn de Canopsis ne supportant pas cette méthode, voici une configuration `haproxy` capable de répondre à la place du serveur web

````
frontend canopsis-webserver
  ...  
  http-response set-header Access-Control-Allow-Origin "http://domaine.tld"
  http-response set-header Access-Control-Allow-Credentials true
  ...
  acl is_options method OPTIONS
  use_backend cors_backend if is_options
  ...

backend cors_backend
  errorfile 503 /etc/haproxy/errors/cors_canopsis_200.http
````

Avec comme contenu du fichier `/etc/haproxy/errors/cors_canopsis_200.http` 

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: http://votredomaininterne
Access-Control-Max-Age: 31536000
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Authorization
Content-Length: 0
Cache-Control: private



```

Un backend `cors_backend` pourra ainsi répondre avec les bons entêtes attendus pour les CORS Policy lorsque la méthode demandée est `OPTIONS`.

Une fois l'appel à la méthode `OPTIONS` passée, les en-têtes HTTP seront ajoutés à la volée pour la réponse aux appels `GET` suivants.