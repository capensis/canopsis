# Exporter Prometheus pour Canopsis

Des **exporters Prometheus** sont disponibles dans Canopsis pour exposer des métriques internes au format compatible Prometheus. Canopsis fournit **deux mécanismes d’export** :

  - un exporter principal, déployé comme service dédié ;
  - un second exporter directement **intégré à l’engine FIFO**, permettant d’exposer les métriques relatives au traitement des événements.

Ces composants facilitent l’intégration complète de Canopsis dans votre système de supervision existant.

## Présentation

### Exporter Canopsis

- **Chemin d’export** : `/metrics`
- **Port par défaut** : `9180` (modifiable via le flag `-port`)
- **Dépôt** : `go-engines-community/cmd/prometheus-exporter`

Cet exporter expose un ensemble de métriques utiles pour suivre l’état de la plateforme et le bon fonctionnement de ses composants internes.

#### Métriques exposées

| Nom de la métrique                                  | Description |
|-----------------------------------------------------|-------------|
| `canopsis_eventfilter_errors`                       | Nombre d’erreurs dans les filtres d’événements |
| `canopsis_opened_alarms`                            | Nombre d’alarmes ouvertes |
| `canopsis_resolved_alarms`                          | Nombre d’alarmes résolues |
| `canopsis_active_entities`                          | Nombre d’entités actives |
| `canopsis_disabled_entities`                        | Nombre d’entités désactivées |
| `canopsis_user_connections`                         | Nombre de connexions utilisateur |
| `canopsis_enabled_users`                            | Nombre d’utilisateurs activés |
| `canopsis_event_filters`                            | Nombre total de filtres d’événements |
| `canopsis_active_pbehavior`                         | Nombre de comportements périodiques actifs |
| `canopsis_meta_alarms_rules`                        | Nombre de règles de méta-alarmes |
| `canopsis_dynamic_infos_rules`                      | Nombre de règles d’informations dynamiques |
| `canopsis_engine_status{engine_name="<nom>"}`       | Statut d’un moteur Canopsis (1 = actif, 0 = inactif) |
| `canopsis_last_exploitation_mod_time{type="<type>"}`| Horodatage de dernière modification d’un élément d’exploitation |

#### Démarrage de l’exporter

L’exporter peut être lancé avec les options suivantes :

| Option                    | Valeur par défaut | Description |
|---------------------------|-------------------|-------------|
| `-version`                | `false`           | Affiche la version et quitte |
| `-port`                   | `9180`            | Port d’écoute du serveur |
| `-d`                      | `false`           | Active les logs de debug |
| `-updateMetricsInterval`  | `10s`             | Fréquence de mise à jour des métriques |


#### Variables d’environnement

L’exporter nécessite que les variables d’environnement suivantes soient définies à l’exécution :

- `CPS_REDIS_URL` : URI de connexion à Redis.
- `CPS_MONGO_URL` : URI de connexion à MongoDB.

#### Préférence de lecture MongoDB

Afin de réduire la charge sur le **noeud primaire MongoDB**, il est recommandé d’utiliser l’option `readPreference=secondary` dans l’URI de connexion. Cela permet à l’exporter de lire les données à partir des **noeuds secondaires** du replica set, sans impacter les opérations d’écriture sur le primaire.

**Exemple d’URI :**

```
mongodb://cpsmongo:canopsis@localhost:27017,localhost:27018,localhost:27019/canopsis?replicaSet=cps&readPreference=secondary
```

