# Connecteur prometheus

## Description

Convertit des alertes de l'Alertmanager de Prometheus en évènements Canopsis.

## Principe de fonctionnement

L'Alertmanager de Prometheus peut être configuré pour envoyer les alertes
via un [*webhook*][webhook].

Le contenu du message du *webhook* n'est pas configurable dans Prometheus ;
l'intégration vers chaque outil de destination nécessite donc un programme
adéquat pour transformer et envoyer les données.

Dans ce contexte, le connecteur **prometheus** est un programme qui :

- *Écoute* les requêtes HTTP POST envoyées par l'Alertmanager de Prometheus

    Port d'écoute par défaut : 8080/tcp (configurable)

- *Lit et transforme* les données reçues pour en faire des évènements Canopsis

    Le message envoyé par l'Alertmanager (JSON) est décodé et un évènement
    Canopsis est construit pour chaque alerte. Le contenu placé dans l'évènement
    Canopsis est une combinaison de constantes et de valeurs tirées de
    l'alerte Prometheus. Cette association est configurable.

- *Envoie* les évènements à Canopsis, c'est-à-dire publie des messages dans le
  bus AMQP (instance RabbitMQ de Canopsis)

## Intégration du connecteur

Le connecteur et sa documentation (installation, configuration, utilisation)
sont disponibles dans le dépôt
[canopsis-connectors/connector-prometheus][upstream].

### Installation

Deux méthodes d'installation ou d'exécution sont proposées :

- Installation en tant que service sur un système de production
- Exécution en tant que conteneur *Docker* avec l'image fournie :
  [canopsis/canopsis-connector-prometheus][docker-image]

### Configuration connecteur

Quelle que soit la méthode d'installation choisie, la configuration du
connecteur passe par le renseignement du fichier `config.yml`, qui sert à :

- Indiquer l'URL AMQP où le connecteur doit envoyer les évènements
- Définir le token de securité du webhook
- Définir le contenu des évènements et la correspondance des attributs
  (alerte Prometheus -> évènement Canopsis)

Il est également possible de changer le port d'écoute avec l'argument `--port`
du connecteur.

Un exemple complet de fichier `config.yml` est fourni avec le code du
connecteur.

L'installation et la configuration du connecteur sont documentées
[à la racine du dépôt][upstream].

### Configuration Prometheus

Côté Prometheus, dans la configuration du daemon Alertmanager, il faut définir
un *receiver* avec une configuration *webhook* en indiquant l'URL du connecteur.

Exemple (pour une installation du connecteur sur le même serveur que
l'Alertmanager) :

```yaml
receivers:
- name: 'prometheus'
  webhook_configs:
  - url: 'http://127.0.0.1:8080/webhook'
    http_config:
      bearer_token: 
```

Note : l'extrait ci-dessus est un exemple, à intégrer au sein de votre propre
configuration alertmanager. Pour que le connecteur prometheus reçoive
des messages, le *receiver* nommé doit être utilisé dans votre définition du
routage des notifications. Voir pour cela la page sur la
[configuration de l'alertmanager][alertmanager-config] dans la documentation
officielle de Prometheus.

[upstream]: https://git.canopsis.net/canopsis-connectors/connector-prometheus
[alertmanager-config]: https://prometheus.io/docs/alerting/latest/configuration/
[docker-image]: https://git.canopsis.net/docker/community/container_registry/209
[webhook]: https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
