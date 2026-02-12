# Connecteur Icinga2 vers Canopsis (connector-icinga2)

Convertit les problèmes dans Icinga2 en événements Canopsis.

## Prerequis

 - Icinga2 **2.7.0** ou supérieur

## Introdution

Le connecteur _connector-icinga2_ interroge l'API Icinga2 (port par défaut
`tcp/5665`) pour collecter les problèmes, les convertir et envoyer les
événements à Canopsis. Ce connecteur publie des messages directement dans le
bus AMQP (instance RabbitMQ de Canopsis).

## Intégration du connecteur

Le connecteur et sa documentation sont disponibles dans le dépôt
[canopsis-connectors/connector-icinga2][git-repo].

Un compte utilisateur API Icinga2 avec les permissions appropriées est
nécessaire pour authentifier les requêtes du connecteur Icinga2. Pour plus
d'information sur l'authentification avec l'API Icinga2, consulter la
[documentation Icinga][doc-icinga-api].

## Installation

Deux méthodes d'installation ou d'exécution sont proposées :

 - Installation en tant que service sur un système Linux de branche Red Hat
   (paquets RPM)
 - Exécution en tant que conteneur Docker avec l'image fournie :
   [canopsis/connector-icinga2][docker-image]

### Paquets RPM

Ajouter les dépôts Canopsis :

=== "RHEL 8"

    ```sh
    cat << EOF > /etc/yum.repos.d/canopsis.repo
    [canopsis-community]
    name = canopsis community
    baseurl=https://nexus.canopsis.net/repository/canopsis/el8/community/
    gpgcheck=0
    enabled=1
    EOF

    cat << EOF > /etc/yum.repos.d/canopsis-connectors.repo
    [canopsis-connectors]
    name=Canopsis connectors repository
    baseurl=https://nexus.canopsis.net/repository/canopsis-connectors/el8/
    gpgcheck=0
    enabled=1
    EOF
    ```

=== "RHEL 9"

    ```sh
    cat << EOF > /etc/yum.repos.d/canopsis.repo
    [canopsis-community]
    name = canopsis community
    baseurl=https://nexus.canopsis.net/repository/canopsis/el9/community/
    gpgcheck=0
    enabled=1
    EOF

    cat << EOF > /etc/yum.repos.d/canopsis-connectors.repo
    [canopsis-connectors]
    name=Canopsis connectors repository
    baseurl=https://nexus.canopsis.net/repository/canopsis-connectors/el9/
    gpgcheck=0
    enabled=1
    EOF
    ```

Installer le connecteur Icinga2 :

```sh
dnf makecache
dnf install canopsis-connector-icinga2
```

Éditer le fichier de configuration
`/etc/canopsis-connectors/icinga2/config.yml` afin de l'adapter à
l'environnement cible.

Activer le service au démarrage du système et démarrer le connecteur

```sh
systemctl enable --now canopsis-connector-icinga2.service
```

### Docker

Exemple de définition de service pour Docker Compose :

```yaml
services:
  connector-icinga2:
    image: docker.canopsis.net/docker/community/connector-icinga2:<TAG>
    volumes:
      - ./config/config.yml:/config.yml
    restart: on-failure
```

Un exemple complet de fichier [config.yml][config] est fourni avec le code du
connecteur.

## Résultat

Lorsque le connecteur démarre correctement, il indique les événements
ci-dessous dans son log (cas où les clefs `forward_ack` et `forward_downtime`
ont pour valeur `false`) :

```
> acknowledgement event stream types forward disabled
> downtime event stream types forward disabled
> connector started
```

Pour rappel, toute la documentation du connecteur Icinga2 est disponible sur le
[dépôt Git][git-repo].

[doc-icinga-api]: https://icinga.com/docs/icinga-2/latest/doc/12-icinga2-api/#authentication
[git-repo]: https://git.canopsis.net/canopsis-connectors/connector-icinga2
[docker-image]: https://git.canopsis.net/docker/community/container_registry/305
[config]: https://git.canopsis.net/canopsis-connectors/connector-icinga2/-/blob/main/config/config.yml
