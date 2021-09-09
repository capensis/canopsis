# Configuration avancée du reverse proxy HTTP Nginx de Canopsis

Le *reverse proxy* HTTP ([Nginx](https://nginx.org)) fournit l'accès à l'interface web de Canopsis.

## Configuration par défaut

La configuration de Nginx proposée par défaut avec Canopsis se situe dans le fichier `/etc/nginx/conf.d/default.conf` (contenu dans l'image `canopsis/nginx`, en environnement Docker).

Ce fichier de configuration est amené à évoluer avec les mises à jour de Canopsis.

Actuellement, cette configuration apporte :

*  un accès unique à l'interface Canopsis par le port HTTP `80` standard ; 
*  une mise en cache de certains éléments (fichiers CSS, fichiers JavaScript, images), afin d'améliorer le temps de chargement de l'interface ;
*  une compression à la volée de la plupart des ressources, afin d'en accélérer le téléchargement dans les navigateurs ;
*  des [entêtes de sécurité CORS](https://developer.mozilla.org/fr/docs/Web/HTTP/CORS), nécessaires pour certains applicatifs.

!!! information
        La configuration Nginx par défaut déploiée avec Canopsis est consultable sur le [dépôt Gitlab](https://git.canopsis.net/canopsis/canopsis-community/-/tree/develop/community/deploy-ansible/playbook/roles/canopsis/templates/nginx).

## Configuration additionnelle

Vous pouvez compléter le fichier de configuration fourni, ou créer un nouveau fichier `.conf` dans le répertoire `/etc/nginx/conf.d/`, comportant vos propres ajouts.

Vous pouvez notamment vous en servir afin :

*  [d'ajouter des certificats HTTPS](https://nginx.org/en/docs/http/configuring_https_servers.html) ;
*  [d'activer le protocole HTTP/2](https://nginx.org/en/docs/http/ngx_http_v2_module.html) ;
*  [d'ajouter des entêtes](https://nginx.org/en/docs/http/ngx_http_headers_module.html) dont vous pourriez avoir besoin ;
*  etc.

Canopsis ne supporte néanmoins que la configuration Nginx proposée par défaut.
