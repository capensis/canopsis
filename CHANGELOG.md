# Canopsis - Changelog

This document references all changes made to Canopsis since 2017/08/21. Ticket titles are kept in their original language.

## Unreleased

### Experimental features

- [Alarms list]: Alarms on resources can be hidden when their parent component is down.

## Canopsis 3.7.0 - Due date : 2018-12-27

 - [Documentation] Add documentation for ackcentreon task
 - [Tooling] Fix debian installation method
 - [UI] Add tabs system inside views
 - [UI] Add default views feature for roles and users
 - [UI/ServiceWeather] Add a parameter that, when clicking on a service weather tile, allows choosing between opening a modal with the observer's entities list and opening an alarm list related to the observer
 - [UI/ServiceWeather] Fix random screen freezing
 - [UI/Alarm List] Fix bug that fetched the alarms twice
 - [UI/Context] Fix timestamp problem at pbehavior creation
 - [UI/Context] Add a search bar used to search through an entity's informations
 - [UI/Filters] Add a CRUD managing the event filter rules

## Canopsis 3.6.0 - Due date : 2018-12-13

 - [Documentation] Add documentation for global status check with the healthcheck route
 - [Documentation] Add documentation for mix-filters and entity duplication
 - [Documentation] Various cosmetic improvements
 - [Go] Fix requests to the Observer ticketing API
 - [Go] Fix statecounter steps handling
 - [Tooling] Fix error handling in canopsinit
 - [Tooling] Fix issue with unused parameters in init
 - [UI] Add default views and rights for UIv3
 - [UI/Context] Add "Manage infos" panel for watchers
 - [UI/Context] Fix a bug with resources expand
 - [UI/Context] Add "Clone" action on entities and watchers
 - [UI/Context] Filter required on watcher's creation (at least one valid rule in the filter)
 - [UI/Rights] Fix a bug with confirmation panel not closing when submitting rights
 - [UI/Events] Add an `origin: 'canopsis'` parameter with all events coming from Canopsis UI
 - [UI/Version] Add Canopsis version number on side-bar
 - [UI/Filters] Add "Mix filters" feature
 - [UI/Alarm list] Simplify default sort column selector in settings
 - [UI/Top bar] Fix a bug with group editing on the top bar

## Canopsis 3.5.0 - Due date : 2018-11-29

- [Documentation]: Add a new documentation
- [Python]: Add a new route to fetch a list of entities with their current alarm
- [Python]: Automatically recover from the loss of the primary member of a MongoDB replicaset
- [Go]: Prevent a crash when a snooze has no duration
- [Go]: Add an option to automatically acknowledge the resources below a component
- [Go]: End the implementation of the eventfilter service
- [Go]: Automatically create a ticket when a new alarm is created if the flag autoDeclareTickets is given to the axe engine
- [ServiceWeather]: Fix the message set in the events sent when an action is triggered
- [Tooling]: Update the configuration of catag to handle the new canopsis project
- [Tooling]: Add a VERSION.txt file inside canopsis, display it on the prompt inside the canopsis env and add an API to retrieve it through HTTP
- [Tooling]: Add the missing engine-action in the push_docker_image.sh script

## Canopsis 3.4.0 - Due date : 2018-11-15

- [Go]: Fix make release command
- [Go]: IncidentKey in Service Ticket API is now optional
- [Go]: Introducing new Observer API driver for ticket creation
- [Go]: New actions for event-filter engine (copy and set_entity_info_from_template)
- [Python]: Fix /get-alarms route which was limited to 50 elements, and returned an overestimation of total
- [Python]: New Healthcheck API to check Canopsis status (service connexions, engine statues...)
- [Python]: Update Heartbeat engine docstring
- [ServiceWeather]: Fix paused pbehavior icon not always correctly display

## Canopsis 3.3.0 - Due date : 2018-10-31

