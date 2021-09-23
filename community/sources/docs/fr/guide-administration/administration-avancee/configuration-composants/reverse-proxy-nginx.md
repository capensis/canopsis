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
        La configuration Nginx par défaut déployée avec Canopsis est consultable sur le [dépôt Gitlab](https://git.canopsis.net/canopsis/canopsis-community/-/tree/develop/community/deploy-ansible/playbook/roles/canopsis/templates/nginx).

## Configuration additionnelle

### Activation d'HTTPS, HTTP/2 et les Websockets

À partir de Canopsis 4.4.0, une configuration activant HTTPS, HTTP/2 et les Websockets est disponible, mais désactivée par défaut.

Consultez le [Guide d'activation d'HTTPS](../../webserver/https.md) pour en savoir plus.

### Modifications personnelles

Vous pouvez compléter le fichier de configuration fourni, ou créer un nouveau fichier `.conf` dans le répertoire `/etc/nginx/conf.d/`, comportant vos propres ajouts. Veillez à toujours vous synchroniser avec la dernière configugration officielle de Nginx après chaque mise à jour de Canopsis.

Notez cependant que Canopsis ne supporte que la configuration Nginx proposée par défaut.
