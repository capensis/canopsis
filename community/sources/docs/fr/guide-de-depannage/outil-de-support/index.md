# Outil de support

Lorsqu'un cas de support se présente sur un environnement Canopsis, plusieurs pistes peuvent être explorées.  
Il vous est alors demandé de consulter 

* Les logs des moteurs
* Les dernières modifications réalisées au niveau des règles de gestion (comportements périodiques, règles d'eventfilter, scénarios, etc.)
* Des données en base
* L'état des files d'attentes du bus RabbitMQ
* D'autres éléments

Canopsis met désormais à disposition un outil d'aide au support dont les premières fonctionnalités sont :

* Export de données MongoDB et TimescaleDB
* Export des [Métriques techniques](../../guide-de-depannage/metriques-techniques/)
* Export de métriques spécifiques à RabbitMQ
* Visualisation des données exportées sous forme de dashboards Grafana

![dashboard_tech_metrcis](./img/dashboard_tech_metrics_1.png)


[TOC]

## Types de données exportées

### Action log

Toutes les modifications de règles de gestion au sein de Canopsis sont soumises à des enregistrements dans une table TimescaleDB `action_log`.  

Les données sont représentées avec la structure suivante :

|  Column    |            Type             | Description                                       |
| ---------- | --------------------------- | ------------------------------------------------- |
| id         | bigint                      | id interne                                        |
| type       | smallint                    | type                                              |
| value_type | character varying(255)      | Nom usuel du type                                 |
| value_id   | character varying(255)      | L'id de l'objet manipulé                          |
| author     | character varying(255)      | L'id de l'auteur qui a manipulé l'objet           |
| time       | timestamp without time zone | les dates et heures de la modification de l'objet |
| data       | jsonb                       | Un JSON représentant les données manipulées       |

Voici un exemple d'enregistrement d'un "action_log" :

```sql
SELECT * FROM action_log;
```

!!! tip "Astuce"
    Si vous utilisez un environnement Docker, vous pouvez exécuter cette requête avec la commande suivante : 
    
    ```sh
    CPS_EDITION=pro docker compose exec timescaledb psql postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis -c "SELECT * FROM action_log"
    ```


| id | type | value_type  | value_id | author                               | time                | data |
| -- | ---- | ----------- | -------- | ------------------------------------ | ------------------- | ---- | 
| 14 | 0    | eventfilter | regle1   | e641bcd8-16d4-471d-be36-78061c9a3004 | 2024-10-30 13:12:22 | {"_id": "regle1", "type": "drop", "rrule": "", "author": "e641bcd8-16d4-471d-be36-78061c9a3004", "config": {}, "created": 1730293942, "enabled": true, "exdates": [], "updated": 1730293942, "priority": 0, "exceptions": [], "description": "regle1", "event_pattern": [[{"cond": {"type": "eq", "value": "aa"}, "field": "component"}]], "external_data": null, "entity_pattern": [], "resolved_exdates": null, "corporate_entity_pattern": "", "corporate_entity_pattern_title": ""} |

Voici une liste non exhaustive des différents "value_type" que vous pourrez retrouver dans l'export :

"kpi_filter", "statesetting", "resolverule", "eventfilter", "role", "map", "pattern", "eventrecord", "jobconfig", "colortheme", "idlerule", "entity", "pbehavior", "pbehaviorreason", "metaalarmrule", "user", "instruction", "entityservice", "declareticketrule", "alarmtag", "job", "widget", "viewtab", "entitycategory", "widgettemplate", "broadcastmessage", "pbehaviorexception", "linkrule", "snmprule", "flappingrule", "dynamicinfo", "viewgroup", "widgetfilter", "pbehaviortype", "view", "playlist", "scenario"

### Métriques techniques

Les Métriques techniques sont décrites dans cette [documentation](../../guide-de-depannage/metriques-techniques/).  

### Métriques RabbitMQ

La récupération des métriques spécifiques à RabbitMQ s'appuie sur une service tiers, Prometheus. 
Ce service Prometheus n'est pour le moment pas intgéré de base dans une distribution Canopsis.  

Voici néanmoins un exemple de configuration de celui-ci :

```yml
volumes:
  prometheusdata:
    driver: local

services:
  prometheus:
    image: prom/prometheus:${PROMETHEUS_TAG}
    container_name: prometheus_dev
    ports:
      - "9090:9090"
    volumes:
      - prometheusdata:/prometheus
      - ./files-pro/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      # should explicitly define paths when use command config, it fails instead.
      - --config.file=/etc/prometheus/prometheus.yml
      - --storage.tsdb.path=/prometheus
      - --web.enable-admin-api
```

Le fichier `./files-pro/prometheus/prometheus.yml` doit contenir les éléments suivants :

```yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'rabbitmq'
    static_configs:
      - targets:
          - 'rabbitmq:15692'
    metric_relabel_configs:
      - source_labels: [__name__, queue]
        regex: 'rabbitmq_resident_memory_limit_bytes;(.*)|rabbitmq_process_resident_memory_bytes;(.*)|rabbitmq_disk_space_available_bytes;(.*)|rabbitmq_channel_messages_acked_total;Engine.*|rabbitmq_queue_messages_published_total;Engine.*|rabbitmq_channel_messages_delivered_ack_total;Engine.*|rabbitmq_queue_messages_ready;Engine.*|rabbitmq_queue_messages_unacked;Engine.*|rabbitmq_queue_messages_ready_bytes;Engine.*|rabbitmq_queue_messages_unacked_bytes;Engine.*'
        action: keep
```

Le fichier de configuration de RabbitMQ doit également être modifié pour prendre en compte la directive "prometheus.return_per_object_metrics".  
Le fichier `./files-pro/rabbitmq/rabbitmq.conf` doit donc contenir : 

```
prometheus.return_per_object_metrics = true`
```

La section du service "rabbitmq" doit également être modifiée pour prendre en compte le fichier de configiration RabbitMQ.

```yml
  rabbitmq:
    <<: *backend_base
    image: rabbitmq:${RABBITMQ_TAG}
    # Prevent rabbit from loosing data on stack down and up
    # See : https://git.canopsis.net/canopsis/canopsis-pro/-/issues/4005
    hostname: rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
    env_file:
      - rabbitmq.env
    volumes:
      - rabbitmqdata:/var/lib/rabbitmq
      - ./files-pro/rabbitmq/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf
    oom_score_adj: -1000
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 30s
      timeout: 30s
      retries: 3
```

Grâce à ces configurations, vous serez en mesure d'activer l'export des métriques RabbitMQ.

## Export de données

### Description des options d'export

L'outil d'export peut être lancé avec plusieurs options qui sont présentées ci-après :


* `--without-action-log`

    -	**Action:** Exclut les "action log" de l'export.
    -	**Description:** Les "action log" correspondent à l'ensemble des actions réalisées via l'API de Canopsis concernant les règles de gestion. On y retrouve par exemple les créations, modifications, suppressions des comportements périodiques, des règles d'enrichissement, des scénarios, etc.

* `--without-tech-metrics`

    -	**Action:** Exclut les "métriques techniques" de l'export.
    -	**Description:** Les [Métriques techniques](../../guide-de-depannage/metriques-techniques/) permettent d'effectuer des analyses poussées sur les temps de traitement des événements dans chaque moteur, les consommations de ressources, etc. 

* `--without-rabbitmq-metrics`

    -	**Action:** Exclut les "métriques spécifiques à RabbitMQ" de l'export.
    -	**Description:** Des métriques spécifiques à RabbitMQ peuvent être récupérées depuis un backend Prometheus (Voir section à compléter).

* `--with-mongodb-dump`

    -	**Action:** Inclut un un dump MongoDB dans l'export.
    -	**Description:** Un dump MongoDB est réalisé à l'aide de la commande `mongodump`. Ce dump n'est pas anonymisé et peut contenir des informations sensibles.

* `--exclude-mongodb-collections`

    -	**Action:** Exclut des collections spécifiques du dump MongoDB. Il doit être utilisé en même temps que l'option `--with-mongodb-dump`.
    -	**Description:** Il s'agit de définir une liste de collections à exclure du dump qui sera exporté.
    -	**Exemple:**
    ```yaml
    services:
      export:
        ...
        volumes:
          - ./export-support:/export
        command:
          - --with-mongodb-dump
          - --exclude-mongodb-collections=collection1,collection2
    ```

* `--with-postgres-dump`

    -	**Action:** Inclut un dump TimescaleDB dans l'export.
    -	**Description:** Un dump TimescaleDB est réalisé à l'aide de la commande `pg_dump`. Ce dump n'est pas anonymisé et peut contenir des informations sensibles.

* `--exclude-postgres-tables`

    -	**Action:** Exclut des collections spécifiques du dump TimescaleDB. Il doit être utilisé en même temps que l'option `--with-postgres-dump`.
    -	**Description:** Il s'agit de définir une liste de collections à exclure du dump qui sera exporté.
    -	**Exemple:**
   ```yaml
    services:
      export:
        ...
        volumes:
          - ./export-support:/export
        command:
          - --with-postgres-dump
          - --exclude-postgres-tables=table1,table2
   ```

### Variables d'environnement

| Variable              | Description |
| --------------------- | ----------- |
| CPS_POSTGRES_URL      | Chaine de connexion PostgreSQL, obligatoire |
| CPS_MONGO_URL         | Chaine de connexion MongoDB, obligatoire si l'option --with-mongodb-dump est utilisée |
| CPS_POSTGRES_TECH_URL | Chaine de connexion PostgreSQL pour les métriques tehcniques, obligatoire lorsque l'option --without-tech-metrics n'est pas utilisée |
| CPS_PROMETHEUS_URL    | Chaine de connexion Prometheus, obligatoire lorsque l'option --without-rabbitmq-metrics n'est pas utilisée |
| PROMETHEUS_DATA_DIR   | Répertoire des données Prometheus, obligatoire lorsque l'option --without-rabbitmq-metrics n'est pas utilisée |


### Lancement d'un export

=== "Docker Compose"

    Le fichier de référence `docker-compose.override.yml` contient la section suivante 
    
    ```yml
    export:
        <<: *initial_config_base
        profiles:
          - support
        image: ${DOCKER_REPOSITORY}${CPS_EDITION}/canopsis-support-tool-export:${CANOPSIS_IMAGE_TAG}
        volumes:
          - ./export-support:/export
        command:
          - --with-postgres-dump
          - --with-mongodb-dump
          - --without-rabbitmq-metrics
    ```

    Le service d'export peut donc être lancé de la manière suivante :

    ```sh
    CPS_EDITION=pro docker compose --profile support up -d export
    ```
    
    Le résultat de l'opération sera visible dans le répertoire `export-support` dans lequel vous trouverez une archive au format tar.gz.
    
    Si vous avez activé les [Métriques RabbitMQ](#metriques-rabbitmq), alors vous devez modifier la configuration du service d'export comme suit : 

    ```yml
    export:
        <<: *initial_config_base
        profiles:
          - support
        image: ${DOCKER_REPOSITORY}${CPS_EDITION}/canopsis-support-tool-export:${CANOPSIS_IMAGE_TAG}
        environment:
          - PROMETHEUS_DATA_DIR=/prometheus    
          - CPS_PROMETHEUS_URL=http://prometheus:9090
        volumes:
          - ./export-support:/export
          - prometheusdata:/prometheus
        command:
          - --with-postgres-dump
          - --with-mongodb-dump
          #- --without-rabbitmq-metrics
    ```


=== "Paquets RHEL 8"

    A venir

=== "Helm"

    A venir


## Exploitation des résultats

### Import de données

#### Environnement d'import

Les archives précédemment créées par l'outil d'export doivent être importées sur un environnement hors Canopsis.  
L'exploitation des données et en effet totalement indépendante.

L'environnement d'exploitation est composé de :

* Grafana
* Prometheus
* TimescaleDB
* L'outil d'import lui-même

Voici un fichier de configuration "docker compose" permettant d'instancier cet environnement.

```yml
volumes:
  grafanadata:
    driver: local
  grafanaconfig:
    driver: local
  prometheusdata:
    driver: local
  timescaledata:
    driver: local

services:
  import:
    image: ${DOCKER_REPOSITORY}/pro/canopsis-support-tool-import:${CANOPSIS_IMAGE_TAG}
    container_name: import
    depends_on:
      - prometheus
      - grafana
      - timescaledb
    env_file: .env
    environment:
      - PROMETHEUS_DATA_DIR=/prometheus
    volumes:
      - prometheusdata:/prometheus
      - ./files/archive:/archive
  prometheus:
    image: prom/prometheus:${PROMETHEUS_TAG}
    container_name: prometheus_support
    ports:
      - "9092:9090"
    volumes:
      - prometheusdata:/prometheus
  grafana:
    image: grafana/grafana:${GRAFANA_TAG}
    ports:
      - "3002:3000"
    volumes:
      - grafanadata:/var/lib/grafana
      - grafanaconfig:/etc/grafana
      - ./files/grafana/provisioning:/etc/grafana/provisioning
      - ./files/grafana/dashboards:/etc/dashboards
    env_file: .env
    restart: unless-stopped
  timescaledb:
    image: timescale/timescaledb:${TIMESCALEDB_TAG}
    env_file: .env
    environment:
      - PGPORT=5435
    ports:
      - "5435:5435"
    volumes:
      - timescaledata:/var/lib/postgresql/data
    restart: unless-stopped
```

Le fichier d'environnement associé est le suivant (.env) :

```
COMPOSE_PROJECT_NAME=canopsis-support

GF_SECURITY_ADMIN_PASSWORD=canopsis

POSTGRES_USER=cpspostgres
POSTGRES_PASSWORD=canopsis
POSTGRES_DB=canopsis

POSTGRES_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@timescaledb:5435/${POSTGRES_DB}

TIMESCALEDB_TAG=2.15.1-pg15
PROMETHEUS_TAG=v2.53.0
GRAFANA_TAG=10.0.10

CANOPSIS_IMAGE_TAG=24.10.0
DOCKER_REPOSITORY=docker.canopsis.net/docker
```

#### Lancement d'un import

L'archive précédemment exportée doit être placée dans le répertoire `./files/archive`.  
**Seule une archive** peut être présente dans ce répertoire.  

Notez également que l'import d'un snapshot Prometheus peut prendre plus d'une minute dans certains cas. 
Si le dashboard RabbitMQ est vide, il se peut que le snapshot soit en cours d'import.

```sh
docker compose up -d
```

!!! warning "Avertissement"
    Seul un export peut être traité à la fois.  
    Veillez à détruire un environnement existant avant de démarrer une nouvelle analyse.


### Dashboards

Lorsque vous avez démarré l'environnement d'import, un service Grafana a été démarré sur le port 3002.  
Vous pouvez y accéder avec votre navigateur.  

L'authentification suivante est requise : admin/canopsis

Une fois authentifié, vous accédez à un répertoire de dahsboards `support` dans lequel 3 dashboards sont disponibles.

**Action Log**

![dashboard_action-log](./img/dashboard_action-log.png)


**RabbitMQ**

![dashboard_rabbitmq](./img/dashboard_rabbitmq.png)

**Tech metrics**

![dashboard_tech-metrics](./img/dashboard_tech-metrics.png)