- [Alarm List]: Fix bug where selected filter is not saved in userpreferences
- [Chore]: Some configuration cleanup in sources, configs and docker files
- [Docker]: Change rabbitmq base image and add some new envvars for image version selection
- [email2canopsis]: Converter now invalidate a translation if there is not match
- [Go]: Change build dir (from /tmp to ./build)
- [Go]: New event-filter engine ! (to filter/translate events based on some rules)
- [Go]: New ticket engine ! (to create tickets through external APIs)
- [Go]: Standardization of heartbeat ids generation
- [Python]: Add new event-filter and ticket API (for new engines)
- [ServiceWeather]: Add default message on Validate and Ack actions
- [snow2canopsis]: New connector which read Service Now API to import informations into Canopsis context
- [UI]: Fix canopsis-next not correctly builded with docker

## Canopsis 3.2.5 - Due date : 2018-10-17

- [Doc]: translating recent Upgrade documentations to english
- [Docker] : new CPS_LOGGING_LEVEL envvar to change loglevel in dockerized engines
- [email2canopsis] : now decode encoded subject line
- [Go]: new Action engine (especially pbehaviors from regex)
- [Go]: send alarm resolution informations to stats engine
- [Python] : fix « Socket Error 104 » while engines communicate with rabbitmq
- [Tool] : env2cfg can now handle mongo replicaset option
- [UI] : add rights for items in engine menu
- [UI] : fix `has_active_pb` flag not correctly calculated
- [UI] : fix timeline not correctly showing all informations (especially on ack, ticket and canceled events)
- [UI] : integrate listalarm, timeline and querybuilder bricks in central repo

## Canopsis 3.2.4 - Due date : 2018-10-05

- [Chore]: add some system dependencies to ansible installation
- [Go]: fix bagoting alarms never closed if cropped
- [Go]: send axe statistics to stasng engine
- [Python]: change some amqp publishers to pika to prevent odd reconnections
- [Python]: fix has_active_pb flag no corretly show all linked pbehaviors
- [Python]: fix pbehavior not corretly handle timezone change (one day gap)
- [Python]: fix performance concern on alert consultation

## Canopsis 3.2.3 - Due date : 2018-09-17

- [Email Connector]: handle base64 encoded parts
- [Email Connector]: option which use redis to resend last known event on run
- [Service weather]: Fix icons consistency between a watcher and his modal

## Canopsis 3.2.2 - Due date : 2018-09-12

- [Go]: Fix che event enrichment with new entities
- [Go]: Fix unused debug flag
- [Python]: Fix pbehavior based desynchronisation on weather route

## Canopsis 3.2.1 - Due date : 2018-09-12

- [CAT] Simplification of the statistics API (!74, !77, !78)
- [CAT] Add trends, SLA, sorting, aggregations and filtering by author to the statistics API  (!75, !76, !80, !87)
- [CAT] Add support for monthly periods to the statistics API (!79)
- [CAT] Add statistics for ongoing alarms and current state of an entity (!81, 82)
- [CAT] Fix an issue with the state of the alarms for statistics on resolved alarms (!83)
- [CAT] Add API route to get the history of the state of an entity (!85)
- [Python]: Fix error on /trap route
- [Python]: Fix amqp driver to revive stucked Consumers
- [Service weather]: Fix wrong pbehavior and maintenance icons used
- [UI]: Fix role default view that cannot be modified

## Canopsis 3.2 - Due date : 2018-09-01

- [Python]: fix rruled pbehavior computation
- [Python]: two bugfix on /get-entities route
- [Service weather]: fix unfoldable item when it contains a %
- [UI]: add ellipsis on hostgroups field
- [UI]: add rights on filters (create, list, modify)
- [UI]: hide "Restore alarm" button on open alarm list

## Canopsis 3.1 - Due date : 2018-08-17
- [Python]: add a feature that track the change of the longoutput field and alter the behavior of the output field of an alarm.
- [Python]: fix a bug that duplicate an alarm.
- [Build]: fix an error during the compilation of canopsis-next.
- [UI]: fix a right issue in the listalarm brick that allow the done button to be displayed without any right.

## Canopsis 3.0.1 - Due date : 2018-08-03

