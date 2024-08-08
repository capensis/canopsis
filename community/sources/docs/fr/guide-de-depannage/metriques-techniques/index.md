# Métriques techniques

Lorsqu'un problème de performance est identifié, il peut être assez complexe de trouver et de comprendre son origine.  
Canopsis met à disposition un ensemble de métriques qui peuvent orienter le diagnostic.  

Pour l'heure, ces métriques peuvent vous être demandées par le service de support.  

## Mise en oeuvre

### Remarques sur la mise en place des métriques techniques

L'exemple qui suit relate la mise en place sur une installation mono-noeud. Dans le cadre d'une installation
multi-noeud, veuillez porter une attention particulière aux points suivants:

* La génération d'un export depuis la webui nécessite la présence du binaire `pg_dump` du paquet `postgresql` sur 
l'hôte exécutant le moteur `canopsis-api`.
* L'exécution de la commande `canopsis-reconfigure` nécessite la présence du binaire `mongo` du paquet `mongodb-org-shell` 
sur l'hôte.
* L'accès à la base de données PostgreSQL distante est contrainte par la configuration du fichier `pg_hba.conf` de celui-ci.

### Préparation de la base de données

Les données collectées sont stockées dans une base `timescaledb` et donc `postgresql`.  

!!! warning

    Cette base de donnée **DOIT** être différente de celle utilisée pour les KPI Canopsis.

=== "Docker Compose"

    ```sh
    CPS_EDITION=pro docker compose exec timescaledb psql postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/postgres
    postgres=# CREATE database canopsis_tech_metrics;
    postgres=# \c canopsis_tech_metrics
    canopsis_tech_metrics=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis_tech_metrics=# SET password_encryption = 'scram-sha-256';
    canopsis_tech_metrics=# CREATE USER cpspostgres_tech_metrics WITH PASSWORD 'canopsis';
    canopsis_tech_metrics=# exit
    ```

    Définir la chaine de connexion à cette base de données dans le fichier `compose.env`

    ```
    CPS_POSTGRES_TECH_URL=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis_tech_metrics
    ```

