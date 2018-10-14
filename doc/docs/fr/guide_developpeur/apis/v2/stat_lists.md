# API Listes Statistiques

## Route

#### URL

`POST /api/v2/stats/lists/<nom de la liste>`

#### Paramètres

Cette route accepte en requête un objet JSON contenant les paramètres
suivants :

 - `tstop` : un timestamp indiquant le fin de la période pour laquelle la
   liste doit être calculées. Ce timestamp doit correspondre à une heure
   pile (e.g. 12:00, et non 12:03).
 - `duration` : la durée de la période, représentée par une chaîne
   `"<n><unité>"`, avec `<n>` un entier et `<unité>` une unité de temps (`h`,
   `d`, `w` ou `m`).
 - `mfilter` : un filtre mongodb, filtrant les entités pour lesquelles les
   listes doivent être calculées.
 - `parameters` : un objet contenant les paramètres spécifiques à la
   liste calculée. Ces paramètres sont précisés dans la documentation de
   chacune des listes.


#### Réponse

En cas de succès, la réponse est un tableau d'objets JSON contenant :

 - un champ `values` contenant les valeurs de la liste.
 - un champ `entity` contenant l'entité pour laquelle la liste a été
   calculée.


## Listes

### Intervalles passés dans un état

La liste `state_intervals` renvoie une liste contenant les intervalles de temps
passés dans chaque état par les entités.

#### Paramètres

Cette statistique accepte les paramètres suivants (à indiquer dans le champ
`parameters` d'une requête).

 - `states` (optionnel) : Un tableau contenant les états à prendre en compte
   (par exemple `[3]` pour ne renvoyer que les intervalles pendant lesquelles
   l'entité était en état critique).

### Exemple

```javascript
POST /api/v2/stats/lists/state_intervals
{
	"mfilter": {
		"type": "resource",
	},
    "tstop": 1534716000,
	"duration": "1h"
}
```

```javascript
[
    {
        "values": [
            {
                "duration": 969,
                "start": 1536242399,
                "state": 3,
                "end": 1536243119
            },
            {
                "duration": 326,
                "start": 1536243368,
                "state": 0,
                "end": 1536243694
            },
            {
                "duration": 23,
                "start": 1536243694,
                "state": 2,
                "end": 1536243708
            },
            {
                "duration": 2282,
                "start": 1536243717,
                "state": 0,
                "end": 1536245999
            }
        ],
        "entity": {
            "_id": "resource1/component1",
            "type": "resource"
            "impact": [
                "service"
            ],
        }
    },
    // ...
]
```
