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

    Port d'écoute par défaut : `8080/tcp` (configurable)

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

Un jeton est requis pour la connexion entre le connecteur et prometheus alertmanager.

Les variables sont `PROMETHEUS_BEARER_TOKEN` et `bearer_token`, elles ont besoin d'une chaîne identique sur ces 2 variables pour initier l'identification.

### Installation

Deux méthodes d'installation ou d'exécution sont proposées :

- Installation en tant que service sur un système de production
- Exécution en tant que conteneur *Docker* avec l'image fournie :
  [canopsis/canopsis-connector-prometheus][docker-image] et l'exemple de configuration [docker compose](https://git.canopsis.net/canopsis-connectors/connector-prometheus#deployment-docker-compose)

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
      bearer_token: token_example
```
!!! Note 

    L'extrait ci-dessus est un exemple, à intégrer au sein de votre propre configuration alertmanager. 

    Pour que le connecteur prometheus reçoive des messages, le *receiver* nommé doit être utilisé dans votre définition du
    routage des notifications. Voir pour cela la page sur la [configuration de l'alertmanager][alertmanager-config] dans la documentation
    officielle de Prometheus.

### Configuration connecteur

Quelle que soit la méthode d'installation choisie, la configuration du
connecteur passe par le renseignement d'un fichier `config.yml` à passer en paramètre au démarrage du connecteur et qui sert à :

- Indiquer l'URL AMQP où le connecteur doit envoyer les évènements
- Définir le token de securité du webhook
- Définir le contenu des évènements et la correspondance des attributs
  (alerte Prometheus -> évènement Canopsis)

Il est également possible de changer le port d'écoute avec l'argument `--port`
du connecteur.

Un exemple complet de fichier [config.yml](https://git.canopsis.net/canopsis-connectors/connector-prometheus/-/blob/master/config/config.yml) est fourni avec le code du
connecteur.

Pour ajouter des informations dans la configuration de votre connecteur depuis Prometheus, consultez alertmanager (ici : `<alertmanager_IP>:9093/api/v1/alerts` ) afin d'identifier les informations sur lesquelles vous pouvez réaliser des mappings.
Sur la base des informations fournies par l'alertmanager, voicii un exemple de données renvoyées par Prometheus pour illustrer notre exemple de configuration :

```yaml
{
  "status": "success",
  "data": [
    {
      "labels": {
        "alertname": "InstanceDown",
        "instance": "node_exporter:9100",
        "job": "node1",
        "severity": "critical"
      },
      "annotations": {
        "description": "node_exporter:9100 of job node1 has been down for more than 1 minute.",
        "title": "Instance node_exporter:9100 down"
      },
      "startsAt": "2023-04-11T16:29:54.61318652Z",
      "endsAt": "2023-04-12T07:44:59.61318652Z",
      "generatorURL": "http://7e9d35e4e800:9090/graph?g0.expr=up+%3D%3D+0&g0.tab=1",
      "status": {
        "state": "active",
        "silencedBy": [],
        "inhibitedBy": []
      },
      "receivers": [
        "connector-prometheus"
      ],
      "fingerprint": "ee8077593a410082"
    }
  ]
}
```


#### Directives liées aux `Output`

`output_length` et `long_output_length` prennent comme valeur le nombre maximum de caractères à partir duquel l'`output` et le `long_output` de l'événement seront tronqués.

```yaml
output_length: 255
long_output_length: 1024
```

#### Directives liées aux mapping des champs

Le connecteur permet de réaliser différents mapping entre les informations provenant de Prometheus et celles envoyées dans l'`event` généré vers Canopsis

!!! Warning
    
    Les valeurs de `event_type` et `source_type` ne peuvent pas être modifiés car le connecteur ne peut envoyer que des événements de type `check` et ayant comme `source_type` la valeur `resource`.
    
    Les valeurs des champs `connector` et `connector_name` seront toujours à spécifier de façon statique
    
    ```yaml
    # connector specifies a type of the connector.
    connector: prometheus
    # connector_name specifies a unique name of the connector.
    connector_name: node_exporter1
    ```


##### type : copy

Le type `copy` est utilisé pour récupérer la valeur statique renvoyée par prometheus.

```yaml
component:
  type: copy
  value: labels.instance
```

Le résultat obtenu et envoyé à Canopsis sera :

`"component": "node_exporter:9100"`

##### type : set

Le type `set` permet de définir une chaîne de valeur constante.

```yaml
component:
  type: set
  value: My Component
```

Le résultat obtenu et envoyé à Canopsis sera :

`"component": "My Component"`

##### type : template

L'utilisation du template go est fournie par les fonctions suivantes :

- `lowercase` : permet de mettre en minuscules les lettres dans une sortie attendue.

- `uppercase` : permet de mettre en majuscules les lettres dans une sortie attendue.

- `replace` : permet de remplacer un champ correspondant à une regex par un autre

- `trim` : supprime les espaces blancs initiaux et finaux d'une chaîne de caractères.

- `split` : permet de récupérer une sous-chaîne spécifique d'une chaîne qui est séparée en plusieurs parties par une chaîne spécifique.

- `regex_map_keys` : permet de récupérer la valeur associée à la première clé d'un dictionnaire qui correspond à une expression régulière spécifiée.

- `map_as_keys` : utilisé pour vérifier si une carte contient une clé spécifiée et retourner un booléen en conséquence.

- `json` : utilisé pour convertir n'importe quelle valeur en chaîne JSON et la retourner.

###### Exemple de template

* Méthode `uppercase`

```yaml
output:
  type: template
  field: labels.severity
  value: Statut : {{ uppercase .Field }}
```

Le résultat obtenu et envoyé à Canopsis sera :

`"output": "Statut : CRITICAL"`

* Capture group via une `regexp`

```yaml
output:
  type: template
  field: annotations.title
  regexp: Instance (?P<Substr>.* ?)($|,)
  value: "Prometheus {{ .RegexMatch.Substr }}"
```
Le résultat obtenu et envoyé à Canopsis sera :

`"output": "Prometheus node_exporter:9100 down"`

##### Utilisation d'un Préfixe

Pour éviter les collisions avec les champs internes à Cnaopsis, il est possible d'utiliser l'option de `prefix` sur les données de type `extra_infos` :

Vous pouvez utiliser un préfixe pour les informations supplémentaires contenusi dans la variable : `extra_infos_prefix:` qui a pour but de changer le nom des champs, mais si vous décidez de ne pas l'utiliser, un avertissement apparaît dans les logs lors du lancement du conteneur.

!!! Note

    Logs sans préfixe

    ```go
    canopsis-pro-connector_prometheus-1 | 2023-04-11T15:52:14Z WRN app/cmd/main.go:125 > extra_infos_prefix is empty
    canopsis-pro-connector_prometheus-1 | 2023-04-11T15:52:14Z INF app/cmd/main.go:188 > connector started
    ```

###### Exemple de préfixe

```yaml
extra_infos:
   type_ack:
     type: set
     valeur: auto
extra_infos_prefix: prom_
```
Le résultat obtenu et envoyé à Canopsis sera :

`"prom_type_ack": "auto"`

L'installation et la configuration complète du connecteur sont documentées
[à la racine du dépôt][upstream].

[upstream]: https://git.canopsis.net/canopsis-connectors/connector-prometheus
[alertmanager-config]: https://prometheus.io/docs/alerting/latest/configuration/
[docker-image]: https://git.canopsis.net/docker/community/container_registry/209
[webhook]: https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