=== "Paquets RPM"

    Suite à l’installation de Canopsis via les packages RPM, le binaire de l’exporter Prometheus est disponible dans `/opt/canopsis/bin/prometheus-exporter`. 

    Vous pouvez lancer l’exporter directement en adaptant les variables d’environnement à votre contexte (MongoDB, Redis) :

    ```sh
    CPS_MONGO_URL=mongodb://cpsmongo:canopsis@localhost:27017/canopsis?replicaSet=rs0&readPreference=secondary \
    CPS_REDIS_URL=redis://:canopsis@localhost:6379/0 \
    /opt/canopsis/bin/prometheus-exporter
    ```

    Pour intégrer l’exporter au démarrage et au système de supervision local, créez un service `systemd`. 
    
    !!! warning "Adapter variables d'environnements"
        Veillez à adapter les variables d'environnements `CPS_MONGO_URL` et `CPS_REDIS_URL` à votre contexte

    Exemple d’unité à placer dans `/usr/lib/systemd/system/canopsis-engine-go@prometheus-exporter.service` :

    ```sh
    cat <<EOF > /usr/lib/systemd/system/canopsis-engine-go@prometheus-exporter.service
    [Unit]
    Description=Canopsis Go Prometheus Exporter
    PartOf=canopsis.service
    After=canopsis.service
    ReloadPropagatedFrom=canopsis.service
    After=network.target
    Documentation=https://doc.canopsis.net/guide-de-depannage/supervision/exporter-prometheus

    [Service]
    User=canopsis
    Group=canopsis
    WorkingDirectory=/opt/canopsis
    Environment=HOME=/opt/canopsis
    ExecStart=/usr/bin/env /opt/canopsis/bin/prometheus-exporter -port 9180
    Restart=always
    RestartSec=1
    StartLimitBurst=300
    Environment="CPS_MONGO_URL=mongodb://cpsmongo:canopsis@localhost:27017/canopsis?replicaSet=rs0&readPreference=secondary"
    Environment="CPS_REDIS_URL=redis://:canopsis@localhost:6379/0"

    [Install]
    WantedBy=canopsis.service
    EOF
    ```

    Puis rechargez systemd et activez le service :

    ```sh
    systemctl reload daemon
    systemctl enable --now canopsis-engine-go@prometheus-exporter.service
    ```

=== "Docker"

    Un service spécifique est distribué dans l'environnement de référence. Il est attaché au profil "prometheus".  
    Pour démarrer l'exporter, veuillez exécuter la commande suivante :  

    ```sh
    CPS_EDITION=pro docker compose --profile=prometheus up -d prometheus-exporter
    ```

    Les logs associés sont alors :

    ```sh
    CPS_EDITION=pro docker compose --profile=prometheus logs  prometheus-exporter
    prometheus-exporter-1  | 2025-04-08T14:10:00Z INF git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo/mongo.go:532 > replica set is detected, transactions are enabled
    prometheus-exporter-1  | 2025-04-08T14:10:00Z INF git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/cmd/prometheus-exporter/main.go:140 > prometheus exporter started
    ```

=== "Helm"

    L’exporter Prometheus peut être activé directement via les valeurs du chart Canopsis.

    La fonctionnalité est exposée dans la section `metrics`, qu’il est possible de surcharger dans votre fichier `customer-values.yaml`.

    Pour activer l’exporter, ajoutez la configuration suivante :

    ```yaml
    metrics:
      enabled: true        # Active le service exporter
      port: 9180           # Port d’exposition du endpoint /metrics
      serviceMonitor:
        enabled: true      # Crée automatiquement un ServiceMonitor (Prometheus Operator)
    ```

    Ensuite, appliquez le déploiement via Helm :

    ```sh
    helm upgrade --install canopsis-prod canopsis/canopsis-pro \
    -f customer-values.yaml
    ```

    Une fois déployé, **un pod dédié sera créé**, correspondant à l’exporter Prometheus, ainsi qu’un **ServiceMonitor**, si vous utilisez le Prometheus Operator.

    D’autres paramètres du bloc `metrics` peuvent être surchargés selon vos besoins.
    
    La liste complète est disponible dans le **README du chart Canopsis**, incluant notamment la configuration du service, des annotations et des labels supplémentaires pour l’intégration dans votre environnement Prometheus.

### Exporter FIFO

- **Chemin d’export** : `/metrics`
- **Port par défaut** : `9180` (modifiable via le flag `-prometheusExporterPort`)
- **Dépôt** : `go-engines-community/cmd/engine-fifo`

L’exporter FIFO est directement intégré au moteur engine-fifo.

Il expose ses métriques internes sur un endpoint Prometheus lorsqu’il est activé via les options dédiées.

#### Métriques exposées