- [CAT]: New stats in statsng: most_alarms_impacting, worst_mtbf, alarm_list, longest_alarms (!67 et !70)
- [CAT]: Add a "periods" parameter to the statsng API (!69)
- [CAT]: New stats in statsng: most_alarms_created, alarms_impacting, state_list (!71)
- [CAT]: send stats with alarm duration
- [CAT]: default sats value
- [Snooze]: faster snooze application
- [Periodic refresh]: fix infinite loop with periodic refresh
- [Link list]: remove linklist
- [Links]: fix links in service weather
- [UI]: views in canopsis-next
- [UI]: alarms in canopsis-next

## Canopsis 3.0.0 - Due date : 2018-07-20

- [Chore]: Finalizing new building process
- [CAT]: New stats in statsng: time_in_state, availability, maintenance and mtbf
- [CAT]: Provisionning for docker demos
- [Deploy]: Automated standalone deployment
- [Docker]: Multiple cleanup and fixes
- [Go]: Fix an uncatched exception in alarm
- [Go]: Large refacto on engines and services
- [Go]: New engine interface
- [Go]: New Done action
- [Go]: Support watcher
- [Python]: Huge code cleanup, upgrade to pymongo 3.6
- [Python]: New Done action
- [Python]: New metric engine
- [Python]: New stat event: statstateinterval
- [Python]: New webserver handler using Flask
- [Python]: Small refacto and optimization on linkbuilder classes
- [UI]: Beginning Canopsis-next integration (context view)
- [UI]: Fix brickmanager when used with canopsis user
- [UI]: New documentation subprocess (visible on search field in Alarm list)
- [UI]: New 'unknown' state marker
- [Alarm list]: Natural search only on visible columns
- [Alarm list]: Possibility to search on ticket number
- [Service weather]: Add rights for actions
- [Service weather]: Fix blink on pbehaviored watchers
- [Service weather]: Fix refresh (again)

## Canopsis 2.7.0 - Due date : 2018-06-28

- [API]: new routes to manage frontend views
- [CAT]: add routes and make some bugfixes for the new engine statsng
- [CAT]: fix priviledge escalation and introduction of a default right group
- [Go]: large service/adapter refacto
- [Go]: fix timeout on mongo requests
- [Service weather]: fixed an issue where periodical refresh process never refresh anymore
- [Service weather]: add permission on action
- [UI]: a new 'done' action is avalaible, which simply mark an alarm as done
- [UI]: fix /sessionstart route
- [UI]: fix snmp view that can be broken
- [Webserver]: authentification trough WebSSO

## Canopsis 2.6.8 - Maintenance release - Due date : 2018-06-20

- [Service weather]: fixed an issue where the popup overlay could stay when the popup was closed, freezing the view
- [CAT]: ported the email2canopsis connector to python3 and fixed an issue with pattern matching caused by the new python version.
- [pbehaviors]: refactored the pbheaviors internal API to unify processing and avoid inconsistencies

## Canopsis 2.6.7 - Maintenance release - not released

This release was replaced by the 2.6.8 version to integrate an urgent fix on the pbehaviors.

## Canopsis 2.6.6 - Maintenance release -  Due date : 2018-06-08

- [Connector]: email2canopsis now with python3
- [Docker]: add some usefull tools in debian-9 docker image
- [EventFilter]: now can filter on subkeys
- [Go]: fix crash when a non-check event arrive on an empty alarm
- [Go]: add printEventOnError flag on engines
- [Go]: update dependency, and notablly migrating from mgo to globalsign/mgo
- [Go]: fix event-entity not always correctly associated
- [Go]: fix PBehavior event message never acked in Che
- [Go]: fix inaccurate messages in Che
- [Python]: API can now produce correct HttpError on catched exception (instead of 200)
- [Python]: limit /get-alarms response size by removing steps
- [Python]: upsert mode on import context route
- [Python]: alarms without a corresponding entity are not shadowed anymore
- [Python]: fix active pbehavior search
- [Python]: two new events (statcounterinc and statduration) for statistics purpose
- [Python]: remove ticket engine from amqp2engines.conf
- [UI]: fix BNF research on entity properties
- [UI]: fix sort on current_state_duration
- [Service Weather]: add new rights for actions in service weather


