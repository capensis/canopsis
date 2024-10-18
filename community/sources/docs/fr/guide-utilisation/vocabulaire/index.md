# Vocabulaire des termes de Canopsis

## Alarme

Une *alarme* est le résultat du traitement des [évènements](#evenement) par un [moteur](#moteur). Elle sert à signaler un problème.

Une alarme est liée à une [entité](#entite) de type [composant](#composant), [ressource](#ressource) ou [service](#service). La combinaison d'un [connecteur](#connecteur), d'un [nom de connecteur](#nom-de-connecteur), d'un composant et d'une ressource crée une alarme unique. Si l'un de ces éléments change, une alarme différente est créée.

Une alarme peut connaître de multiples changements de [criticité](#criticite), [priorité](#priorite) et de [statut](#statut), et subir une suite d'actions (acquittement, mise en veille, changement de criticité, annulation, etc.), [utilisateurs](../interface/widgets/bac-a-alarmes/actions.md) ou [automatiques](../menu-exploitation/scenarios.md). L'ensemble de ces changements et de ces actions constitue le *cycle d'alarme*.

Les alarmes peuvent être affichées à l'aide d'un [widget Bac à alarmes](../interface/widgets/bac-a-alarmes/index.md).

Vous pouvez consulter la [structure des alarmes](../../guide-developpement/base-de-donnees/periodical-alarm.md) présente dans le Guide de développement.

## Battement

Un [moteur](#moteur) effectue une tâche périodique appelée *battement* (ou *beat*) à un intervalle régulier. L'intervalle typique est de 1 minute.

## Composant

Un *composant* peut être soit :

* Un type d'[entité](#entite) créé après le traitement d'un [évènement](#evenement).
* Le champ `component` d'un [évènement](#evenement). Le plus souvent, il s'agit d'une machine ou d'un périphérique réseau (serveur, routeur, etc.). Une [alarme](#alarme) peut être rattachée à ce composant.

## Context-Graph

Le *context-graph* est un schéma relationnel entre les [entités](#entite) de Canopsis. Il sert à grapher leur contexte. Il s'appuie sur les notions de [`impact` et `depends`](../../guide-developpement/base-de-donnees/default-entities.md#context-graph). Il est présent au sein de chaque [entité](#entite) et est accessible au travers du [widget Explorateur de contexte](../interface/widgets/contexte/index.md).

## Connecteur

Un *connecteur* peut être soit :

* Un type d'[entité](#entite) créé suite au traitement d'un [évènement](#evenement). Il est le fruit de la concaténation des champs `connector` et `connector_name`.
* Le champ `connector` d'un évènement. Le plus souvent, il s'agit du nom du logiciel qui envoie ses données à Canopsis. Il sert à créer l'entité [connecteur](#connecteur).
* Un [script ou un programme](../../interconnexions/index.md#connecteurs) permettant d’envoyer à Canopsis des évènements à partir de sources d'informations extérieures.

## Criticité

Une [alarme](#alarme) a une *criticité*, indiquant la gravité de l'incident.

Il y a actuellement 4 criticités possibles :

* 0 - Info (quand en cours) / OK (quand résolue), de type stable.
* 1 - Mineure (*minor*), de type alerte.
* 2 - Majeure (*major*), de type alerte.
* 3 - Critique (*critical*), de type alerte.

## Enrichissement

L'*enrichissement* est l'action d'ajouter des informations supplémentaires à une structure de données.

On peut enrichir :

* Un [évènement](#evenement) via l'[event-filter du moteur `engine-che`](../menu-exploitation/filtres-evenements.md).
* Une [entité](#entite) via l'[event-filter du moteur `engine-che`](../menu-exploitation/filtres-evenements.md), l'[Explorateur de contexte](../interface/widgets/contexte/index.md) ou les [drivers](../../interconnexions/index.md#drivers).
* Une [alarme](#alarme) via le [moteur `engine-dynamic-infos`](../menu-exploitation/informations-dynamiques.md).

## Entité

Les *entités* servent à structurer les [alarmes](#alarme). Elles sont liées entre elles via le [context-graph](#context-graph). Elles peuvent permettre, via l'[enrichissement](#enrichissement) de conserver des données statiques (emplacement du serveur, nom du client, etc.).

Les entités sont accessibles au travers du [widget Explorateur de contexte](../interface/widgets/contexte/index.md).

Les entités ont les propriétés suivantes :

| Type d'entité | Résulte du traitement d'un [évènement](#evenement) | Peut être lié à une [alarme](#alarme)|
|---------------|--------------------------------------|---------------------------|
|[composant](#composant)|✅            |✅         |
|[connecteur](#connecteur)|✅          |❌                        |
|[service](#service)|❌                       |✅         |
|[ressource](#ressource)|✅            |✅         |

Vous pouvez consulter la [structure d'une entité](../../guide-developpement/base-de-donnees/default-entities.md) présente dans le Guide de développement.

## Évènement

Un *évènement* est un message arrivant dans Canopsis.

Il est formaté en JSON et peut être de plusieurs types, avec leurs propres structures. 

Les évènements de type `check` peuvent provenir d'une source externe, d'un [connecteur](../../interconnexions/index.md#connecteurs) ([email](../../interconnexions/Transport/Mail.md), [SNMP](../../interconnexions/Supervision/SNMPtrap.md), etc.) ou de Canopsis lui-même. Ils aboutissent à la création ou la mise à jour d'une [alarme](#alarme) dans le [Bac à alarmes](../interface/widgets/bac-a-alarmes/index.md).

## Impact

Une [entité](#entite) de [service](#service) a un *niveau d'impact* permettant de calculer la [priorité](#priorite) des [alarmes](#alarme) liées à l'entité.

Ce niveau d'impact permet aussi de définir la couleurs de l'alarme ou de la tuile liée au service dans la [météo de services](#meteo).

## Météo

La [*météo des services* est un widget](../interface/widgets/meteo-des-services/index.md) qui permet d'avoir une vue globale sur l'état d'un ensemble d'[entités](#entite). Pour cela, elle affiche des tuiles dont la couleur est représentative de la [priorité](#priorite) des [alarmes](#alarme) calculée par le [service](#service) lié.

## Moteur

Un *moteur* Canopsis consomme les [évènements](#evenement) entrants pour les traiter, puis les acheminer vers les moteurs suivants. Ils effectuent également une tâche périodique au [battement](#battement) et consomment leurs enregistrements en base de données lorsqu'ils sont disponibles.

## Nom de connecteur

Un *nom de connecteur* (ou `connector_name`) est le champ d'un [évènement](#evenement). Le plus souvent, il s'agit du nom du logiciel qui envoie ses données à Canopsis, complété par sa localisation ou sa numérotation (`superviseur_lille` ou `superviseur_5` par exemple). Il sert à créer l'entité [connecteur](#connecteur).

## Priorité

Une [alarme](#alarme) a une *priorité* qui est le produit de la [criticité](#criticite) d'une alarme et du niveau d'impact d'une [entité](#entite) liée. Cette priorité est recalculée par le [service](#service) lié à chaque changement de criticité de l'alarme.

## Ressource

Une *ressource* peut être soit :

- Un type d'[entité](#entite) créé suite au traitement d'un [évènement](#evenement). Il est le fruit de la concaténation des champs `resource` et `component`.
- Un champ d'un évènement. Le plus souvent, il s'agit du nom de la vérification effectuée (RAM, DISK, PING, CPU, etc.). Une [alarme](#alarme) peut être rattachée à une ressource.

## Service

Un *service* peut être soit :

* un type d'[entité](#entite) constituant l'arbre de dépendances, auquel peut être ajoutée une catégorie. Il s'agissait anciennement des watchers/observateurs.
* le nom de certains composants de Canopsis n'étant pas des [moteurs](#moteur). Par exemple, `canopsis-api` est un service et non pas un moteur, puisqu'il ne consomme pas d'évènements.
* dans le cadre d'une installation de type paquets, le nom d'une [unité systemd](https://access.redhat.com/documentation/fr-fr/red_hat_enterprise_linux/7/html/system_administrators_guide/sect-managing_services_with_systemd-unit_files) lançant un composant de Canopsis.

## Statut des alarmes

Une [alarme](#alarme) a un *statut*, indiquant la situation dans laquelle se trouve l'alarme indiquant un incident.

Il y a actuellement 5 statuts possibles :

* 0 - Fermée
    * Une alarme est considérée *fermée* (*off*) si elle est stable. C'est-à-dire que sa [criticité](#criticite) est stable à 0.
* 1 - En cours
    * Une alarme est considérée *en cours* (*ongoing*) si sa criticité est dans un état d'alerte (supérieur à 0).
* 2 - Furtive
    * Une alarme est considérée *furtive* (*stealthy*) si sa criticité est passée d'alerte à stable dans un délai spécifié.
    * Si la criticité de cette alarme est de nouveau modifiée durant le délai spécifié, elle est toujours considérée *furtive*.
    * Une alarme restera *furtive* pendant une durée spécifiée et passera à *fermée* si la dernière criticité était 0, *en cours* s'il s'agissait d'une alerte, ou *bagot* si elle se qualifie en tant que tel.
* 3 - Bagot
    * Une alarme est considérée *bagot* (*flapping*) si elle est passée d'une criticité d'alerte à un état stable un nombre spécifique de fois sur une période donnée.
* 4 - Annulée
    * Une alarme est considérée *annulée* (*cancel*) si l'utilisateur l'a signalée comme telle à partir de l'interface utilisateur.
    * Après une période donnée, une alarme marquée comme *annulée* changera de criticité pour être considérée comme *résolue*.

### Relations entre les statuts d'une alarme

```mermaid
sequenceDiagram
participant Fermée
participant En cours
participant Furtive
participant Bagot
participant Annulée
participant Résolue

Fermée ->> En cours:   
Note over Fermée, En cours: Envoi d'un événement <br/> de type check <br/> et criticité différente de 0

En cours ->> Furtive:   
Note over En cours, Furtive: Changement de criticité de stable à alerte <br/> une ou plusieurs fois pendant <br/> la durée de StealthyIntervale

En cours ->> Bagot:   
Note over En cours, Bagot: Changement de criticité de stable à alerte X fois durant Y secondes <br/> (où X est égal à la valeur de FlappingFreqLimit et Y à FlappingInterval)

En cours ->> Annulée:   
Note over En cours, Annulée: Envoi d'un événement de type cancel ou action utilisateur dans l'interface

En cours ->> Résolue:   
Note over En cours, Résolue: Envoi d'un événement de type done et délai de 15 minutes (valeur fixe)

Furtive ->> Bagot:   
Note over Furtive, Bagot: Le nombre de changements de criticité <br/> atteint la valeur du nombre d'oscillations

Furtive ->> Résolue:   
Note over Furtive, Résolue: Criticité stable à la fin de de la durée de StealthyInterval

Furtive ->> En cours:   
Note over Furtive, En cours: Criticité de niveau alerte à la fin de <br /> la durée de StealthyInterval

Bagot ->> En cours:   
Note over Bagot, En cours: Le délai depuis le dernier changement de criticité <br/> est supérieur à la période de bagot

Annulée ->> Résolue:   
Note over Annulée, Résolue: Automatique après un délai égal à <br/> la durée de CancelAutosolveDelay
```

Note : cliquez sur les liens suivants pour accéder aux informations relatives aux variables utilisées dans ce diagramme : [`StealthyInterval`](../../guide-administration/administration-avancee/modification-canopsis-toml.md#section-canopsisalarm), [`Nombre d'oscillations`](../../guide-utilisation/menu-exploitation/regles-bagot.md#anatomie-dune-regle-de-bagot), [`Durée de bagot`](../../guide-utilisation/menu-exploitation/regles-bagot.md#anatomie-dune-regle-de-bagot)
et [`CancelAutosolveDelay`](../../guide-administration/administration-avancee/modification-canopsis-toml.md#section-canopsisalarm).