| Nom de la métrique                                         | Description                                                    |
| ---------------------------------------------------------- | -------------------------------------------------------------- |
| `canopsis_events_counter{type="ack"}`                      | Nombre d’événements d’acquittement (ack).                      |
| `canopsis_events_counter{type="ackremove"}`                | Nombre d’annulations d’acquittement.                           |
| `canopsis_events_counter{type="assocticket"}`              | Nombre d’événements d’association d’un ticket.                 |
| `canopsis_events_counter{type="autoinstructionactivate"}`  | Nombre d’auto-instructions activées.                           |
| `canopsis_events_counter{type="cancel"}`                   | Nombre d’annulations d’événements ou d’actions.                |
| `canopsis_events_counter{type="changestate"}`              | Nombre de changements d’état sur les entités ou alarmes.       |
| `canopsis_events_counter{type="check"}`                    | Nombre d’événements de type “check” (vérifications).           |
| `canopsis_events_counter{type="comment"}`                  | Nombre de commentaires ajoutés.                                |
| `canopsis_events_counter{type="contextupdate"}`            | Nombre de mises à jour de contexte sur les entités/événements. |
| `canopsis_events_counter{type="declareticketwebhook"}`     | Nombre de déclenchements du webhook de création de ticket.     |
| `canopsis_events_counter{type="entitytoggled"}`            | Nombre d’activation/désactivation d’entités.                   |
| `canopsis_events_counter{type="entityupdated"}`            | Nombre de mises à jour d’entités.                              |
| `canopsis_events_counter{type="instructionaborted"}`       | Nombre d’instructions abandonnées.                             |
| `canopsis_events_counter{type="instructioncompleted"}`     | Nombre d’instructions terminées avec succès.                   |
| `canopsis_events_counter{type="instructionfailed"}`        | Nombre d’instructions échouées.                                |
| `canopsis_events_counter{type="instructionpaused"}`        | Nombre d’instructions mises en pause.                          |
| `canopsis_events_counter{type="instructionresumed"}`       | Nombre d’instructions reprises.                                |
| `canopsis_events_counter{type="instructionstarted"}`       | Nombre d’instructions démarrées.                               |
| `canopsis_events_counter{type="junittestcaseeupdated"}`    | Nombre de mises à jour de cas de test JUnit.                   |
| `canopsis_events_counter{type="junittestsuiteupdated"}`    | Nombre de mises à jour de suites de tests JUnit.               |
| `canopsis_events_counter{type="metaalarm"}`                | Nombre d’événements liés aux méta-alarmes.                     |
| `canopsis_events_counter{type="metaalarmattachchildren"}`  | Nombre d’ajouts d’enfants à une méta-alarme.                   |
| `canopsis_events_counter{type="metaalarmchildactivate"}`   | Nombre d’activations d’un enfant dans une méta-alarme.         |
| `canopsis_events_counter{type="metaalarmchilddeactivate"}` | Nombre de désactivations d’un enfant dans une méta-alarme.     |
| `canopsis_events_counter{type="metaalarmdetachchildren"}`  | Nombre de détachements d’enfants d’une méta-alarme.            |
| `canopsis_events_counter{type="noevents"}`                 | Nombre d’occurrences sans événement associé.                   |
| `canopsis_events_counter{type="pbhenter"}`                 | Nombre d’entrées dans un pbehavior.                            |
| `canopsis_events_counter{type="pbhleave"}`                 | Nombre de sorties d’un pbehavior.                              |
| `canopsis_events_counter{type="pbhleaveandenter"}`         | Nombre de transitions sortie/entrée de pbehavior.              |
| `canopsis_events_counter{type="recomputeentityservice"}`   | Nombre de recalculs de service pour une entité.                |
| `canopsis_events_counter{type="resolve_cancel"}`           | Nombre d’événements de résolution par annulation.              |
| `canopsis_events_counter{type="resolve_close"}`            | Nombre d’événements de fermeture d’alarme.                     |
| `canopsis_events_counter{type="resolve_deleted"}`          | Nombre d’événements liés à la suppression d’alarme.            |
| `canopsis_events_counter{type="run_delayed_scenario"}`     | Nombre de scénarios différés exécutés.                         |
| `canopsis_events_counter{type="snooze"}`                   | Nombre d’événements de mise en pause (snooze).                 |
| `canopsis_events_counter{type="trigger"}`                  | Nombre de déclenchements d’événements.                         |
| `canopsis_events_counter{type="uncancel"}`                 | Nombre d’annulations de “cancel” (réactivation).               |
| `canopsis_events_counter{type="unsnooze"}`                 | Nombre de sorties de snooze.                                   |
| `canopsis_events_counter{type="updatestatus"}`             | Nombre de mises à jour de statut (status update).              |