## Canopsis 2.6.5 - Maintenance release -  Due date : 2018-05-18

- [Go]: Porting steps cropping from python
- [Go]: Huge refactoring on error handling
- [Go]: Events with a timestamp as a float are now correctly parsed
- [Python]: step cropping compatibility update
- [Python]: bugfix where a source_type component is treated as a resource
- [UI]: fix context info presentation in the context explorer
- [UI]: bugfix on hide resource
- [UI]: fix navbar rights
- [Service Weather]: fix possibility to ack twice, change invalidate to cancel on thumb down button
- [Alarms list]: creation date is now presented in full format, even if its the present day


## Canopsis 2.6.4 - Maintenance release -  Due date : 2018-05-03

- [Setup] A version of the amqp2engines.conf file is provided for the High performance engines
- [Go]: fix error code when an engine crashes
- [Go]: pbehavior events are now correctly parsed by Che and passed to the pbehavior engine
- [Go]: generation of an unique id on every alarm
- [Go]: fix a crash when impacts and depends are not initialized in an entity object
- [Go]: Removed some useless logs from the Context
- [Event filter]: fixed an issue that blocked snooze events
- [Heartbeat]: fix heartbeat never closed on special conditions
- [Stat]: perfomance update on the engine
- [Stat]: fix handling special character in requests
- [UI]: fix presentation of pbehavior (in service-weather), ellipsis and snooze
- [UI]: fix permissions on the "List pbehaviors" action
- [CAT]: fix a crash with datametrie engine
- [Alarms list]: the  'pbehavior' icon is now displayed only when a pbeheavior is active on the entity, not when it is configured
- [Alarms list]: add the snooze End date on the Snooze tooltip
- [Alarms list]: several fixes on the ACK and report ticket action
- [Alarms list]: Pbehaviors can now be filtered.
- [Alarms list]: fix an issue where a pbehavior id could become null with periodic refresh enabled
- [Service weather]: Resources are now greyed out when a Component has a pbehavior
- [Context Graph]: Fixed an issue that prevented the expansion of an item of the Context Graph explorer


## Canopsis 2.6.3 - Maintenance release -  Due date : 2018-04-26

- [Engines] : Added UNACK, Uncancel, keep state actions to the High performance engines
- [Engines] : Added a default author in all alarms steps in the High performance engines
- [Go]: gracefull initialization of docker env and fix amqp2engines.conf on docker
- [Go]: Context and alarm creation performance improvements
- [Go]: code refactoring, fix alarm duplication and preloading
- [UI]: update rights on default view
- Multiple fixes on Pbheavior and Event_filter


## Canopsis 2.6.2 - Maintenance release -  Due date : 2018-04-23

### Fixes

- [Alarms list] Fixed a date formatting issue in the alarms list that made the `last_update_date` column appear with a 1 month delay
- [Service weather] reworked the Ticket action to fix a display issue caused by the new "save on exit" workflow
- [Alarms list] Mass actions now correctly get their rights/permissions (inherited from the rights applied on the actions on the single alarm)
- [setup] Canopsinit now requires a flag `--authorize-reinit` to perform any destructive modification to the database as an extra security
- [APIs] the "enabled" flag on all entities is now active
- [Engines] removed alarms caching in the Che engine to avoid  alarms duplication
- [Rights management] : fixed a rights issue with the massive actions on a limited account
- [Service weather] : fixed the components display to put long names on 2 lines instead of truncating it


## Canopsis 2.6.1 - Due date: 2018-04-20
**Not released due to regression**


## Canopsis 2.6.0 - Due date: 2018-04-18

This release introduced the new High performance engines and allowed the renaming on the Alarms list columns.

### Functional changes

- [Alarms list] : columns can now be renamed by the user
- [Alarms list] : pbehaviors can now be filtered
- [Alarms list] : added rights management to some new actions that were missing it
- [Alarms list] : Added "duration" and "current_state_duration" columns that display the total duration of the alarm and the duration of the current state of the alarm (respectively)
- [Alarms list] : the date is now displayed even if the the alarm was created today
- [pbehaviors] : the form has been simplified to give a behavior closer to a Calendar event
- [setup] : replaced the `schema2db` and `canopsis-filldb` commands with the new Canopsinit command (**Warning: see UPGRADING_2.6.md document**)
- [CAT] : The Datametrie connector can now filter some alarms based on their criticity
- [CAT] : The datametrie connector can now use the local date

