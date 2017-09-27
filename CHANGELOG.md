# Canopsis - Changelog

This document references all changes made to Canopsis since 2017/08/21. Ticket titles are kept in their original language.


## Canopsis 2.4.6 and CAT 2.5.0 (25/09/2017)


Canopsis 2.4.6 is a maintenance release for the 2.4 branch of canopsis.



### Functional changes - CAT


- [#393](https://git.canopsis.net/canopsis/canopsis/issues/393) Feat(Auth) : Compatibilité SAMLV2 
- [#375](https://git.canopsis.net/canopsis/canopsis/issues/375) Feat(SNMP) : Les traps SNMP anomalies ne remontent pas


###Bug fixes - CAT

- [#375](https://git.canopsis.net/canopsis/canopsis/issues/375) Fix(SNMP): Les traps SNMP anomalies ne remontent pas

### Functional changes

- [#394](https://git.canopsis.net/canopsis/canopsis/issues/394) feat(UI) : permettre l'ajout d'onglets dropdown dans la vue Header
- [#392](https://git.canopsis.net/canopsis/canopsis/issues/392) feat(Context-graph) : création d'une route ws  pour update du  context

###Bug fixes and other non-functional changes

- [#391](https://git.canopsis.net/canopsis/canopsis/issues/391) fix(Context-Graph) : La route post retourne parfois des doublons
- [#378](https://git.canopsis.net/canopsis/canopsis/issues/378) fix(web): Blocage appli 
- [#377](https://git.canopsis.net/canopsis/canopsis/issues/377) fix(Météo de service) : Impossible de faire les actions sur les alarmes
- [#376](https://git.canopsis.net/canopsis/canopsis/issues/376) fix(Météo de service) : La mise en pause d'un scénario en alarme ne remet pas en vert l'application
- [#374](https://git.canopsis.net/canopsis/canopsis/issues/374) fix(Bac à Alarmes) : Désynchro entre une alarme et son historique & fermeture de l'alarme pas toujours prise en compte
- [#349](https://git.canopsis.net/canopsis/canopsis/issues/349) fix(Météo de service) : Nouvelle desynchro entre entités sur 2.4.5
- [#347](https://git.canopsis.net/canopsis/canopsis/issues/347) fix(Météo de service) : création de pbehavior depuis le service weather
- [#371](https://git.canopsis.net/canopsis/canopsis/issues/371) fix(perfs) : Perf sur queue alerts en 2.4.5
- [#364](https://git.canopsis.net/canopsis/canopsis/issues/364) fix(pbehavior) : Engine pbehavior, déconnection et reconnection en boucle
- [#342](https://git.canopsis.net/canopsis/canopsis/issues/342) fix(pbehavior) : création lente
- [#333](https://git.canopsis.net/canopsis/canopsis/issues/333) fix(Bac à Alarmes) : les boutons d'action de masse ne fonctionnent que si 1 alarme est sélectionnée
- [#326](https://git.canopsis.net/canopsis/canopsis/issues/326) fix(métriques) : Route /api/context/metric et recherches par nom
- [#320](https://git.canopsis.net/canopsis/canopsis/issues/320) fix(pbehavior): Modale pbehavior - pas de création de pbheavior
- [#318](https://git.canopsis.net/canopsis/canopsis/issues/318) fix(pbehavior) : Il est possible de créer un pbehavior avec une rrule invalide
- [#317](https://git.canopsis.net/canopsis/canopsis/issues/317) fix(pbehavior) : check des rrules avant insertion


## Canopsis 2.4.5 (25/08/2017)


### Functional changes


- feat(Météo de service): amélioration serviceweather hauteur des tuiles
- [#345](https://git.canopsis.net/canopsis/canopsis/issues/345) trad(Météo de service) : Traduction française de la météo de Service
- [#337](https://git.canopsis.net/canopsis/canopsis/issues/337) feat(Météo de service) : afficher un compteur de temps avant le prochain changement sur une tuile en alarme
- [#323](https://git.canopsis.net/canopsis/canopsis/issues/323) feat(Météo de service) : templatiser les modales du widget
- [#314](https://git.canopsis.net/canopsis/canopsis/issues/314) feat(Bac à alarmes) : pouvoir afficher les infos de la resource affectée par l'alarme
- [#309](https://git.canopsis.net/canopsis/canopsis/issues/309) feat(Bac à alarmes) : ajouter dans les alarmes un champ de la cause d'alarme
- [#305](https://git.canopsis.net/canopsis/canopsis/issues/305) feat(baselines) : Intégration des baselines avec le service weather
- [#295](https://git.canopsis.net/canopsis/canopsis/issues/295) feat(alarmes) :  ajout de dates de création et de date de dernier changement dans une alarme


### Bug fixes and other non-functional changes:


- [#344](https://git.canopsis.net/canopsis/canopsis/issues/344) fix(engines) : Lors de la création d'un pbehavior avec le widget service weather, le cleaner event crash.
- [#339](https://git.canopsis.net/canopsis/canopsis/issues/339) fix(Mété de service) :  Erreur 404 sur api/v2/weather/watcher
- [#338](https://git.canopsis.net/canopsis/canopsis/issues/338) fix(session) : storage non instancié
- [#330](https://git.canopsis.net/canopsis/canopsis/issues/330) fix(Météo de service) : PBehaviors non fonctionnels
- [#329](https://git.canopsis.net/canopsis/canopsis/issues/329) fix(Bac à alarmes) : Le champs actions ne fonctionne pas si la colonne extra_details n'est pas affichée
- [#327](https://git.canopsis.net/canopsis/canopsis/issues/327) fix(Météo de service) : mixin customsendevent : conflit entre versions service-weather et listalarms
- [#324](https://git.canopsis.net/canopsis/canopsis/issues/324) fix(Bac à alarmes) : probleme de Timestamp
- [#310](https://git.canopsis.net/canopsis/canopsis/issues/310) fix(metrics) :  changer les adapters pour récupérer les métriques
- [#298](https://git.canopsis.net/canopsis/canopsis/issues/298) fix(pbehaviors): crash de l'engine pbehaviors lors de la creation de pbehavior |
- [#296](https://git.canopsis.net/canopsis/canopsis/issues/296) fix(runtime): amqp2engines* dans hypcontrol si on se trouve dans un dossier avec un amqp2engines.conf qui existe
- [#294](https://git.canopsis.net/canopsis/canopsis/issues/294) fix(global) : Internationalisation non fonctionnelle
- [#293](https://git.canopsis.net/canopsis/canopsis/issues/293) fix(Bac à alarmes) : La recherche ne fonctionne pas sur la 2.4.4
- [#291](https://git.canopsis.net/canopsis/canopsis/issues/291) fix(Bac à alarmes) : Le tri automatique sur les dates ne fonctionne pas
- [#290](https://git.canopsis.net/canopsis/canopsis/issues/290) fix(Bac à alarmes) : Pas de rafaichissement automatique de la vue lorsqu'une action est effectuée sur une alarme
- [#289](https://git.canopsis.net/canopsis/canopsis/issues/299) fix(snooze) : fonctionnalité inopérante
- [#287](https://git.canopsis.net/canopsis/canopsis/issues/287) refact(configuration): [interne] remplacer Configurable

## Canopsis 2.4.0 :

- feat(content-graph) : new database structure that stores the real-world topology of a supervised system as an object graph, allowing us to identify all entities impacted by an alarm
- feat (backend) : Canopsis now generates alarmes based on the status of events. This new object keep tracks of the full history of a real-world alarm
- feat(UI) : new "Service weather" UI brick that can display the status of up to 120 entities on a single window
- feat (UI): new "Alarms list" UI brick as a replacement of the old events list.


