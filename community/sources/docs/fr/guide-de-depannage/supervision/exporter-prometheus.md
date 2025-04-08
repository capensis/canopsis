# Exporter Prometheus pour Canopsis

Un **exporter Prometheus** est disponible dans Canopsis pour exposer des métriques internes au format compatible Prometheus. Ce composant permet d’intégrer facilement Canopsis dans votre système de supervision existant.

## Présentation

- **Chemin d’export** : `/metrics`
- **Port par défaut** : `9180` (modifiable via le flag `-port`)
- **Dépôt** : `go-engines-community/cmd/prometheus-exporter`

Cet exporter expose un ensemble de métriques utiles pour suivre l’état de la plateforme et le bon fonctionnement de ses composants internes.

## Métriques exposées

| Nom de la métrique                           | Description |
|----------------------------------------------|-------------|
| `canopsis_enrichment_errors`                 | Nombre d’erreurs d’enrichissement |
| `canopsis_active_alarms`                     | Nombre d’alarmes actives |
| `canopsis_closed_alarms`                     | Nombre d’alarmes clôturées |
| `canopsis_active_entities`                   | Nombre d’entités actives |
| `canopsis_disabled_entities`                 | Nombre d’entités désactivées |
| `canopsis_connected_users`                   | Nombre d’utilisateurs connectés |
| `canopsis_active_users`                      | Nombre d’utilisateurs actifs |
| `canopsis_event_filters`                     | Nombre de filtres d’événements |
| `canopsis_active_pbehavior`                  | Nombre de comportements prédictifs actifs |
| `canopsis_meta_alarms_rules`                 | Nombre de règles de méta-alarmes |
| `canopsis_dynamic_infos_rules`               | Nombre de règles d’informations dynamiques |
| `canopsis_engine_status{engine_name=}`       | Statut des moteurs (1 = actif, 0 = inactif) |
| `canopsis_last_exploitation_mod_time{type=}` | Date de dernière modification des éléments d’exploitation |

## Démarrage de l’exporter

L’exporter peut être lancé avec les options suivantes :

| Option                    | Valeur par défaut | Description |
|---------------------------|-------------------|-------------|
| `-version`                | `false`           | Affiche la version et quitte |
| `-port`                   | `9180`            | Port d’écoute du serveur |
| `-d`                      | `false`           | Active les logs de debug |
| `-updateMetricsInterval`  | `10s`             | Fréquence de mise à jour des métriques |

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
