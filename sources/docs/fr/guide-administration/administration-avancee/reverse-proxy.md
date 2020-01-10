# Configuration avancée du reverse proxy HTTP de Canopsis

Par défaut, l'interface Canopsis et ses API sont accessibles sur le port 8082, par le biais du [serveur HTTP Gunicorn](https://gunicorn.org).

À partir de [Canopsis 3.35.0](../../notes-de-version/3.35.0.md), un *reverse proxy* HTTP ([Nginx](https://nginx.org)) est aussi proposé par défaut, afin d'apporter une configuration plus avancée du service HTTP lié à Canopsis.

## Activation de Nginx dans une installation Canopsis

### Mise en place dans une installation Docker

En environnement Docker, une image [`canopsis/nginx`](https://hub.docker.com/repository/docker/canopsis/nginx) est publiquement disponible depuis Canopsis 3.35.0.

Vous devez, pour cela, posséder une section de ce type dans votre fichier Docker Compose :

```yaml
nginx:
  image: canopsis/nginx:${CANOPSIS_IMAGE_TAG}
  ports:
    - "80:80"
  env_file:
    - compose.env
  environment:
    - TARGET=http://webserver:8082
  depends_on:
    - "webserver"
  restart: unless-stopped
```

La variable `TARGET` correspond à l'URL d'écoute du service Gunicorn, qui s'appelle `webserver` et qui utilise le port `8082` par défaut.

### Mise en place dans une installation par paquets

Dans une installation par paquets, la commande [`canoctl`](../installation/installation-paquets.md) s'occupe d'installer un serveur Nginx, pour toute nouvelle installation à partir de Canopsis 3.35.0.

## Configuration de Nginx

### Configuration par défaut

La configuration de Nginx proposée par défaut avec Canopsis se situe dans le fichier `/etc/nginx/conf.d/default.conf` (contenu dans l'image `canopsis/nginx`, en environnement Docker).

Ce fichier de configuration est amené à évoluer avec les mises à jour de Canopsis.

Actuellement, cette configuration apporte :

*  un accès unique à l'interface Canopsis par le port HTTP `80` standard ; 
*  une mise en cache de certains éléments (fichiers CSS, fichiers JavaScript, images), afin d'améliorer le temps de chargement de l'interface ;
*  une compression à la volée de la plupart des ressources, afin d'en accélérer le téléchargement dans les navigateurs ;
*  des [entêtes de sécurité CORS](https://developer.mozilla.org/fr/docs/Web/HTTP/CORS), nécessaires pour certains applicatifs.

### Configuration additionnelle

Vous pouvez compléter le fichier de configuration fourni, ou créer un nouveau fichier `.conf` dans le répertoire `/etc/nginx/conf.d/`, comportant vos propres ajouts.

Vous pouvez notamment vous en servir afin :

*  [d'ajouter des certificats HTTPS](https://nginx.org/en/docs/http/configuring_https_servers.html) ;
*  [d'activer le protocole HTTP/2](https://nginx.org/en/docs/http/ngx_http_v2_module.html) ;
*  [d'ajouter des entêtes](https://nginx.org/en/docs/http/ngx_http_headers_module.html) dont vous pourriez avoir besoin ;
*  etc.

Canopsis ne supporte néanmoins que la configuration Nginx proposée par défaut.

## Utilisation d'un autre serveur HTTP que Nginx

Vous pouvez aussi remplacer Nginx par tout autre serveur HTTP capable de faire du *reverse proxy*, tel qu'[Apache HTTPD](https://httpd.apache.org) ou [HAProxy](https://www.haproxy.org), en fonction de vos besoins.

Vous devez, pour cela, désactiver le service Nginx fourni avec Canopsis (`systemctl disable nginx` ou arrêt de l'image Docker `canopsis/nginx`).

Vous pouvez ensuite transposer la configuration `/etc/nginx/conf.d/default.conf` existante vers la syntaxe de votre propre serveur HTTP. L'utilisation d'un autre serveur HTTP que Nginx est possible, mais n'est cependant pas officiellement supportée.