#### Démarrage de l’exporter

L’exporter fait partie de l'engine FIFO, il peut être lancé avec les options suivantes :

| Option                     | Valeur par défaut | Description                      |
|----------------------------|-------------------|----------------------------------|
| `-enablePrometheusExporter`| `false`           | Activer ou désactiver l'exporter |
| `-prometheusExporterPort`  | `9180`            | Port d’écoute du serveur         |

=== "Paquets RPM"

    À l’installation de Canopsis via les **packages RPM**, l’exporter FIFO n’est pas un binaire séparé : il est directement inclus dans `/opt/canopsis/bin/engine-fifo`.

    Pour activer l’exporter, il faut **surcharger le service systemd de l’engine FIFO**.

    Exécuter :

    ```sh
    systemctl edit canopsis-engine-go@engine-fifo.service
    ```

    Ajouter :

    ```
    [Service]
    ExecStart=
    ExecStart=/usr/bin/env /opt/canopsis/bin/engine-fifo -enablePrometheusExporter=true -prometheusExporterPort=9180
    ```

    Recharger systemd :

    ```sh
    systemctl daemon-reload
    ```

    Redémarrer l'engine FIFO :

    ```sh
    systemctl restart canopsis-engine-go@engine-fifo.service
    ```

=== "Docker"

    Dans un déploiement Docker, l’exporter FIFO doit être activé en surchargeant le service `fifo`.
    La méthode recommandée consiste à utiliser un fichier :

    `docker-compose.override.yml`

    Cela permet d’ajouter ou modifier uniquement les paramètres nécessaires sans toucher au `docker-compose.yml` d’origine.

    Editer le fichier docker-compose.override.yaml et y ajouter :

    ```yaml
    fifo:
      command: /engine-fifo -enablePrometheusExporter=true -prometheusExporterPort=9180
      # Si vous souhaitez exposer le port sur l'hôte
      ports:
      - "9180:9180"
    ```

    Redémarrer uniquement le service FIFO :

    ```sh
    docker compose up -d fifo
    ```

=== "Helm"

    Avec Helm, l’exporter FIFO se configure directement dans les valeurs du chart Canopsis.

    Ajouter à votre fichier customer-values.yaml :

    ```yaml
    fifo:
      metrics:
        enabled: true
        serviceMonitor:
          enabled: true
      port: 9180
    ```

    Ensuite, appliquer le déploiement via Helm :

    ```sh
    helm upgrade --install canopsis-prod canopsis/canopsis-pro \
    -f customer-values.yaml
    ```

    Une fois appliqué :

      * l’engine FIFO expose `/metrics` sur le port 9180
      * un **ServiceMonitor** est créé si le Prometheus Operator est présent

## Configuration de Prometheus

Pour intégrer ces exporters à Prometheus, ajoutez la configuration suivante dans votre fichier `prometheus.yml` :

```yaml
scrape_configs:
  - job_name: 'canopsis_exporter'
    scrape_interval: 15s
    static_configs:
      - targets:
            - 'your-exporter-host:9180'
            - 'fifo-host:9180'
```

!!! tip "Conseil"
    Veillez à ce que la valeur de scrape_interval dans Prometheus soit supérieure ou égale à -updateMetricsInterval côté exporter. Sinon, Prometheus risque de collecter des valeurs identiques ou obsolètes.

## Dashboard Grafana

![dashboard grafana](./img/dashboard-grafana.png)

Un dashboard Grafana clé en main, permettant de visualiser l’ensemble des métriques exposées par les exporters de Canopsis, est disponible dans notre dépôt GitLab interne.

Si vous souhaitez y accéder ou l’importer dans votre instance Grafana, contactez notre équipe afin d’obtenir les accès au repository.
