# Structure des événements

## Focus AMQP

* Vhost : `canopsis`
* Routing key : `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`
* Exchange : `canopsis.events`
* Options de l'Exchange : type : "topic", durable : true, auto_delete : false
* Content Type : `application/json`

## Structure basique d'un événement

Voici la structure de base d'un [événement](../../guide-utilisation/vocabulaire/index.md#evenement), commune à tous les [types d'événements](#liste-des-types-devenements).

```json
{
  "event_type": "",      // Type de l'événement - type `string`
  "source_type": "",     // Source de l'événement ("component" ou "resource") - type `string`
  "connector": "",       // Type de connecteur (gelf, nagios, snmp, ...) - type `string`
  "connector_name": "",  // Identifiant du connecteur - type `string`
  "component": "",       // Nom du composant - type `string`
  "resource": "",        // Nom de la ressource (si `source_type` = "resource") - type `string`

  // Champs optionnels
  "timestamp": 0,        // Timestamp UNIX de l'émission de l'événement - type `integer`
  "output": "",          // Message court - type `string`
  "long_output": ""      // Message détaillé - type `string`
}
```

!!! note "Remarques sur le champ `timestamp`"
    
    * Champ optionnel mais fortement recommandé.
    * Si absent, Canopsis prendra l'heure courante.
    * Depuis la version 4.5, seuls les événements dont le `timestamp` est dans la fenêtre **Now ± 24h** sont conservés.
    * Dans les alarmes créées :
      * `creation_date` = maintenant
      * `last_update_date` = valeur du `timestamp`


## Liste des types d'événements

| Type d'événement | Fonction | Origine | Interne ? |
|:---------------------|:---------|:--------|:----------|
| ack | Acquitte une alarme | UI / Engine-axe | Non |
| ackremove | Retire un acquittement | UI / Engine-axe | Non |
| assocticket | Associe un ticket | UI / Engine-axe | Non |
| cancel | Annule une alarme | UI / Engine-axe | Non |
| check | Crée ou met à jour une alarme | Connecteurs divers | Non |
| contextupdate | Met à jour l'entité sans changer la date du dernier événement | Connecteurs divers | Non |
| comment | Ajoute un commentaire | UI | Non |
| declareticket | Déclare un ticket sur l'alarme | UI | Non |
| done | Marque une alarme comme résolue | ??? | Non |
| changestate | Change et vérrouille la criticité d'une alarme | UI / Engine-axe | Non |
| snooze | Snooze une alarme (mise en pause) | UI / Engine-axe | Non |
| unsnooze | Retire le snooze | Engine-axe | Oui |
| uncancel | Annule une annulation | ??? | Non |
| metaalarm | Crée une méta-alarme | Engine-correlation | Oui |
| metaalarmupdated | Met à jour une méta-alarme | Engine-correlation | Oui |
| pbhenter / pbhleave / pbhleaveandenter | Gestion des comportements périodiques (pbh) | Engine-pbehavior | Oui |
| pbhcreate | Crée un comportement PBH | Engine-axe | Oui |
| resolve_done / resolve_cancel / resolve_close / resolve_deleted | Résolution d'alarmes | Moteurs divers | Oui |
| updatestatus | Mise à jour de statut | Engine-axe | Oui |
| manual_metaalarm_group / manual_metaalarm_ungroup / manual_metaalarm_update | Gestion manuelle des méta-alarmes | UI | Oui |
| activate | Activation d'une alarme | Engine-axe | Oui |
| run_delayed_scenario | Exécution différée d'un scénario | Engine-action | Oui |
| instructionstarted / instructionpaused / instructionresumed / instructioncompleted / instructionfailed / instructionaborted / autoinstructionstarted / autoinstructioncompleted / autoinstructionfailed / autoinstructionalreadyrunning | Statut d'exécution des remédiations | Engine-remediation / API | Oui |
| recomputeentityservice | Recalcul du graphe de service | API | Oui |
| updateentityservice | Mise à jour du cache service | API | Oui |
| entityupdated | Notification de modification d'entité | Engine-che / API | Oui |
| entitytoggled | Activation / Désactivation d'une entité | API | Oui |
| alarmskipped | Alarme ignorée pendant le recompute | Engine-service | Oui |
| junittestsuiteupdated / junittestcaseupdated | Mise à jour de tests Junit | connector-junit | Oui |
| noevents | Création d'une alarme par règle idle-rule | Engine-axe | Oui |

!!! warning "Attention"
    Les événements internes doivent être utilisés avec précaution voire non utilisés. Leur usage manuel peut entrainer un dysfonctionnement de Canopsis !


## Format public des événements

### Champs obligatoires

Un événement doit être un objet JSON comportant :

* `event_type`     : Type de l'événement (`string`)
* `source_type`    : Source de l'événement (`string`) * "component" ou "resource"
* `connector`      : Type de connecteur (`string`)
* `connector_name` : Nom du connecteur (`string`)
* `component`      : Composant concerné (`string`)
* `resource`       : Ressource (obligatoire si `source_type` = "resource") (`string`)
* `timestamp`      : Timestamp UNIX (`integer`) (optionnel mais recommandé)


### Détail par type d'événement

#### Evénement `check`

Champs additionnels :

* `state` : Sévérité/Criticité de l'alarme (0 = INFO, 1 = MINOR, 2 = MAJOR, 3 = CRITICAL) - Obligatoire
* `output` : Message court - Optionnel
* `long_output` : Message détaillé - Optionnel
* `close_delay` : Délai de résolution automatique en secondes - Optionnel 
* `perf_data` : [Métrique de performance](../../guide-utilisation/interface/widgets/graphiques/index.md)

```json
{
  "event_type": "check",
  "state": 2,
  "output": "CPU usage high",
  "long_output": "CPU usage over 90% for 5 minutes",
  "timestamp": 1639403732,
  "perf_data": "cpu=20%;80;90;0;100"
}
```

!!! info "Information close_delay"
    Le champ `close_delay` permet à une alarme d'être cloturée automatiquement après un délai dans le cas où aucune contre alarme ne serait émise vers Canopsis

!!! info "Information perf_data"
    Le format attendu pour le champ `perf_data` est le suivant :  
    ```text
    'label'=value[UOM];[warn];[crit];[min];[max]
    ```
    Exemple :
    ```json
    "perf_data": "cpu=20%;80;90;0;100"
    ```
     
    - Seules la **valeur** et l’**unité** sont prises en compte
    - Les valeurs `warn`, `crit`, `min`, `max` sont **ignorées**
    - Les métriques sont associées aux entités Canopsis via leur `entity_id`


#### Evénement `contextupdate`

Même structure que `check`, mais sans effet sur l'état de l'alarme.  
Permet d'enrichir une entité.

```json
{
  "event_type": "contextupdate",
  "output": "Update metadata",
  "timestamp": 1639403732
}
```


#### Evénements `ack`, `ackremove`, `assocticket`, `comment`, `cancel`, `uncancel`, `snooze`, `changestate`

Selon le type :  

* `author` : Auteur de l'action - Optionnel
* `user_id` : ID utilisateur - Optionnel
* `output` : Commentaire - Optionnel
* `ticket` : Numéro de ticket (assocticket uniquement)
* `state` : Nouvelle criticité forcée (changestate uniquement)
* `duration` : Durée en secondes (snooze uniquement)

Exemple `ack` :
```json
{
  "event_type": "ack",
  "author": "admin",
  "output": "Maintenance planned",
  "timestamp": 1639403732
}
```

!!! note "Note"
    Si `author` ou `user_id` ne sont pas fournis, Canopsis utilise ceux de l'utilisateur authentifié.

