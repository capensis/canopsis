# Journalisation des actions utilisateurs

Canopsis permet de journaliser certaines actions réalisées par un utilisateur dans un journal.

## Actions et types pris en charge

Actions  | Description  
--|---
`create`  | Création d'un objet (ex : Création d'un groupe de vues)
`update`  | Mise à jour d'un objet (ex : Mise à jour d'une règle Heartbeat)
`delete`  | Suppression d'un objet (ex : Suppression d'une règle de corrélation)


Types  | Description  
--|---
`user`                  | [Les utilisateurs](../../guide-utilisation/menu-administration/utilisateurs.md)
`role`                  | [Les rôles](../../guide-utilisation/menu-administration/roles.md)
`playlist`              | [Les listes de lecture](../../guide-utilisation/menu-administration/listes-de-lecture.md)
`eventfilter`           | [Les règles de filtrage/enrichissement](../../guide-utilisation/menu-exploitation/filtres-evenements.md)
`scenario`              | [Les règles de scénarios](../../guide-utilisation/menu-exploitation/scenarios.md)
`metaalarmrule`         | [Les règles de méta alarmes/corrélation](../../guide-utilisation/menu-exploitation/regles-metaalarme.md)
`dynamicinfo`           | [Les règles d'enrichissement d'alarmes](../../guide-utilisation/cas-d-usage/enrichissement.md)
`entity`<br/>`entityservice`<br/>`entitycategory` | [Gestion des entités](../../guide-utilisation/interface/widgets/contexte/index.md)
`pbehavior`<br/>`pbehaviortype`<br/>`pbehaviorreason`<br/>`pbehaviorexception`  | [Les comportements périodiques](../../guide-utilisation/cas-d-usage/comportements_periodiques.md)
`instruction`<br/>`job`<br/>`jobconfig`             | [Les objets de remédiation](../../guide-utilisation/remediation/index.md)
`statesetting`          | [Les paramètres](../../guide-utilisation/menu-administration/parametres-de-calculd-etat-sévérité.md)
`broadcastmessage`      | [La diffusion de messages](../../guide-utilisation/menu-administration/diffusion-de-messages.md)
`idlerule`              | [Les règles d'inactivité](../../guide-utilisation/menu-exploitation/regles-inactivite.md)
`view`<br/>`viewgroup`<br/>`viewtab`  | [Les vues](../../guide-utilisation/interface/vues/index.md)
`widget`<br/>`widgetfilter`<br/>`widgettemplate` | [Les widgets](../../guide-utilisation/interface/widgets/index.md)
`resolverule`           | [Les règles de résolution](../../guide-utilisation/menu-exploitation/regles-resolution.md)
`flappingrule`          | [Les règles de bagot](../../guide-utilisation/menu-exploitation/regles-bagot.md)
`kpi_filter`            | [Les KPI](../../guide-utilisation/menu-administration/kpi.md)
`pattern`               | [Les patterns](../../guide-utilisation/interface/widgets/bac-a-alarmes/index.md#filtres)
`map`                   | [Les cartographies](../../guide-utilisation/menu-administration/cartographie.md)
`snmprule`              | [Les règles SNMP](../../guide-utilisation/menu-exploitation/regles-snmp.md)
`declareticketrule`     | [Les règles de déclarations de tickets](../../guide-utilisation/menu-exploitation/regles-declaration-tickets.md)
`linkrule`              | [Les règles de générateurs de liens](../../guide-utilisation/menu-exploitation/generateur-liens.md)
`alarmtag`              | [Les tags d'alarmes](../../guide-utilisation/menu-administration/gestion-des-tags.md)
`eventrecord`           | [Les enregistrements d'évenements](../../guide-utilisation/menu-administration/enregistrements-d-evenements.md)
`externaldata`          | [Les données externes](../../guide-utilisation/menu-exploitation/donnees-externes.md)

## Récupérer les logs des actions depuis la base TimescaleDB

Par ailleurs, une table TimescaleDB, `action_log`, propose de conserver les actions effectuées sur l'ensemble des objets.

Le format de la table est le suivant:
Champs  | Description  
--|---
`id` |  ID interne à TimescaleDB.
`type`| L'action réalisée<br/>0 - create<br/>1 - update<br/>2 - delete 
`value_type` | Le type de l'action. 
`value_id` | L'identifiant de la tâche. 
`author` | L'auteur est l'utilisateur qui a effectué l'action.
`time` | Le moment où l'action s'est effectuée
`data` | L'ensemble des informations techniques sur l'action qui a été réalisée

Exemple d'un log:
```
psql (15.7)
Type "help" for help.

canopsis=> select * from action_log;
...
  5 |    0 | eventfilter | action_log_test                      | d7f523ac-ae68-4e95-92ec-9fc80abf3c28 | 2025-05-22 07:08:11 | {"_id": "action_log_test", "type": "enrichment", "rrule": "", "author": "d7f523ac-ae68-4e95-92ec-9fc80abf3c28", "config": {"actions": [{"name": "software", "type": "set_entity_info", "value": "canopsis", "description": "add software name"}], "on_failure": "pass", "on_success": "pass"}, "created": 1747897691, "enabled": true, "exdates": [], "updated": 1747897691, "priority": 0, "exceptions": [], "description": "This is an event filter rule", "event_pattern": [], "external_data": [], "entity_pattern": [[{"cond": {"type": "eq", "value": "capensis"}, "field": "infos.customer", "field_type": "string"}]], "resolved_exdates": null, "corporate_entity_pattern": "", "corporate_entity_pattern_title": ""}
...

```