### Experimental features

- [Engines] : New High Performance engines for heavily loaded environments (**experimental**)
- [Engines] : reimplemented the last_event_date feature on the High performance engines
- [engines] : The engine "stats" (High performance version) can now log the actions and their autors for audit purposes

### Bug fixes

- [Alarms list] : fixed an issue where rights were not saved properly in the admin
- [Alarms list] : fixed an issue that prevented pbehaviors to be saved properly on the High performance engines
- [Alarms list] : fixed an issue that could record a ticket number with the `0` value with the High performance engines
- [Engines] fixed an issue where the cancel Action did not close the alarm with the High performance engines


## Canopsis 2.5.12 (Sprint 03.16) - Due date : 2018-03-16

### Functional changes

- [Service Weather] Added the Fast ACK action in the service Weather widget
- [Service Weather] Added an ACK icon in the service weather titles, to identify entities that have acknowledged alarms
- [Service Weather] Added OK/NOK events statistics in the info popup
- [APIs] APIs can now be requested using Basic Auth, without having to request the authentication route first
- [Timeline] automatic actions now have an Author: "system"
- [list alarms] added "Cancel alarm" and "Change criticity" actions, with correct translation
- [Event filter] added "state" and "state_type" fields in the usable fields list
- [Installation] New Installation method based on RPM/deb packages

### Bug fixes

- [Service weather] Fixed an issue where the watcher state was incorrectly impacted by paused applications
- [Service weather] Fixed a type issue that could prevent the UI to display correctly
- [HA Tools] Fixed an issue of the High Performance engines when the RabbitMQ connection was lost
- [HA Tools] Reduced the downtime of the Canopsis UI when the MongoDB primary instance changes
- [CAT] Fixed an SAMLv2 installation issue on CentOS 7/RHEL 7
- [CAT] Fixed an race condition where SNMP rules could be missing when a trap was received while the engine's rules update function was running
- [CAT] Fixed an issue where an SNMP trap with a routing key > 255 chars could crash the SNMP engine and block the whole AMQP queue

## Canopsis 2.5.6 (sprint 02.2) - Due date : 2018-02-02

### Bug fixes
 - [#598](https://git.canopsis.net/canopsis/canopsis/issues/598) - Certaines entités ne sont pas importées correctement
 - [#593](https://git.canopsis.net/canopsis/canopsis/issues/593) - Pouvoir interroger les alarmes Resolved avec les dates des alarmes
 - [#589](https://git.canopsis.net/canopsis/canopsis/issues/589) - Probleme Dataset sur widget ServiceWeather
 - [#580](https://git.canopsis.net/canopsis/canopsis/issues/580) - Stealthy calculation
 - [#568](https://git.canopsis.net/canopsis/canopsis/issues/568) - Modifier la couleur d’affichage de la plate-forme de pré-prod
 - [#565](https://git.canopsis.net/canopsis/canopsis/issues/565) - Le bouton "Remove alarm (poubelle/corbeille)" n'apparaît pas sur les alarmes critiques
 - [#558](https://git.canopsis.net/canopsis/canopsis/issues/558) - Recherche avec des int
 - [#529](https://git.canopsis.net/canopsis/canopsis/issues/529) - Pouvoir supprimer des alarmes
 - [#528](https://git.canopsis.net/canopsis/canopsis/issues/528) - Disposer d'une information sur la date de l'alarme

### Functional and other changes
 - [#599](https://git.canopsis.net/canopsis/canopsis/issues/599) - Nettoyage des engines
 - [#566](https://git.canopsis.net/canopsis/canopsis/issues/566) - Remapper l'output "Lost 100%" en "Equipement injoignable"
 - [#594](https://git.canopsis.net/canopsis/canopsis/issues/594) - Validation des actions dans le popup MDS



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
