# Exporter Prometheus pour Canopsis

Un **exporter Prometheus** est disponible dans Canopsis pour exposer des métriques internes au format compatible Prometheus. Ce composant permet d’intégrer facilement Canopsis dans votre système de supervision existant.

## Présentation

- **Chemin d’export** : `/metrics`
- **Port par défaut** : `9180` (modifiable via le flag `-port`)
- **Dépôt** : `go-engines-community/cmd/prometheus-exporter`

Cet exporter expose un ensemble de métriques utiles pour suivre l’état de la plateforme et le bon fonctionnement de ses composants internes.

## Métriques exposées

## Métriques exposées

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

## Démarrage de l’exporter

L’exporter peut être lancé avec les options suivantes :

| Option                    | Valeur par défaut | Description |
|---------------------------|-------------------|-------------|
| `-version`                | `false`           | Affiche la version et quitte |
| `-port`                   | `9180`            | Port d’écoute du serveur |
| `-d`                      | `false`           | Active les logs de debug |
| `-updateMetricsInterval`  | `10s`             | Fréquence de mise à jour des métriques |


## Variables d’environnement

L’exporter nécessite que les variables d’environnement suivantes soient définies à l’exécution :

- `CPS_REDIS_URL` : URI de connexion à Redis.
- `CPS_MONGO_URL` : URI de connexion à MongoDB.

### Préférence de lecture MongoDB

Afin de réduire la charge sur le **noeud primaire MongoDB**, il est recommandé d’utiliser l’option `readPreference=secondary` dans l’URI de connexion. Cela permet à l’exporter de lire les données à partir des **noeuds secondaires** du replica set, sans impacter les opérations d’écriture sur le primaire.

**Exemple d’URI :**

```
mongodb://cpsmongo:canopsis@localhost:27017,localhost:27018,localhost:27019/canopsis?replicaSet=cps&readPreference=secondary
```

=== "Paquets RPM"

    A venir

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

    A venir

## Configuration de Prometheus

Pour intégrer cet exporter à Prometheus, ajoutez la configuration suivante dans votre fichier `prometheus.yml` :

```yaml
scrape_configs:
  - job_name: 'canopsis_exporter'
    scrape_interval: 15s
    static_configs:
      - targets: ['your-exporter-host:9180']
```

!!! tip "Conseil"
    Veillez à ce que la valeur de scrape_interval dans Prometheus soit supérieure ou égale à -updateMetricsInterval côté exporter. Sinon, Prometheus risque de collecter des valeurs identiques ou obsolètes.





