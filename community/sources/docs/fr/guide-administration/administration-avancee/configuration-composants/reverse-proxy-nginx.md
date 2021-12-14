# Configuration avancée du reverse proxy HTTP Nginx de Canopsis

Le *reverse proxy* HTTP [Nginx](https://nginx.org) fournit l'accès à l'interface web de Canopsis. Il relaie aussi les API REST de Canopsis.

## Configuration par défaut

Le fichier principal de configuration de Nginx est `/etc/nginx/conf.d/default.conf`.

Ce fichier de configuration évolue lors des mises à jour de Canopsis.

Actuellement, cette configuration apporte :

* un accès unique à l'interface Canopsis par le port HTTP `80` standard ;
* un relais vers les API REST fournies par `canopsis-api` ;
* une mise en cache de certains éléments (fichiers CSS, fichiers JavaScript, images), afin d'améliorer le temps de chargement de l'interface ;
* une compression à la volée de la plupart des ressources, afin d'en accélérer le téléchargement dans les navigateurs ;
* des [entêtes de sécurité CORS](https://developer.mozilla.org/fr/docs/Web/HTTP/CORS), nécessaires pour certains applicatifs ;
* une prise en charge optionnelle d'HTTPS, HTTP/2 et des Websockets (voir ci-dessous).

!!! information
        La configuration Nginx par défaut déployée avec Canopsis est consultable sur le [dépôt Gitlab](https://git.canopsis.net/canopsis/canopsis-community/-/tree/develop/community/deploy-ansible/playbook/roles/canopsis/templates/nginx).

## Configuration additionnelle

### Changement du nom de serveur hôte HTTP (`server_name`)

À partir de Canopsis 4.4.0, Nginx est configuré pour utiliser le nom de serveur `localhost`, par défaut.

Suivez la procédure suivante, si le service HTTP doit être accessible avec un autre nom.

=== "Paquets CentOS 7"

    Éditez la variable `canopsis_server_name` du fichier `/etc/nginx/conf.d/default.conf`.

    Par exemple :
    ```nginx
    set $canopsis_server_name "canopsis.mon-si.fr";
    server_name $canopsis_server_name;
    ```

    Puis, rechargez le service `nginx` (`systemctl reload nginx`).

=== "Docker Compose"

    Si vous voulez éviter de surcharger l'intégralité du fichier `/etc/nginx/conf.d/default.conf`, vous pouvez modifier la variable d'environnement `CPS_SERVER_NAME` dans le fichier `compose.env` lié à votre Compose :

    ```ini
    CPS_SERVER_NAME=canopsis.mon-si.fr
    ```

    Puis, redémarrez le conteneur `nginx`.

### Activation d'HTTPS, HTTP/2 et les Websockets

À partir de Canopsis 4.4.0, une configuration activant HTTPS, HTTP/2 et les Websockets est disponible, mais n'est pas encore activée par défaut.

Consultez le [Guide d'activation d'HTTPS](reverse-proxy-nginx-https.md) pour en savoir plus.

### Modifications personnelles

Vous pouvez compléter le fichier de configuration fourni, ou créer un nouveau fichier `.conf` dans le répertoire `/etc/nginx/conf.d/`, comportant vos propres ajouts. Veillez à toujours vous synchroniser avec la dernière configugration officielle de Nginx après chaque mise à jour de Canopsis.

Notez cependant que seule la configuration Nginx proposée par défaut est prise en charge.