=== "Paquets RHEL 8"

    ```sh
    psql postgresql://cpspostgres_tech_metrics:canopsis@timescaledb:5432/postgres
    postgres=# CREATE database canopsis_tech_metrics;
    postgres=# \c canopsis_tech_metrics
    canopsis_tech_metrics=# CREATE EXTENSION IF NOT EXISTS timescaledb;
    canopsis_tech_metrics=# SET password_encryption = 'scram-sha-256';
    canopsis_tech_metrics=# CREATE USER cpspostgres_tech_metrics WITH PASSWORD 'canopsis';
    canopsis_tech_metrics=# exit
    ```

    Définir la chaine de connexion à cette base de données dans le fichier `/opt/canopsis/etc/go-engines-vars.conf

    ```
    CPS_POSTGRES_TECH_URL=postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis_tech_metrics
    ```

### Activation des collectes

La collecte ne sera active que lorsque l'option sera activée dans le fichier [canopsis.toml](../../../guide-administration/administration-avancee/modification-canopsis-toml/#section-canopsistech_metrics) 

```
[Canopsis.tech_metrics]
Enabled = true
DumpKeepInterval = "1h"
```

L'option `DumpKeepInterval` définit le délai au dela duquel un dump demandé depuis l'interface graphique sera supprimé.

=== "Docker Compose"

    Surchargez l'entrée du service `reconfigure` dans votre fichier d'override `docker-compose.override.yml` comme suit (notez l'activation de `-migrate-tech-postgres`) :

    ```sh
    reconfigure:
      command: /canopsis-reconfigure -migrate-postgres=true -migrate-mongo=true -migrate-tech-postgres -edition ${CPS_EDITION} -conf /canopsis-${CPS_EDITION}.toml -override /opt/canopsis/etc/conf.d/canopsis-override.toml
      volumes:
        - ./files-pro/reconfigure/reconfigure.override.toml:/opt/canopsis/etc/conf.d/canopsis-override.toml
    ```

    Ajouter le bloc de configuration dans votre fichier d'override `files-pro/reconfigure/reconfigure.override.toml` puis exécutez un `up -d`.

    ```sh
    CPS_EDITION=pro docker compose up -d
    ```

=== "Paquets RHEL 8"

    Ajouter le bloc de configuration dans votre fichier d'override `canopsis.override.yml` puis exécutez `canopsis-reconfigure`

    ```sh
    set -o allexport ; source /opt/canopsis/etc/go-engines-vars.conf
    /opt/canopsis/bin/canopsis-reconfigure -migrate-postgres=true -migrate-mongo=true -edition pro -migrate-tech-postgres
    ```

### Génération d'un export à destination du support

Le service de support peut être amené à vous demander un export des métriques techniques.  
Pour le réaliser, rendez-vous dans le menu `Administration->Bilan de santé` dans l'onglet `Métriques des moteurs`.

![tech_metrics1](./img/tech_metrics1.png)
![tech_metrics2](./img/tech_metrics2.png)
![tech_metrics3](./img/tech_metrics3.png)

## Description des métriques

<table>
   <tr>
      <th>
         Métrique
      </th>
      <th>
         Description
      </th>
      <th>
         SQL associé
      </th>
   </tr>
   <tr>
      <td>
         fifo_queue
      </td>
      <td>
         Cette métrique permet d’investiguer sur des problèmes de performance liés au flux.
         Collectée une fois par minute
         43.200 rows per month = 60 minutes * 24 hours * 30 days
      </td>
      <td>
         <pre>
   CREATE TABLE IF NOT EXISTS fifo_queue
    (
        time   TIMESTAMP NOT NULL,
        length INT       NOT NULL
    );
    SELECT create_hypertable('fifo_queue', 'time', if_not_exists => TRUE);
    SELECT add_retention_policy('fifo_queue', INTERVAL '30 days');
        </pre>
      </td>
   </tr>
   <tr>
      <td>
         cps_event
      </td>
      <td>
         Temps de traitement depuis engine-fifo jusque engine-action.  
         Cette métrique est représentative de la performance réelle de traitement. La consommation d’espace disque anticipée est très importante, nous limitons cette métrique à quelques jours.
      </td>
      <td>
         <pre>
    CREATE TABLE IF NOT EXISTS cps_event
    (
        time     TIMESTAMP   NOT NULL,
        type     VARCHAR(30) NOT NULL,
        interval INT         NOT NULL
    );
    SELECT create_hypertable('cps_event', 'time', if_not_exists => TRUE);
    
    CREATE MATERIALIZED VIEW cos_event_summary_hourly
                (time, type, count, interval)
                WITH (timescaledb.continuous)
    AS
    SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
    FROM cps_event
    GROUP BY time_bucket('1 hour', time), type;
    
    SELECT add_continuous_aggregate_policy('cps_event_summary_hourly', '12 hours', '1 hour', '1 hour');
    
    CREATE MATERIALIZED VIEW _event_summary_daily
                (day, type, count, interval)
                WITH (timescaledb.continuous)
    AS
    SELECT time_bucket('1 day', time), type, count(*), avg(interval)
    FROM cps_event
    GROUP BY time_bucket('1 day', time), type;
    
    SELECT add_continuous_aggregate_policy('cps_event_summary_daily', '49 hours', '1 hour', '1 hour');
    
    SELECT add_retention_policy('cps_event', INTERVAL '49 hours');
    SELECT add_retention_policy('cps_event_summary_hourly', INTERVAL '7 days');
 
        </pre>
      </td>
   </tr>
   <tr>
      <td>
         fifo_event
      </td>
      <td>
         Temps de traitement engine-fifo.  
         Cette métrique aide au diagnostic de bottleneck
      </td>
      <td>
         <pre>
     CREATE TABLE IF NOT EXISTS fifo_event
    (
        time     TIMESTAMP   NOT NULL,
        type     VARCHAR(30) NOT NULL,
        interval INT         NOT NULL
    );
    SELECT create_hypertable('fifo_event', 'time', if_not_exists => TRUE);
</pre>
      </td>
   </tr>
   <tr>
      <td>
         che_event
      </td>
      <td>
         Temps de traitement engine-che.  
         Pour cette métrique, l’accent est mis sur la fréquence de création ou de mise à jour des entités dans le context-graph. Cela devrait permettre à terme de pouvoir proposer des optimisations algorithmiques du moteur.
      </td>
      <td>
         <pre>
    CREATE TABLE IF NOT EXISTS che_event
    (
        time                TIMESTAMP   NOT NULL,
        interval            INT         NOT NULL,
        type                VARCHAR(30) NOT NULL,
        entity_type         VARCHAR(30) NOT NULL,
        is_new_entity       BOOLEAN     NOT NULL,
        is_infos_updated    BOOLEAN     NOT NULL,
        is_services_updated BOOLEAN     NOT NULL
    );
    SELECT create_hypertable('che_event', 'time', if_not_exists => TRUE);
</pre>
      </td>
   </tr>
   <tr>
      <td>
         axe_event
      </td>
      <td>
         Temps de traitement engine-axe.  
         On regarde ici les temps machines occupés par la gestion d’alarme pure. Il s’agit de repérer quelles opérations liées aux alarmes sont consommatrices et amener une comparaison avec les autres moteurs de la chaîne.
      </td>
      <td>
         <pre>
    CREATE TABLE IF NOT EXISTS axe_event
    (
        time              TIMESTAMP   NOT NULL,
        interval          INT         NOT NULL,
        type              VARCHAR(30) NOT NULL,
        entity_type       VARCHAR(30) NOT NULL,
        alarm_change_type VARCHAR(30) NOT NULL
    );
    SELECT create_hypertable('axe_event', 'time', if_not_exists => TRUE);
</pre>
      </td>
   </tr>
   <tr>
      <td>
         correlation_event, service_event, dynamic_infos_event, action_event
      </td>
      <td>
         Temps de traitement engine-correlation, engine-service, engine-dynamic-infos, engine-action.  
         Il s’agit de contrôler le pourcentage de temps de traitement occupé par les corrélations d’événements, les services, les enrichissements d’alarmes ainsi que les scénarios.
      </td>
      <td>
         <pre>
    CREATE TABLE IF NOT EXISTS correlation_event
    (
        time     TIMESTAMP   NOT NULL,
        type     VARCHAR(30) NOT NULL,
        interval INT         NOT NULL
    );
    SELECT create_hypertable('correlation_event', 'time', if_not_exists => TRUE);


    CREATE TABLE IF NOT EXISTS service_event
    (
        time     TIMESTAMP   NOT NULL,
        type     VARCHAR(30) NOT NULL,
        interval INT         NOT NULL
    );
    SELECT create_hypertable('service_event', 'time', if_not_exists => TRUE);


    CREATE TABLE IF NOT EXISTS dynamic_infos_event
    (
        time     TIMESTAMP   NOT NULL,
        type     VARCHAR(30) NOT NULL,
        interval INT         NOT NULL
    );
    SELECT create_hypertable('dynamic_infos_event', 'time', if_not_exists => TRUE);


    CREATE TABLE IF NOT EXISTS action_event
    (
        time     TIMESTAMP   NOT NULL,
        type     VARCHAR(30) NOT NULL,
        interval INT         NOT NULL
    );
    SELECT create_hypertable('action_event', 'time', if_not_exists => TRUE);
</pre>
      </td>
   </tr>
   <tr>
      <td>
         axe_periodical
      </td>
      <td>
         Le moteur engine-axe émet énormément d’événements, cette métrique permet d’apporter un diagnostic rapide. (environ 10.080 enregistrements/mois = 60 minutes * 24 hours * 7 days).
      </td>
      <td>
         <pre>
    CREATE TABLE IF NOT EXISTS axe_periodical
    (
        time     TIMESTAMP NOT NULL,
        interval INT       NOT NULL,
        events   INT       NOT NULL
    );
    SELECT create_hypertable('axe_periodical', 'time', if_not_exists => TRUE);
    SELECT add_retention_policy('axe_periodical', INTERVAL '7 days');
</pre>
      </td>
   </tr>
   <tr>
      <td>
         pbehavior_periodical
      </td>
      <td>
         Le moteur pbehavior est souvent victime de dégradation de performance, cette métrique viendra en aide.
      </td>
      <td>
         <pre>
    CREATE TABLE IF NOT EXISTS pbehavior_periodical
    (
        time       TIMESTAMP NOT NULL,
        interval   INT       NOT NULL,
        events     INT       NOT NULL,
        entities   INT       NOT NULL,
        pbehaviors INT       NOT NULL
    );
    SELECT create_hypertable('pbehavior_periodical', 'time', if_not_exists => TRUE);
    SELECT add_retention_policy('pbehavior_periodical', INTERVAL '7 days');
</pre>
      </td>
   </tr>
   <tr>
      <td>
         che_infos
      </td>
      <td>
         Cette métrique permet de déterminer une stratégie d’optimisation selon la fréquence de mise à jour des entités.
      </td>
      <td>
         <pre>
        CREATE TABLE IF NOT EXISTS che_infos
        (
            time TIMESTAMP   NOT NULL,
            name VARCHAR(30) NOT NULL
        );
        SELECT create_hypertable('che_infos', 'time', if_not_exists => TRUE);
        
        CREATE MATERIALIZED VIEW che_infos_summary_daily
                    (day, name, count)
                    WITH (timescaledb.continuous)
        AS
        SELECT time_bucket('1 day', time), name, count()
        FROM che_infos
        GROUP BY time_bucket('1 day', time), name;
        SELECT add_continuous_aggregate_policy('che_infos_summary_daily', '12 hour', '1 hour', '1 hour');
        
        SELECT add_retention_policy('che_infos', INTERVAL '1 day');
        SELECT add_retention_policy('che_infos_summary_daily', INTERVAL '30 days');
</pre>
      </td>
   </tr>
   <tr>
      <td>
         api_requests
      </td>
      <td>
         Cette métrique donne pour chaque endpoint d'API les temps d'exécution
      </td>
      <td>
         <pre>
        CREATE TABLE IF NOT EXISTS api_requests
        (
            time     TIMESTAMP    NOT NULL,
            method   VARCHAR(7)   NOT NULL,
            url      VARCHAR(255) NOT NULL,
            interval INT          NOT NULL
        );
        SELECT create_hypertable('api_requests', 'time', if_not_exists => TRUE);
        SELECT add_retention_policy('api_requests', INTERVAL '30 days');
</pre>
      </td>
   </tr>
</table>
