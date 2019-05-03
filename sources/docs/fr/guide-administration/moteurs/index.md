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
  -autoRecomputeWatchers
        Recalcule automatiquement l'état des watchers chaque minute.
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

Le flag `autoRecomputeWatchers` permet de s'assurer que l'état des watchers est mis à jour à chaque battement du moteur axe.  

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

## Changer le niveau de log d'un moteur

### Avec systemd

Consulter l'état actuel du service :
```shell
systemctl status canopsis-engine@dynamic-alerts.service
```
Il contient les informations suivantes :
```shell
/opt/canopsis/bin/python2 /opt/canopsis/bin/engine-launcher -e canopsis.engines.dynamic -n alerts -w dynamic-alerts -l info
```
Dans cet exemple le niveau de logs du moteur alerts est "info".

Visualiser ensuite la configuration du service avec la commande suivante :
```shell
cat /lib/systemd/system/canopsis-engine@.service
```
Le résultat ressemble à ceci :
```
[Unit]
Description=Canopsis Engine %i
After=network.target
Documentation=https://doc.canopsis.net

[Service]
User=canopsis
Group=canopsis
WorkingDirectory=/opt/canopsis
Environment=VIRTUAL_ENV=/opt/canopsis
Environment=HOME=/opt/canopsis
Environment=PATH=$VIRTUAL_ENV/bin:/bin/:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin
Environment=LOGLEVEL=info
ExecStart=/opt/canopsis/bin/engine-launcher-systemd %i
PIDFile=/var/run/canopsis-engine-%i.pid
Restart=on-failure
Type=simple

[Install]
WantedBy=multi-user.target
```
La ligne qui défini le niveau de log du moteur au démarrage est `Environment=LOGLEVEL=info`.

Avant de modifier la configuration du moteur, copier le fichier de configuration dans le dossier `/etc/systemd/system`.
```shell
cp /lib/systemd/system/canopsis-engine@.service /etc/systemd/system/canopsis-engine@.service
```
Surcharger ensuite la configuration en modifiant la ligne comme ceci `Environment=LOGLEVEL=debug`.

Enfin recharger systemd pour s'assurer que la nouvelle configuration est bien prise en compte.
```shell
systemctl daemon-reload
```
Redémarrer le moteur.
```shell
systemctl stop canopsis-engine@dynamic-alerts.service
systemctl start canopsis-engine@dynamic-alerts.service
```
Et vérifier son nouveau statut :
```shell
systemctl status canopsis-engine@dynamic-alerts.service
```
```shell
/opt/canopsis/bin/python2 /opt/canopsis/bin/engine-launcher -e canopsis.engines.dynamic -n alerts -w dynamic-alerts -l debug
```

#### Méthode alternative (si les manipulations précédentes n'ont pas fonctionné)

Afficher la configuration du moteur avec la commande suivante :
```shell
systemctl cat canopsis-engine@dynamic-alerts.service
```
Puis éditer cette configuration dans un autre terminal :
```shell
systemctl edit canopsis-engine@dynamic-alerts.service
```
Une fenêtre d'édition vide s'affiche, copiez les éléments suivants pour surcharger la configuration actuelle :
```
[Service]
Environment=LOGLEVEL=debug
```
Sauvegarder et quitter l'éditeur.

Terminer en rechargeant systemd pour prendre en compte la modification.
```shell
systemctl daemon-reload
```

# Moteurs obsolètes

*  acknowledgement
*  cancel
*  context
*  eventstore
*  task\_linklist : n'existe plus depuis Canopsis 3.0
*  linklist : n'existe plus depuis Canopsis 3.0, remplacé par les linkbuilders
*  perfdata : n'existe plus depuis Canopsis 3.0, remplacé par metric
