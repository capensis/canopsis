# Reverse Proxy

**Remarque** : Ces informations sont données à titre indicatif.  
La commande `canoctl deploy` propose déjà la mise à dispoition d'un reverse proxy.

## Reverse proxy avec Apache2

Pour utiliser un autre port pour votre serveur web, vous pouvez utiliser [Apache "mod_proxy"](https://httpd.apache.org/docs/2.4/fr/mod/mod_proxy.html)

### Installation

Installer Apache2 si cela n'est pas déjà fait.

#### Ubuntu / Debian
```
apt-get install apache2
```

#### CentOS / RHEL
```
yum install httpd
```

### Configurer le VHost (replace ``<DNS_NAME>``)

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

```
service apache2 restart or service httpd restart
update-rc.d apache2 defaults or chkconfig httpd on 
```

## Reverse Proxy avec Ngnix

### Ubuntu / Debian
```
apt-get install nginx
``` 

### CentOS / REHL 
```
yum install nginx
```

**Don't forget to add epel repositories for RedHat/CentOS**

### Configurer le VHost (replace ``<DNS_NAME>``)

```
server {
	listen 80;
	server_name <DNS_NAME>;

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
