## PBehaviors

Ou `Periodical Behaviors`.

Il ne permettent pour le moment que de créer un équivalent des « downtimes », à savoir l’information comme quoi une entité n’est pas active.

### Caractérisation

Un PB se caractérise par les informations suivantes :

 * `name`: `string` type de pbehavior : `downtime` est le seul supporté
 * `filter`: `string` filtre mongo sur les entités du context graph
 * `author`: `string` auteur du PB
 * `enabled`: `bool` activer ou désactiver le PB, pour qu’il ne soit pas traité même si on se trouve pendant la période du PB
 * `comments`: `null` ou liste de commentaires.
 * `rrule`: `string` champs texte défini par la [RFC 2445](https://tools.ietf.org/html/rfc2445)
 * `tstart`: `int` timestamp fournissant la date de départ, recalculé à partir de la RRULE si présente.
 * `tstop`: `int` timestamp fournissant la date de fin, recalculé à partir de la RRULE si présente.
 * `type_`: `str` (optionnel) type de pbehavior (pause, maintenance...).
 * `reason`: `str` (optionnel) raison pour laquel ce pbehavior est posé.

Dans le cas où la RRULE est absente, `tstart` et `tstop` font office de plage d’activation du PB, sans récurrence.

Dans le cas où la RRULE est présente, `tstart` et `tstop` seront recalculés afin de refléter la récurrence.

### Création

 * API HTTP : voir l’API `/pbehavior/create`
 * Event de type `pbehavior` : créer à partir des champs cités en introduction
