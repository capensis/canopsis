# Moteur `engine-action` (Community)

Le moteur `engine-action` permet de déclencher conditionnellement des actions sur des alarmes.

## Utilisation

En édition Community, la file du moteur est placée juste après le moteur [`engine-axe`](moteur-axe.md).

En édition Pro, la file du moteur est placée juste après le moteur [`engine-webhook`](moteur-webhook.md).

### Options du moteur

La commande `engine-action -help` liste toutes les options acceptées par le moteur.

## Fonctionnement

### Types d'action

Les types d'actions disponibles sont :

* `changestate`, qui correspond à un évènement `changestate` : change et verrouille la criticité de l'alarme jusqu'à sa résolution.
* `pbehavior`, met en place un [comportement périodique](moteur-pbehavior.md).
* `snooze`, qui correspond à un évènement `snooze` : pose une mise en veille automatique sur l'alarme.
* `ack`, qui correspond à un événement `ack` : pose un acquittement sur l'alarme.
* `ackremove`, qui correspond à un événement `ackremove` : supprime l'acquittement sur l'alarme.
* `assocticket`, qui correspond à un événement `assocticket` : associe un ticket à l'alarme.
* `declareticket`, qui correspond à un événement `declareticket` : déclarer un ticket pour l'alarme.
* `cancel`, qui correspond à un événement `ackremove` : annule l'alarme.

### Paramètres généraux

Une action est composée d'un JSON contenant les paramètres suivants :

* `_id` : optionnel. Identifiant de l'action. S'il n'est pas fourni par l'utilisateur il sera généré automatiquement. Le champ est de type `string`.
* `type` : obligatoire. Type d'action (voir [section précédente](#types-daction)). Ce champ est de type `string`.
* `parameters` : obligatoire. [Paramétrage spécifique](#parametres-specifiques) à chaque type d'action.
* `delay` : optionnel. Délai avant l'exécution de l'action. Les unités acceptées sont celles utilisées par le langage [Golang](https://golang.org/pkg/time/#ParseDuration) soit `s`, `m`, `h` pour secondes, minutes et heures respectivement. Le champ est de type `string`.
* `hook` : obligatoire. Il est composé des paramètres suivants :
    * [`patterns`](moteur-che-event_filter.md#patterns) : optionnel. Conditions sur les champs des alarmes (`alarm_patterns`), des entités (`entity_patterns`) ou des évènements (`event_patterns`) dans lesquelles l'action doit être appelée.
    * [`triggers`](../architecture-interne/triggers.md) : obligatoire. Ils servent comme point de déclenchement pour les actions automatisées, en général lors de la réception d'un évènement.

!!! attention
    Les [`triggers`](../architecture-interne/triggers.md) `declareticketwebhook`, `resolve` et `unsnooze` n'étant pas déclenchés par des évènements, ils ne sont pas utilisables avec les `event_patterns`.

### Paramètres spécifiques

#### Changestate

* `state` : obligatoire. Criticité dans laquelle sera verrouillée l'alarme (0 : Info, 1 : Minor, 2 : Major, 3 : Critical). Le champ est de type `integer`.
* `output` : optionel. Commentaire du changestate. Le champ est de type `string`.
* `author` : optionel. Auteur du changestate. Le champ est de type `string`.

#### PBehavior

* `rrule` : optionnel. Règle de récurrence pour le pbehavior. Le champ est de type `string`.
* `enabled` : obligatoire. Détermine si le pbehavior est actif ou non. Le champ est de type `boolean`.
* `author` : obligatoire. Auteur du pbehavior. Le champ est de type `string`.
* `name` : obligatoire. Nom du pbehavior. Le champ est de type `string`.
* `tstart` : obligatoire. Date de début du pbehavior. Le champ est de type `integer` en UNIX timestamp.
* `tstop` : obligatoire. Date de fin du pbehavior. Le champ est de type `integer` en UNIX timestamp.
* `type_` : obligatoire. Type du pbehavior. Le champ est de type `string`.
* `reason` : obligatoire. Raison du pbehavior. Le champ est de type `string`.
* `timezone` : optionnel. Timezone du pbehavior. Le champ est de type `string`.
* `comments` : optionnel. Commentaire du pbehavior. Le champ est de type `array`.
* `author` : optionnel. Auteur du commentaire. Le champ est de type `string`.
* `message` : optionnel. Commentaire du pbehavior. Le champ est de type `string`.
* `exdate` : optionnel. Date d'expiration, peut être composée de plusieurs valeurs. Le champ est de type `integer` en UNIX timestamp.

#### Snooze

* `message` : optionnel. Commentaire de la mise en veille. Le champ est de type `string`.
* `duration` : optionnel. Durée de la mise en veille en secondes. Le champ est de type `integer`.
* `author`: optionnel. Auteur de la mise en veille. Le champ est de type `string`.

#### Ack

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Ackremove

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Assocticket

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Declareticket

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.

#### Cancel

* `output`: optionnel. Commentaire de l'acquittement. Le champ est de type `string`.
* `author`: optionnel. Auteur de l'acquittement. Le champ est de type `string`.
