# Canopsis - Changelog

This document references all changes made to Canopsis since 2017/08/21. Ticket titles are kept in their original language.



## Canopsis 2.5.5 (sprint 01.19) - Due date : 2018-01-19

### Bug fixes 
 - [#579](https://git.canopsis.net/canopsis/canopsis/issues/579) - Impossible de créer un pbehavior dans l'explorateur de context
 - [#564](https://git.canopsis.net/canopsis/canopsis/issues/564) - [API] get-alarms ne remonte pas tous les résultats
 - [#563](https://git.canopsis.net/canopsis/canopsis/issues/563) - [HardLimit] la hardlimit empêche toutes les actions sur une alarme
 - [#485](https://git.canopsis.net/canopsis/canopsis/issues/485) - bac a alarme création de pbhavior pop up calendrier qui ne s'affiche pas

 ### Functional and other changes
 - [#573](https://git.canopsis.net/canopsis/canopsis/issues/573) - Bac à alarme - recherche insensible à la casse
 - [#526](https://git.canopsis.net/canopsis/canopsis/issues/526) - Pouvoir trier les tuiles du widget service weather
 - [#525](https://git.canopsis.net/canopsis/canopsis/issues/525) - [Météo] Remonter les statistiques d'un scénario
 - [#524](https://git.canopsis.net/canopsis/canopsis/issues/524) - [Météo de service]Disposer d'une information sur la date de l'alarme



## 2.5.4 (Sprint 01.6) - Due date : 2018-01-06

### Bug fixes
 - [#543](https://git.canopsis.net/canopsis/canopsis/issues/543) - [Engine Alerts] Lorsqu'une alerte a atteint sa hard limit, le beat processing plante dans check_alarm_filters
 - [#540](https://git.canopsis.net/canopsis/canopsis/issues/540) - [doc] lancer un filldb update à chaque mise à jour de Canopsis
 - [#538](https://git.canopsis.net/canopsis/canopsis/issues/538) - [docker/CAT] La brique SNMP n'est pas installée
 - [#537](https://git.canopsis.net/canopsis/canopsis/issues/537) - [CRUD Context] Le schéma d'édition d'une entité n'est pas bon
 - [#556](https://git.canopsis.net/canopsis/canopsis/issues/556) - Declare ticket: missing number
 - [#555](https://git.canopsis.net/canopsis/canopsis/issues/555) - [PE] Dysfonctionnement du workflow des alarmes filter
 - [#553](https://git.canopsis.net/canopsis/canopsis/issues/553) - Missing display_name
 - [#548](https://git.canopsis.net/canopsis/canopsis/issues/548) - Recherche naturelle non fonctionnelle sur les display name
 - [#546](https://git.canopsis.net/canopsis/canopsis/issues/546) - Probleme de couleur sur les tuiles en MAJOR
 - [#545](https://git.canopsis.net/canopsis/canopsis/issues/545) - Pb de désynchronisation de statut entre scénario et Application
 - [#544](https://git.canopsis.net/canopsis/canopsis/issues/544) - Pb de bagot sur les alarmes. Elles ne se cloturent plus.
 - [#542](https://git.canopsis.net/canopsis/canopsis/issues/542) - Pbehavior crash when filter is in a bad format
 - [#539](https://git.canopsis.net/canopsis/canopsis/issues/539) - [Bac à alarmes] Des pbehaviors expirés remontent sur les alarmes en cours
 - [#531](https://git.canopsis.net/canopsis/canopsis/issues/531) - Problème avec la recherche naturelle
 - [#513](https://git.canopsis.net/canopsis/canopsis/issues/513) - le formulaire de login  doit catch le 401


### Functional and other changes
 - [#516](https://git.canopsis.net/canopsis/canopsis/issues/516) - Dummy authentication
 - [#536](https://git.canopsis.net/canopsis/canopsis/issues/536) - Ajouter le résutlat des tests dans le template
 - [#535](https://git.canopsis.net/canopsis/canopsis/issues/535) - Ajouter prérequis dans le template
 - [#533](https://git.canopsis.net/canopsis/canopsis/issues/533) - Retirer le bouton PAUSE sur les alarmes CLOSED
 - [#541](https://git.canopsis.net/canopsis/canopsis/issues/541) - Nettoyage des test



## Canopsis 2.5.3 (Sprint 12.15) Due date : 2017-12-15

**Not released due to blocking issue. This release was tagged on gitlab but not distrubuted. All issues were reported in the 2.5.4 release**


## Canopsis 2.5.2(11.30) - Due date : 2017-12-05

### Bug fixes 

 - [#499](https://git.canopsis.net/canNombre d'issues pour la milestone  "2.5.1 (11.03)" : 11
opsis/canopsis/issues/499) - [PE] la météo ne s'affiche pas pour les applications "standard"
 - [#518](https://git.canopsis.net/canopsis/canopsis/issues/518) - Saml2 import error
 - [#515](https://git.canopsis.net/canopsis/canopsis/issues/515) - Erreur 500 avec avec un identifiants inconnu
 - [#509](https://git.canopsis.net/canopsis/canopsis/issues/509) - Erreur de login webserver
 - [#497](https://git.canopsis.net/canopsis/canopsis/issues/497) - attribut creation_date  en avance par rapport aux dates des évènements sur certaines alarmes.
 - [#450](https://git.canopsis.net/canopsis/canopsis/issues/450) - Probleme de fonctionnement sur la recherche dans la bac à alarmes


### Functional and other changes

  - [#510](https://git.canopsis.net/canopsis/canopsis/issues/510) - Ajouter une date de dernier événement reçu à chaque alarme.
  - [#466](https://git.canopsis.net/canopsis/canopsis/issues/466) - Pouvoir disposer d'id d'alarmes exploitable dans le bac à alarmes



## Canopsis 2.5.1 (11.03) - Due date : 2017-11-03


### Bug fixes


 - [#446](https://git.canopsis.net/canopsis/canopsis/issues/446) - Le paquet dm.xmlsec.binding-1.3.3.tar.gz ne s'installe pas sur CentOS 7
 - [#427](https://git.canopsis.net/canopsis/canopsis/issues/427) - watcher ne se recalcule pas a la fin d'un pbehavior => desync
 - [#383](https://git.canopsis.net/canopsis/canopsis/issues/383) - soucis utf-8 sur les trap laposte
 - [#372](https://git.canopsis.net/canopsis/canopsis/issues/372) - Pas de données sur la route /perfdata
 - [#470](https://git.canopsis.net/canopsis/canopsis/issues/470) - Problème de performance de la route trap
 - [#451](https://git.canopsis.net/canopsis/canopsis/issues/451) - La configuration du popup du widget alarme ne conserve pas les données
  - [#399](https://git.canopsis.net/canopsis/canopsis/issues/399) - probleme utf8  quand on met un display name avec un accent
 - [#432](https://git.canopsis.net/canopsis/canopsis/issues/432) - UI ne charge pas
 - [#430](https://git.canopsis.net/canopsis/canopsis/issues/430) - bug(CRUD context-graph): Crudcontext adapter en doublon
 - [#413](https://git.canopsis.net/canopsis/canopsis/issues/413) - Probleme d'affichage sur le bac à alarmes
 - [#412](https://git.canopsis.net/canopsis/canopsis/issues/412) - Formatage des heures dans l'historique du bac à alarmes

 - [#401](https://git.canopsis.net/canopsis/canopsis/issues/401) - [Bac à alarmes] Lorsqu'un filtre est présent, le premier chargement de la vue échoue.


### Functional and other changes

 - [#472](https://git.canopsis.net/canopsis/canopsis/issues/472) - MAJ Image Docker MongoDB
 - [#406](https://git.canopsis.net/canopsis/canopsis/issues/406) - Installation mono-package python
 - [#478](https://git.canopsis.net/canopsis/canopsis/issues/478) - [Alarms] Optimisation des perfs du beat processing
 - [#357](https://git.canopsis.net/canopsis/canopsis/issues/357) - feat(context-graph): Basic editing of an entity
 - [#356](https://git.canopsis.net/canopsis/canopsis/issues/356) - feat(CRUD context-graph) : list entities
 - [#336](https://git.canopsis.net/canopsis/canopsis/issues/336) - Remplacer les Configurable par confng
 - [#403](https://git.canopsis.net/canopsis/canopsis/issues/403) - [Bac à alarmes] Recherche naturelle
 - [#402](https://git.canopsis.net/canopsis/canopsis/issues/402) - [Météo de service] Ajouter un bouton FastACK sur la popup "scénario"

## Canopsis 2.4.6 and CAT 2.5.0 (25/09/2017)


Canopsis 2.4.6 is a maintenance release for the 2.4 branch of canopsis.



### Functional changes - CAT


- [#393](https://git.canopsis.net/canopsis/canopsis/issues/393) Feat(Auth) : Compatibilité SAMLV2 
- [#375](https://git.canopsis.net/canopsis/canopsis/issues/375) Feat(SNMP) : Les traps SNMP anomalies ne remontent pas


###Bug fixes - CAT

- [#375](https://git.canopsis.net/canopsis/canopsis/issues/375) Fix(SNMP): Les traps SNMP anomalies ne remontent pas

### Functional  and other changes

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


