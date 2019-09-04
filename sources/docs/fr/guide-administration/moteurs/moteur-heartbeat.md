# Heartbeat

Le moteur heartbeat permet de créer des alarmes en l'absence d'événements correspondant à un pattern donné durant un intervalle donné.

Les heartbeats sont définis dans la collection MongoDB `heartbeat`, et peuvent être ajoutés et modifiés avec l'[API Heartbeat](../../heartbeat/api_v2_heartbeat.md).

## Fonctionnement

La file du moteur est placée juste après l'exchange `canopsis.events`.

Un `heartbeat` est une règle qui définit un pattern d'événement entrant et un intervalle de temps.

Passé cet intervalle de temps, si le moteur n'a pas traité d'événements correspondant au pattern donné, il lèvera automatiquement une alarme.

L'alarme levée reste en cours jusqu'à l'arrivée d'un nouvel événement correspondant au pattern donné.

### Options de l'engine-heartbeat

```
  -d    debug
  -version
        version infos
```

#### Patterns

Les patterns acceptés sont une version simplifiée de ceux utilisés pour l'[event-filter](../../event-filter/index.md).

Les valeurs supportées sont uniquement de type `string` et la seule condition supportée est de type `equal` (pas d'expression régulière, de `not equal`, etc.)

#### Intervalles

Les intervalles de temps peuvent être définis en minutes (`m`) ou en heures (`h`).

#### Alarmes

L'événement généré pour créer une alarme prend la forme suivante :

```json
{
	"resource": "connector:heartbeat_test_1.connector_name:heartbeat_test_1_name",
	"event_type": "check",
	"component": "heartbeats",
	"connector": "heartbeat",
	"source_type": "resource",
	"state": 3,
	"connector_name": "heartbeat"
}
```

Seul le champ `resource` peut varier. Il est la concaténation des patterns de règles appliqués.

Ici : `"connector" : "heartbeat_test_1"` et `"connector_name" : "heartbeat_test_1_name"`.

## Collection

Les heartbeats sont stockés dans la collection Mongo `heartbeat` (voir [API Heartbeat](../../guide-developpement/action/api_v2_heartbeat.md) pour la création des heartbeats).

Un exemple de heartbeat pour générer une alarme si aucun événement avec le connecteur `heartbeat_test_1` et le nom de connecteur `heartbeat_test_1_name` n'ont été traités depuis plus d'une minute.

```json
{
    "_id" : "a9807256ef7e2a138a1caa204accc792",
    "pattern" : {
        "connector" : "heartbeat_test_1",
        "connector_name" : "heartbeat_test_1_name"
    },
    "expected_interval" : "1m"
}
```
