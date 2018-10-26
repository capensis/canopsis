# Sommaire et présentation

- [L'arrêt et la relance des moteurs](activation-desactivation-moteurs.md)  
- [Schéma d'enchainement](schema-enchainement-moteurs.md)  

**TODO (DWU) :** doublons sur certains moteurs, entre cette partie du Guide Administrateur et le Guide Utilisateur ! Exemple : pbehaviors, météo… Ne pas documenter la même chose à 2 endroits !

**(MG) :** Moteur présent dans schéma d'enchaînements non présent dans liste des moteurs et inversement, normal ?
**(MG)/amélioration :** Mettre liens clicables dans tableau moteurs vers les docs.

**TODO (DWU) :** schéma, etc. Plus de détails sur les champs obligatoires aujourd'hui ?   
**TODO :** les différents flags (ex : printEventOnError, publishQueue) ne sont pas du tout documentés.

Les événements envoyés par des connecteurs à Canopsis sont traités à l'aide de moteurs.  
  
Un moteur a **plusieurs rôles**:  

- consommation d'un événement: pour le traiter, puis l'acheminer vers le(s) moteur(s) suivant(s).  
- Effectuer une tâche périodique: appelée «beat», cette tâche sera exécutée à intervalle régulier.  
- Consommation d'un enregistrement lorsque les enregistrements de la base de données sont disponibles.  

Chaque moteur est défini par un ensemble de procédures, utilisé pour effectuer les tâches énumérées ci-dessus.  

Un moteur peut avoir les **propiètées** suivantes :

- un type (the python module to load)  
- un nom (must be unique)  
- Un identifiant (0, 1, 2, 3, ..., must be unique)  
- Un niveau de log (debug, info, warning, or error)  

Le listing des moteur peut être réalisé grace à cette commande : `ps -aux | grep canopsis`

## Liste des moteurs

Ce document rassemble les informations sur les moteurs Canopsis.  

**TODO (DWU) :** noms complets :
* `canopsis-engine@dynamic-alerts.service`
* `canopsis-engine@cleaner-cleaner_alerts.service`
* `canopsis-engine@cleaner-cleaner_events.service`
* `canopsis-engine@dynamic-context-graph.service`
* `canopsis-engine@event_filter-event_filter.service`
* `canopsis-engine@linklist-linklist.service` **existe plus**
* `canopsis-engine@dynamic-pbehavior.service`
* `canopsis-engine@scheduler-scheduler.service`
* `canopsis-engine@selector-selector.service`
* `canopsis-engine@task_dataclean-task_dataclean.service`
* `canopsis-engine@task_importctx-task_importctx.service`
* `canopsis-engine@task_linklist-task_linklist.service` **existe plus**
* `canopsis-engine@task_mail-task_mail.service`
* `canopsis-engine@ticket-ticket.service`
* `canopsis-engine@dynamic-watcher.service`
* `canopsis-engine-cat@task_blabla.service`


### Moteurs NG

(Moteurs Go)

| Moteur         | Description                                                                                       | CAT ?              |
|:---------------|:--------------------------------------------------------------------------------------------------|:------------------:|
| axe            | TODO                                                                                              |                    |
| che            | TODO                                                                                              |                    |
| heartbeat      | Parcourir les entités surveillées pour surveiller l'absence d'information, clore les alarmes.     |                    |
| stat           | Peupler le registre d'activité (donnés et méta-données).                                          |                    |

**TODO (DWU) :** voir avec Lucas ce qui est prévu entre stat et statng ?

### Moteurs Python
 
| Moteur                                                         | Description                                      | CAT ?              |
|:---------------------------------------------------------------|:-------------------------------------------------|:------------------:|
| canopsis-engine@**dynamic-alerts**.service                     | Gérer les alarmes à partir d'événements          |                    | 
| canopsis-engine@**cleaner-cleaner_alerts**.service             | Supprimer des événements invalides               |                    | 
| canopsis-engine@**cleaner-cleaner_events**.service             | Supprime les événements incorrects               |                    |
| canopsis-engine@**dynamic-context-graph**.service              | Stocker les données contextuelles de l'événement |                    | 
| **datametrie**                                                 | Gestion du connecteur datametrie                 | :white_check_mark: | 
| canopsis-engine@**event_filter-event_filter**.service          | Appliquer des règles de filtrage                 |                    | 
| **metric**                                                     | Stocker les données de métrologie de l'événement |                    | 
| canopsis-engine@**dynamic-pbehavior**.service                  | Gestion des périodes de maintenance              |                    |
| canopsis-engine@**scheduler-scheduler**.service                | Envoyer un travail à des gestionnaires de tâches |                    | 
| **snmp**                                                       | Pour la gestion des traps SNMP                   | :white_check_mark: | 
| canopsis-engine@**task_dataclean-task_dataclean**.service      | Gestionnaire pour supprimer anciennes données    |                    | 
| canopsis-engine@**task_importctx-task_importctx**.service      | Gestionnaire des imports de données en masse     |                    | 
| canopsis-engine@**task_mail-task_mail**.service                | Gestionnaire de tâches pour envoyer du courrier  |                    | 
| canopsis-engine@**ticket-ticket**.service                      | Gestion du ticketing                             |                    | 
| canopsis-engine@**dynamic-watcher**.service                    | Gestion des Watchers (groupes de surveillance)   |                    | 

## Flags & Usage

### Utilisation de engine-axe

```
  -d    debug
  -featureHideResources
        Active les features de gestion de ressources cachées.
  -featureStatEvents
        Envoie les événements de statistiques
  -printEventOnError
        Afficher les événements sur les erreurs de traitement.
  -version
        version infos
```

### Utilisation de engine-che

```
revolution/cmd/engine-che/engine-che:
  -consumeQueue string
        Consomme les événements venant de cette queue. (default "Engine_che").
  -createContext
        Active la création de context graph. Activé par défaut.
        WARNING: désactiver l'ancien moteur context-graph lorse que vous l'utilisez. (default true)
  -d    debug
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
  -publishQueue string
        Publish event to this queue. (default "Engine_event_filter")
  -purge
        purge consumer queue(s) before work
  -version
        version infos
```

à terminer.

# Moteurs obsolètes

**TODO (DWU) :** je pense que ça ne doit apparaître que dans les guides de mise à jour (fichiers `UPGRAGING_` à la racine du dépôt canopsis/canopsis, sur `develop`), qui doivent être réécrits et migrés ailleurs dans cette doc.

*  acknowledgement
*  cancel
*  context
*  task_linklist : n'existe plus depuis Canopsis 3.0
*  linklist : n'existe plus depuis Canopsis 3.0, remplacé par les linkbuilders
*  perfdata : n'existe plus depuis Canopsis 3.0, remplacé par metric
