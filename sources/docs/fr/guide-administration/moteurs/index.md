# Sommaire et présentation

- [Activation et désactivation des moteurs](activation-desactivation-moteurs.md)
- [Enchainement des moteurs](schema-enchainement-moteurs.md)

Les évènements envoyés par des connecteurs à Canopsis sont traités à l'aide de moteurs.

Un moteur a **plusieurs rôles** :

*  consommation d'un évènement : pour le traiter, puis l'acheminer vers le(s) moteur(s) suivant(s).
*  effectuer une tâche périodique : appelée « beat », cette tâche sera exécutée à intervalle régulier.
*  consommation d'un enregistrement lorsque les enregistrements de la base de données sont disponibles.

Un moteur peut avoir les **propriétés** suivantes :

*  un type (the python module to load)
*  un nom (must be unique)
*  un identifiant (0, 1, 2, 3, ..., must be unique)
*  un niveau de log (debug, info, warning, or error)

Le listing des moteurs peut être réalisé grâce à cette commande : `systemctl list-units "canopsis*"`

## Liste des moteurs

### Moteurs Go

| Moteur         | Description                                                                      | CAT ?              |
|:---------------|:---------------------------------------------------------------------------------|:------------------:|
| [action](engine-action.md)          | Applique des actions définies par l'utilisateur.                                 |                    |
| axe            | Gère le cycle de vie des alarmes.                                                |                    |
| [axe@**webhooks**](../webhooks/index.md)   | Gère le système de webhooks vers des services externes.                                                | ✅                 |
| che            | Supprime les évènements invalides, gère le contexte, et enrichit les évènements. |                    |
| heartbeat      | Surveille des entités, et lève des alarmes en cas d'absence d'information.       |                    |
| stat           | Calcule des statistiques sur les états des alarmes.                              |                    |

### Moteurs Python

| Moteur                                                         | Description                                              | CAT ?              |
|:---------------------------------------------------------------|:---------------------------------------------------------|:------------------:|
| canopsis-engine@**dynamic-alerts**.service                     | Gère le cycle de vie des alarmes.                        |                    |
| canopsis-engine@**cleaner-cleaner_alerts**.service             | Supprime les évènements invalides.                       |                    |
| canopsis-engine@**cleaner-cleaner_events**.service             | Supprime les évènements invalides.                       |                    |
| canopsis-engine@**dynamic-context-graph**.service              | Stocke les données contextuelles des évènements.         |                    |
| **datametrie**                                                 | Gère le connecteur datametrie.                           | ✅             |
| [canopsis-engine@**event_filter-event_filter**.service](moteur-event_filter.md)          | Applique des règles de filtrage.                         |                    |
| canopsis-engine@**metric-metric**.service                      | Stocke les données de métrologie des évènements.         |                    |
| canopsis-engine@**dynamic-pbehavior**.service                  | Gère les périodes de maintenance.                        |                    |
| canopsis-engine@**scheduler-scheduler**.service                | Envoyer un travail à des gestionnaires de tâches.        |                    |
| [canopsis-engine-cat@**snmp**](moteur-snmp.md)                                                       | Gère les traps SNMP.                                     | ✅             |
| canopsis-engine@**task_dataclean-task_dataclean**.service      | Gestionnaire pour supprimer anciennes données.           |                    |
| canopsis-engine@**task_importctx-task_importctx**.service      | Gestionnaire des imports de données en masse.            |                    |
| [canopsis-engine-cat@**task_ackcentreon-task_ackcentreon**.service](moteur-task_ackcentreon.md)      | ACK descendants vers Centreon.            | ✅ |
| canopsis-engine@**task_mail-task_mail**.service                | Gestionnaire de tâches pour envoyer du courrier.         |                    |
| canopsis-engine@**ticket-ticket**.service                      | Gère les tickets externes.                               |                    |
| canopsis-engine@**dynamic-watcher**.service                    | Gère les watchers (groupes de surveillance).             |                    |
| canopsis-engine-cat@**statsng-statsng**.service                | Calcule des statistiques sur les alarmes et les entités. | ✅             |

## Flags & Usage

### Utilisation de engine-action

```
  -d    debug
  -version
        version infos
```

### Utilisation de engine-axe

```
  -autoDeclareTickets
        Déclare les tickets automatiquement pour chaque alarme. DÉPRÉCIÉ, remplacé par les webhooks.
  -d    debug
  -featureHideResources
        Active les features de gestion de ressources cachées.
  -featureStatEvents
        Envoie les évènements de statistiques
  -postProcessorsDirectory
        Le répetoire contenant les plugins de post-traitement (par défaut ".")
  -printEventOnError
        Afficher les évènements sur les erreurs de traitement.
  -publishQueue
        Publie les événements sur cette queue. (par défaut "Engine_action")
  -version
        version infos
```

### Utilisation de engine-che

```
  -consumeQueue string
        Consomme les évènements venant de cette queue. (default "Engine_che").
  -createContext
        Active la création de context graph. Activé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorse que vous l'utilisez. (default true)
  -d    debug
  -dataSourceDirectory
        The path of the directory containing the event filter's data source plugins. (default ".")
  -enrichContext
        Active l'enrichissment de context graph à partir d'un event. Désactivé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorse que vous l'utilisez. (default true)
  -enrichExclude string
        Liste de champs séparés par des virgules ne faisant pas partie de l'enrichissement du contexte
  -enrichInclude string
        Coma separated list of the only fields that will be part of context enrichment. If present, -enrichExclude is ignored.
  -printEventOnError
        Print event on processing error
  -processEvent
        enable event processing. enabled by default. (default true)
  -publishQueue
        Publie les événements sur cette queue. (default "Engine_event_filter")
  -purge
        purge consumer queue(s) before work
  -version
        version infos
```

### Utilisation de engine-heartbeat

```
  -d    debug
  -version
        version infos
```

### Utilisation de engine-stat

```
  -d    debug
  -version
        version infos
```

# Moteurs obsolètes

*  acknowledgement
*  cancel
*  context
*  eventstore
*  task\_linklist : n'existe plus depuis Canopsis 3.0
*  linklist : n'existe plus depuis Canopsis 3.0, remplacé par les linkbuilders
*  perfdata : n'existe plus depuis Canopsis 3.0, remplacé par metric
