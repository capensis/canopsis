# Installation de l'UIv3 Canopsis (standalone)

Ce guide vous permet de déployer l'interface graphique de Canopsis de manière indépendante du backend Canopsis.  

!!! warning
    La compatibilité de l'UIv3 et des APIs Canopsis est assurée dans chaque version.
    Le fait de déployer une UIv3 et des APIs Canopsis de version différente vous explose à des incompatibilités, des incohérences de comportements et, de manière générale, à des effets de bord

!!! note 
    Cette méthode de déploiement est valable pour des versions de Canopsis >= 3.15.0

## Récupération des sources et build de l'UIv3

Vous devez déterminer la version que vous souhaitez déployer, 3.15.0 dans notre exemple.

````
git clone https://git.canopsis.net/canopsis/canopsis.git -b 3.15.0
cd canopsis/sources/webcore/src/canopsis-next/
````

Vous devez réfléchir à la destination de l'application (au sens RootDirectory d'un serveur web).  
Dans cet exemple, il s'agit de */var/www/html/canopsis-uiv3*

````
yarn install && yarn build --mode standalone --dest /var/www/html/canopsis-uiv3
````

## Configuration du serveur web

Il est nécessaire de service l'application fraichement déployée et d'y appliquer une configuration de *reverse proxy* vers les APIs Canopsis.  
Dans notre cas, il s'agit de *http://canopsis-api:8082/*
Voici un exemple de configuration pour *Nginx*.  

````
server {
	listen 80 default_server;
	listen [::]:80 default_server;
	server_name localhost;

	root /var/www/html/canopsis-uiv3;

	index index.html index.htm index.nginx-debian.html index.php;

        location / {
          try_files $uri $uri/ @rewrites;
        }

        location @rewrites {
          rewrite ^(.+)$ /index.html last;
        }

	location /api {
                proxy_pass         http://canopsis-api:8082/;
        }
}
````

## Tests

A ce stade, vous pouvez utiliser votre navigateur et vous rendre sur l'url de votre serveur web.  
Vous aurez alors à disposition l'UIv3 de Canopsis.
