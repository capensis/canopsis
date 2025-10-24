# Activation de HTTPS dans Canopsis

À partir de Canopsis 4.4.0, une configuration HTTPS est proposée avec [Nginx](reverse-proxy-nginx.md), mais elle n'est cependant pas encore activée par défaut. Ce guide décrit sa configuration et son activation.

## Apports de la configuration HTTPS

La configuration HTTPS proposée dans Nginx vous permet :

* de sécuriser vos échanges HTTP, ce qui est notamment recommandé si votre accès à Canopsis ne se fait pas au travers d'un intranet sécurisé ou d'un VPN ;
* d'activer implicitement la prise en charge d'HTTP/2, pour de meilleures performances web dans certaines conditions ;
* de bénéficier des [Websockets](https://developer.mozilla.org/fr/docs/Web/API/WebSockets_API) et donc de la nouvelle fonctionnalité Healthcheck apparue avec Canopsis 4.4.0.

!!! note
    Les prérequis pour ces fonctionnalités sont les suivants :

    * assigner un [FQDN](https://fr.wikipedia.org/wiki/Fully_qualified_domain_name) à votre service web Canopsis (ex : `canopsis.mon-si.fr`) ;
    * avoir un navigateur [officiellement pris en charge](../../../guide-utilisation/limitations/index.md#compatibilite-des-anciens-navigateurs) et raisonnablement récent ;
    * disposer d'au moins HTTP/1.1 et des Websockets dans ce navigateur ;
    * disposer de TLSv1.2 ou de TLSv1.3 sur vos clients HTTPS (note : TLSv1.3 n'est pas disponible dans les paquets CentOS 7) ;
    * disposer d'OpenSSL sur votre serveur, avec ses dernières mises à jour de sécurité ;
    * utiliser une autorité de certification SSL/TLS en raccord avec les pratiques internes de votre SI (voir ci-dessous) ;
    * mettre en place une politique de renouvellement des certificats et une surveillance de leur expiration.

## Choix du type de certificat HTTPS

De façon générale, l'écosystème HTTPS nécessite la mise en place de certificats devant être reconnus et acceptés par les clients.

Pour cela, trois options s'offrent à vous :

1. mettre en place un certificat signé par une autorité de certification tierce (Gandi, Digicert, Let's Encrypt…) ;
2. mettre en place un certificat signé par l'autorité de certification interne à votre SI (le cas échéant) ;
3. mettre en place un certificat autosigné et déployer des règles d'exception pour ce domaine dans votre SI (**non recommandé**).

Les deux premières options vous imposent de vous rapprocher d'une autorité de certification compétente et de venir compléter la configuration Nginx à l'aide des explications données ci-dessous.

La troisième option ne nécessite pas d'autorité de certification, mais de façon générale cette méthode n'est pas recommandée et elle vous impose d'ajouter une exception de sécurité sur chaque client HTTPS (à votre charge).

## Activation de la configuration HTTPS

La configuration Nginx proposée par défaut s'attend à ce que votre certificat et sa clé privée soient présents aux emplacements suivants :

| Type de fichier | Emplacement |
| --------------- | ----------- |
| Certificat SSL/TLS, au format PEM | `/etc/nginx/ssl/cert.crt` |
| Clé privée, au format PEM | `/etc/nginx/ssl/key.key` |

!!! attention
    Par mesure de sécurité, veillez à ce que le répertoire `/etc/nginx/ssl` soit bien attribué à `root:root` et qu'il dispose bien de permissions restreintes `0700`.

### Ajout d'un certificat sécurisé (recommandé)

En premier lieu, rapprochez-vous de votre autorité de certification afin de connaître la procédure à suivre pour la génération de votre clé privée et pour l'obtention d'un certificat signé.

Vous devez ensuite placer ces fichiers au bon endroit sur votre serveur Canopsis, en fonction de la [méthode d'installation](../../installation/index.md#methodes-dinstallation-de-canopsis) que vous avez choisie.

=== "Paquets CentOS 7"

    Assurez-vous tout d'abord de la bonne restriction des accès à `/etc/nginx/ssl` avec la commande suivante :

    ```sh
    install -d -m 0700 -o root -g root /etc/nginx/ssl
    ```

    Placez ensuite votre certificat dans le fichier `/etc/nginx/ssl/cert.crt` et votre clé privée dans le fichier `/etc/nginx/ssl/key.key`.

=== "Docker Compose"

    L'injection des fichiers attendus se fait à l'aide de volumes.

    Exemple avec `cert.crt` comme certificat et `key.key` comme clé privée :

    ```yaml hl_lines="12-14"
      nginx:
        image: ${DOCKER_REPOSITORY}${COMMUNITY_BASE_PATH}nginx:${CANOPSIS_IMAGE_TAG}
        ports:
          - "80:80"
          - "443:443"
        env_file:
          - compose.env
        depends_on:
          - "api"
        restart: unless-stopped
        volumes:
          #- nginxcerts:/etc/nginx/ssl
          - ./cert.crt:/etc/nginx/ssl/cert.crt:ro
          - ./key.key:/etc/nginx/ssl/key.key:ro
    ```

### Ajout d'un certificat autosigné (non recommandé)

La procédure suivante peut être utilisée afin de mettre en place un certificat autosigné sur votre environnement.

Notez au préalable que :

* Vous devez veiller à renouveler les certificats autosignés vous-mêmes chaque année (avec une tolérance jusqu'à 2 ans).
* L'utilisation de certificats autosignés provoquera l'affichage d'un message dans votre navigateur lors de la connexion à Canopsis. Vous devrez ajouter une exception sur chaque navigateur devant accéder à Canopsis.
* De façon générale, les certificats autosignés n'assurent pas un niveau de sécurité suffisant dans un SI, et ne sont donc **pas recommandés**.

=== "Paquets CentOS 7"

    Sur l'environnement cible, exécutez les commandes suivantes en remplaçant `canopsis.mon-si.fr` par le vrai FQDN de votre service Canopsis :

    ```sh
    install -d -m 0700 -o root -g root /etc/nginx/ssl
    ( umask 077 ; openssl req -x509 -nodes -days 730 -newkey rsa:2048 -sha256 \
        -keyout /etc/nginx/ssl/key.key \
        -out /etc/nginx/ssl/cert.crt \
        -subj "/CN=canopsis.mon-si.fr" )
    ```

=== "Docker Compose"

    Dans le cas d'un certificat autosigné vous n'avez rien à faire : le conteneur `nginx` va en générer un de façon automatique lors de son démarrage, sauf si vous surchargez les fichiers attendus dans `/etc/nginx/ssl` à l'aide d'un volume.

## Activation de la configuration HTTPS

=== "Paquets CentOS 7"

    Éditez le fichier `/etc/nginx/conf.d/default.conf` afin de configurer votre FQDN (ex : `canopsis.mon-si.fr`), et décommentez la ligne `#include /etc/nginx/https.inc` afin d'activer la configuration HTTPS.

    ```nginx hl_lines="1 5"
    set $canopsis_server_name "canopsis.mon-si.fr";
    server_name $canopsis_server_name;

    # Uncomment the next line to enable HTTPS
    include /etc/nginx/https.inc;
    ```

    Puis, redémarrez le service Nginx (`systemctl restart nginx`).

=== "Docker Compose"

    Pour activer HTTPS dans le conteneur vous devrez modifier ses variables d'environnement dans le fichier `compose.env` lié à votre fichier de configuration Compose.

    | Variable | Description |
    | -------- | ----------- |
    | `CPS_SERVER_NAME` | FQDN sur lequel Canopsis sera disponible (ex : `canopsis.mon-si.fr`) |
    | `CPS_ENABLE_HTTPS` | État d'activation d'HTTPS. Si cette variable vaut `true` HTTPS sera activé, dans tous les autres cas il sera désactivé |

    Exemple de `compose.env` pour activer HTTPS avec `canopsis.mon-si.fr` en FQDN :

    ```ini
    CPS_SERVER_NAME=canopsis.mon-si.fr
    CPS_ENABLE_HTTPS=true
    ```

    Puis, redémarrez le conteneur `nginx`.

## Sécurisation supplémentaire des cookies (optionnel)

Si vous pouvez garantir que la totalité de vos clients utilisera uniquement une connexion HTTPS (et non pas HTTP), vous pouvez améliorer la sécurité des cookies en leur ajoutant [l'attribut `Secure`](https://developer.mozilla.org/fr/docs/Web/HTTP/Cookies#cookies_secure_et_httponly), garantissant qu'ils ne pourront plus être manipulés sans HTTPS.

Pour cela, le moteur `canopsis-api` doit être lancé avec l'option `-secure`.

=== "Paquets CentOS 7"

    Exécutez les commandes suivantes pour forcer le moteur `canopsis-api` à être lancé avec l'option `-secure` :

    ```sh
    mkdir -p /etc/systemd/system/canopsis-service@canopsis-api.service.d
    cat > /etc/systemd/system/canopsis-service@canopsis-api.service.d/override.conf << EOF
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/%i -secure
    EOF

    systemctl daemon-reload
    systemctl restart canopsis-service@canopsis-api
    ```

=== "Docker Compose"

    Ajoutez l'option `-secure` aux arguments de `canopsis-api` :

    ```yaml hl_lines="8"
      api:
        ...
        ports:
          - "8082:8082"
        env_file:
          - compose.env
        restart: unless-stopped
        command: /canopsis-api -secure
    ```

    Puis, relancez ce conteneur.

## Utilisation d'un autre applicatif ou équipement pour servir les flux HTTPS

Dans la politique de sécurité de certains SI, la terminaison SSL/TLS doit forcément passer par un applicatif ou par un équipement dédié.

Si vous êtes dans cette situation, vous pouvez réaliser votre propre terminaison SSL/TLS avec vos propres moyens et selon vos propres conditions, mais ce type de modification n'est pas officiellement pris en charge par Canopsis. Tout éventuel ajustement à faire et toute maintenance sont donc à votre charge.
