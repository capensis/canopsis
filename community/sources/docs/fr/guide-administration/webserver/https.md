# Activation de HTTPS dans Canopsis

À partir de Canopsis 4.4.0, une configuration HTTPS est proposée avec [Nginx](../administration-avancee/configuration-services/reverse-proxy-nginx.md), mais elle n'est cependant pas encore activée par défaut. Ce guide décrit son activation et son utilisation.

## Apports de la configuration HTTPS

La configuration HTTPS proposée dans Nginx vous permet :

* de sécuriser vos échanges HTTP, ce qui est notamment recommandé si votre accès à Canopsis ne se fait pas au travers d'un intranet sécurisé ou d'un VPN ;
* d'activer implicitement la prise en charge d'HTTP/2, pour de meilleures performances web dans certaines conditions ;
* de bénéficier des [Websockets](https://developer.mozilla.org/fr/docs/Web/API/WebSockets_API) et donc de la nouvelle fonctionnalité Healthcheck apparue avec Canopsis 4.4.0.

!!! note
    Les prérequis pour ces fonctionnalités sont les suivants :

    * assigner un [FQDN](https://fr.wikipedia.org/wiki/Fully_qualified_domain_name) à votre service web Canopsis (ex : `canopsis.mon-si.fr`) ;
    * avoir un navigateur [officiellement pris en charge](../../guide-utilisation/limitations/index.md#compatibilite-des-anciens-navigateurs) et raisonnablement récent ;
    * disposer d'HTTP/1.1 et des Websockets dans ce navigateur ;
    * disposer de TLSv1.2 ou de TLSv1.3 par vos clients HTTPS (note : TLSv1.3 n'est pas disponible dans les paquets CentOS 7) ;
    * utiliser une autorité de certification SSL/TLS en raccord avec les pratiques internes de votre SI (voir ci-dessous).

## Choix du type de certificat HTTPS

De façon générale, l'écosystème HTTPS nécessite la mise en place de certificats devant être reconnus et acceptés par les clients HTTPS.

Pour cela, trois options s'offrent à vous :

1. mettre en place un certificat signé par une autorité de certification tierce (Gandi, Digicert, Let's Encrypt…) ;
2. mettre en place un certificat signé par l'autorité de certification interne à votre SI (le cas échéant) ;
3. mettre en place un certificat autosigné et déployer des règles d'exception pour ce domaine dans votre SI (**non recommandé**).

Les deux premières options vous imposent de vous rapprocher d'une autorité de certification compétente et de venir compléter la configuration Nginx à l'aide des explications données ci-dessous.

La troisième option ne nécessite pas d'autorité de certification, mais de façon générale cette méthode n'est pas recommandée et elle vous impose d'ajouter une exception de sécurité sur chaque client HTTPS (à votre charge).

## Activation de la configuration HTTPS

### Ajout d'un certificat standard à Nginx

Rapprochez-vous de votre autorité de certification afin de connaître la procédure à suivre pour la génération de votre clé privée et pour l'obtention d'un certificat signé.

La configuration Nginx proposée par défaut s'attend à ce que ces fichiers soient présents aux emplacements suivants :

| Type de fichier | Emplacement |
| --------------- | ----------- |
| Certificat SSL/TLS, au format PEM | `/etc/nginx/ssl/cert.crt` |
| Clé privée, au format PEM | `/etc/nginx/ssl/key.key` |

!!! attention
    Par mesure de sécurité, veillez à ce que le répertoire `/etc/nginx/ssl` soit bien attribué à `root:root` et qu'il dispose bien de permissions restreintes `0700`.

=== "Paquets CentOS 7"

    TODO

=== "Docker Compose"

    Exemple avec `./mon_certificat.crt` comme certificat et `./ma_clef_privee.key` comme clé privée :

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
          - ./mon_certificat.cert:/etc/nginx/ssl/cert.crt:ro
          - ./ma_clef_privee.key:/etc/nginx/ssl/key.key:ro
    ```

### Ajout d'un certificat autosigné à Nginx (non recommandé)

Les commandes suivantes peuvent être utilisées afin de mettre en place un certificat autosigné sur cet environnement.

L'outil `openssl` doit être disponible.

/bin/bash: q: command not found
    Vous devrez renouveler les certificats autosignés tous les ans.

    L'utilisation de certificats autosignés provoquera l'affichage d'un message dans votre navigateur lors de la connexion à Canopsis. Vous devrez *manuellement* ajouter une exception sur chaque navigateur où vous voudrez accéder à Canopsis.

=== "Paquets CentOS 7"

    Sur l'environnement cible, exécutez la commande suivante en tant que `root`, en remplaçant `canopsis.mon-si.fr` par le vrai FQDN de votre service Canopsis :

    ```sh
    ( umask 077 ; openssl req -x509 -nodes -days 730 -newkey rsa:2048 -sha256 \
        -keyout /etc/nginx/ssl/key.key \
        -out /etc/nginx/ssl/cert.crt \
        -subj "/CN=canopsis.mon-si.fr" )
    ```

=== "Docker Compose"

    Dans le cas de certificats autosigné vous n'avez rien à faire, le conteneur Nginx va en générer automatiquement lors de son démarrage, sauf si vous surchargez les fichiers attendus à l'aide d'un volume.

## Activation de la configuration HTTPS

=== "Paquets CentOS 7"

    (TODO : la modification est manuelle)

    TODO : redémarrage du service

=== "Docker Compose"

    Pour activer HTTPS dans le conteneur vous devrez modifier ses variables d'environnement dans le fichier `compose.env`.

    | Variable | Description |
    | -------- | ----------- |
    | `CPS_SERVER_NAME` | FQDN sur lequel Canopsis sera disponible (ex : `canopsis.mon-si.fr`) |
    | `CPS_ENABLE_HTTPS` | État d'activation de HTTPS. Si cette variable vaut `true` HTTPS sera activé, dans tous les autres cas il sera désactivé |

    Exemple de `compose.env` pour activer HTTPS avec `canopsis.mon-si.fr` en FQDN :

    ```python
    CPS_SERVER_NAME=canopsis.mon-si.fr
    CPS_ENABLE_HTTPS=true
    ```

    TODO : redémarrage du service

## Utilisation d'un autre applicatif ou équipement pour servir les flux HTTPS

Dans la politique de sécurité de certains SI, la terminaison SSL/TLS doit forcément passer par un applicatif ou par un équipement dédié.

Si vous êtes dans cette situation, vous pouvez réaliser votre propre terminaison SSL/TLS avec vos propres moyens et selon vos propres conditions, mais ce type de modification n'est pas officiellement pris en charge par Canopsis. Tout éventuel ajustement à faire et toute maintenance sont donc à votre charge.
