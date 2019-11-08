# Vocabulaire

## Alarme

Une *alarme* est le résultat de l'analyse des évènements. Elle historise et résume les changements d'état, les actions utilisateurs (acquittement, mise en pause, etc.).

## Bagot

Un évènement est considéré *Bagot* s'il est passé d'un état d'alerte à un état stable un nombre spécifique de fois sur une période donnée.

## Battement

Un moteur effectue une tâche périodique appelée *battement* (ou « beat » ) à un intervalle régulier de 1 minute.

## Composant

## Connecteur

## Enrichissement

L'*enrichissement* est l'action de récupérer des informations supplémentaires dans un évènement pour venir compléter le contexte (càd l'ensemble des entités).

## Entité

Une *entité* est une abstraction utilisée, entre autre, pour conserver des données statiques. Par exemple, l'entité associée à une alarme va contenir des informations qui ne varient pas d'une alarme à l'autre : emplacement du serveur, lien vers procédure liée à ce type d'alarme...

## État

## Évènement

Un *évènement* est un message arrivant dans Canopsis. Il est formatté en JSON et peut être de plusieurs [types](../../guide-developpement/struct-event.md). Les événements de type [`check`](../../guide-developpement/struct-event.md#event-check-structure) peuvent provenir d'une source externe, d'un [connecteur](../../interconnexions/index.md#connecteurs) ([email](../../interconnexions/Transport/Mail.md), [SNMP](../../interconnexions/Supervision/SNMPtrap.md), etc.) ou de Canopsis lui-même. Ils aboutissent à la création d'une alarme dans le bac à alarmes.

## Météo

La *météo des services* permet d'avoir une vue globale de plusieurs alarmes, par exemple en les regroupant par serveur, par type d'activité, ... C'est au choix, et complètement configurable :)

## Nom de connecteur

## Ressource

## Service

## Statut
