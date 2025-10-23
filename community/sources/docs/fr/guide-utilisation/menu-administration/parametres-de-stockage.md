# Paramètres de stockage

Certaines données accumulées dans Canopsis peuvent être régulées par une politique de stockage.  

!!! Note
    Cette politique de stockage est appliquée une fois par semaine. Vous pouvez définir le jour et l'heure d'exécution dans le fichier de configuration `canopsis.toml`
    ```ini
    [Canopsis.data_storage]
    TimeToExecute = "Sunday,23"
    ```

Le tableau suivant décrit précisément chaque opération.  
La durée et l'unité sont configurables dans l'interface.


| Section | Opération | Action réalisée |
|---------|-----------|-----------------|
| [**Alarmes**](#alarmes)  | Archiver les alarmes résolues | Déplace les alarmes résolues dont `v.resolved` est plus ancien que l'intervalle défini de la collection `resolved_alarms` (TimeToKeepResolvedAlarms) vers la collection `archived_alarms`. |
|  | Supprimer les alarmes résolues | Supprime les alarmes archivées trop anciennes (`v.resolved`) de la collection `archived_alarms`. |
| [**Entités non liées**](#entites) | Archiver les entités sans événements | Déplace les entités "abandonnées" (avec `last_event_date` trop ancien ou `null`) de la collection `default_entities` vers la collection `archived_entities`. |
| [**Consignes**](#consignes) | Supprimer l'historique des consignes | Supprime les anciennes exécutions dans `instruction_execution` et les enregistrements associés de `job_history`. |
|  | Supprimer les statistiques de consignes | Définit une rétention sur la table TimescaleDB `instruction_execution_hourly` afin de supprimer automatiquement les lignes trop anciennes. |
|  | Supprimer le résumé des consignes | Définit une rétention sur la table TimescaleDB `instruction_execution_by_modified_on` afin de supprimer automatiquement les lignes trop anciennes. |
| [**Comportements périodiques**](#comportements-periodiques) | Supprimer les Comportements périodiques inactifs | Supprime les comportements de la collection `pbehavior` : ceux sans `rrule` et avec un `tstop` trop ancien, ou ceux avec `rrule_end` trop ancien. |
| [**JUnit**](#junit) | Supprimer les jeux de tests | Supprime les suites de tests de la collection `junit_test_suite` trop anciennes, efface les médias associés, et nettoie les IDs obsolètes dans les paramètres des widgets JUnit. |
| [**Healthcheck**](#healthcheck) | Supprimer le flux FIFO entrant | Définit une rétention sur la table TimescaleDB `message_rate_hourly` pour supprimer les lignes trop anciennes. |
| [**Webhooks**](#webhooks) | Supprimer l'historique des webhooks | Supprime les enregistrements trop anciens de la collection `webhook_history`. |
|  | Afficher les données d'authentification | Si désactivé, masque dans les logs tous les champs sensibles (tokens, mots de passe) en les remplaçant par `***`. |
| [**Métriques internes**](#metriques-internes) | Supprimer les métriques | Définit une rétention sur les tables TimescaleDB utilisées pour les métriques internes (KPI, performances moteurs). |
| [**Métriques externes**](#metriques-externes) | Supprimer les métriques | Définit une rétention sur la table TimescaleDB `perf_data` (et ses agrégats associés) pour supprimer automatiquement les données trop anciennes. |
| [**Messages d'erreur des filtres d'événements**](#filtres-devenements) | Supprimer les messages d'erreur | Supprime les enregistrements anciens de la collection `eventfilter_failure`. |
| [**Tags externes des alarmes**](#tags-externes) | Supprimer les tags externes | Supprime les tags de la collection `alarm_tag` dont `last_event_date` est trop ancien, puis efface les couleurs associées dans `alarm_tag_color`. |
| [**Enregistrements d'événements**](#enregistrements-devenements) | Supprimer les enregistrements | Supprime les enregistrements anciens de la collection `event_records`. |


## Alarmes

Le cycle de vie d'une alarme respecte le schéma suivant : 

```mermaid
flowchart TB
    A[🔔 Alarme ouverte] --> B[✅ Alarme résolue]

    B --> D{  }
    D -->|Alarme rouverte| A2[🔔 Alarme rouverte]
    D -->|Sinon| E[⏱️ Les alarmes sont archivées<br/>après ce délai]

    E --> F[📦 Alarmes archivées]
    F --> G[⏱️ Les alarmes sont supprimées<br/>après ce délai]
    G --> H[🗑️ Alarmes supprimées]

    %% Styles
    style A stroke:#f33,stroke-width:3px,color:#f33
    style A2 stroke:#f33,stroke-width:3px,color:#f33
    style B stroke:#3c3,stroke-width:3px,color:#3c3
    style F stroke:#888,stroke-width:2px,color:#555
    style H stroke:#888,stroke-width:2px,color:#555
```

L'archivage des alarmes consiste à déplacer les alarmes éligibles (résolues et respectant le délai défini) dans une collection de données dédiée.  
Ces alarmes restent ainsi disponibles pour les administrateurs en cas de besoin.

La suppression des alarmes résolues est quant à elle définitive et a lieu après le délai défini.

Par ailleurs, les alarmes `ouvertes` (collection `periodical_alarm`) et les alarmes `résolues` (collection `resolved`) ne sont désormais plus stockées dans le même espace pour garantir la performance d'accès aux alarmes en cours.  

Le paramètre `TimeToKeepResolvedAlarms` permet de définir le délai à partir duquel une alarme résolue passera de la collection `ouvertes` à la collection `résolues`

Ce paramètre se situe dans le fichier de configuration `canopsis.toml`.

```ini
[Canopsis.alarm]
# TimeToKeepResolvedAlarms defines how long resolved alarms will be kept in main alarm collection
TimeToKeepResolvedAlarms = "720h"
```

## Entités

Les entités désactivées peuvent être :

* archivées : déplacées dans une collection de données dédiée
* supprimées : supprimées définitivement de la collection d'archives

### Archiver les entités désactivées

!!! Attention 
    Une option permet également l'archivage ou la suppression des impacts et dépendances de ces entités.  
    Pour les **connecteurs**, tous les composants et ressources dépendants sont archivés ou supprimés.  
    Pour les **composants**, toutes les ressources dépendantes sont archivées ou supprimées.

!!! Note
    Cette opération n'est pas éligible à l'ordonnancement général et ne peut s'effectuer qu'à la demande 

### Archiver les entités non liées

## Consignes

Les statistiques d'exécution des remédiations sont conservées pendant le délai défini.
Passé ce délai :

* elles sont agrégées par semaine : seul le nombre total d'exécutions hebdomadaires est conservé ;
* puis elles sont supprimées définitivement lorsque la durée de rétention configurée est atteinte.

## Comportements périodiques

Les comportements périodiques sont soumis à une politique de suppression automatique.  
Un comportement peut être supprimé uniquement s'il respecte les conditions suivantes :

* il est inactif,
* il ne possède aucune période planifiée à venir.

Le délai de rétention configuré commence à courir à partir de la fin de la dernière période du comportement.
À l'issue de ce délai, le comportement est définitivement supprimé de la base de données.

## JUnit

Les données associées aux scénarios de tests JUnit (rapports XML, captures d'écran, vidéos, etc.) sont conservées pendant la durée définie.  
Au-delà de ce délai, elles sont supprimées définitivement de la base de données ainsi que des espaces de stockage associés.

## Healthcheck

Les mesures de santé collectées par Canopsis concernant le nombre d'événements entrants (moteur fifo) sont conservées pendant le délai configuré.  
Une fois ce délai atteint, ces données sont automatiquement supprimées afin de limiter l'espace occupé dans la base de données TimescaleDB (message_rate_hourly).

## Webhooks

Les historiques liés aux exécutions de webhooks sont stockés dans la collection `webhook_history`.  
Ils sont conservés pendant la durée configurée.  
Une fois ce délai atteint, ces enregistrements sont automatiquement supprimés.

Il est également possible de choisir si les données d'authentification (tokens, mots de passe, etc.) doivent apparaître en clair ou être masquées (\*\*\*) dans les journaux.

## Métriques internes

Les métriques internes générées par Canopsis (KPI, performances des moteurs, temps de traitement, etc.) sont conservées dans des tables TimescaleDB spécifiques.  
Passé le délai configuré, ces données sont automatiquement supprimées pour éviter une croissance trop importante de la base de données.

## Métriques externes

Les métriques externes issues des événements ([perf_data]((guide-developpement/structures/#detail-par-type-devenement)) sont stockées dans la table TimescaleDB `perf_data`.  
Elles sont conservées pendant la durée définie.  
Au-delà de ce délai, elles sont automatiquement supprimées, y compris leurs éventuels agrégats.

## Filtres d'événements

Lorsque des filtres d'événements génèrent des [erreurs](guide-utilisation/menu-exploitation/filtres-evenements/#gestion-des-erreurs), celles-ci sont stockées dans la collection `eventfilter_failure`.  
Ces messages d'erreur sont conservés pendant la durée configurée.
Une fois ce délai atteint, ils sont automatiquement supprimés.

## Tags externes

Les [tags](guide-utilisation/menu-administration/gestion-des-tags/#tags-presents-dans-les-evenements) créés à partir d'événements peuvent être supprimés passé un délai.

Les [tags](guide-utilisation/menu-administration/gestion-des-tags/#tags-presents-dans-les-evenements) créés à partir d'événements sont conservés dans la collection `alarm_tag`.  
Si leur last_event_date est plus ancien que le délai configuré, ces tags sont supprimés, et leurs couleurs associées sont également nettoyées de la collection `alarm_tag_color`.

## Enregistrements d'événements 

Les [enregistrements d'événements](guide-utilisation/menu-administration/enregistrements-d-evenements/) permettent de rejouer ou analyser a posteriori certains flux d'événements.  
Ils sont conservés jusqu'à ce que la durée configurée soit atteinte, puis ils sont automatiquement supprimés de la collection `event_records`.  
