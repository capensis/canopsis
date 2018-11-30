## Caractérisation

Un pbehavior se caractérise par les informations suivantes.

| Champ  | Type   | Description |
| -------| ------ | ----------- |
| `_id` | ? | Identifiant unique du comportement, généré par MongoDB lui-même. |
| `eids` | liste | Liste d'identifiants d'entité qui correspond au filtre précédent. |
| `name` | string | Type de pbehavior. `downtime` est la seule valeur acceptée. |
| `author` | string | Auteur ou application ayant créé le pbehavior. |
| `enabled`| bool | Activer ou désactiver le pbehavior, pour qu’il puisse être ignoré, même sur une plage active. |
| `comments` | liste | `null` ou une liste de commentaires. |
| `rrule` | string | Champ texte défini par la RFC 2445. |
| `tstart` | int | Timestamp fournissant la date de départ du pbehavior, recalculée à partir de la `rrule` si présente. |
| `tstop` | int | Timestamp fournissant la date de fin du pbehavior, recalculée à partir de la `rrule` si présente. |
| `type_` | string | Optionnel. Type de pbehavior (pause, maintenance…). |
| `reason` | string | Optionnel. Raison pour laquelle ce pbehavior a été posé. |

Un exemple d'évènement pbehavior brut :
```js
{
   "_id" : string,
   "name" : string,
   "filter" : string,
   "comments" : [ {
       "_id": string,
       "author": string,
       "ts": timestamp,
       "message": string
   } ],
   "tstart" : timestamp,
   "tstop" : timestamp,
   "rrule" : string,
   "enabled" : boolean,
   "eids" : [ ],
   "connector" : string,
   "connector_name" : string,
   "author" : string
}
```

## rrules

Dans le cas où la `rrule` est absente, `tstart` et `tstop` font office de plage d’activation du pbehavior, sans récurrence.

Dans le cas où la `rrule` est présente, `tstart` et `tstop` seront recalculés afin de refléter la récurrence.

## Création

API HTTP : voir l’API `/pbehavior/create`

Event de type pbehavior : créé à partir des champs cités en introduction
