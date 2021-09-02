# Moteur `engine-heartbeat` (Community)

Le moteur `engine-heartbeat` permet de créer des alarmes, si des événements suivant un motif donné n'ont pas été rencontrés depuis un intervalle donné.

## Utilisation

### Options du moteur

La commande `engine-heartbeat -help` liste toutes les options acceptées par le moteur.

## Fonctionnement

Un heartbeat est une règle qui définit un modèle d'événement entrant et un intervalle de temps.

Passé cet intervalle de temps, si le moteur n'a pas traité d'événements correspondant au modèle donné, il lèvera automatiquement une alarme.

L'alarme levée reste en cours jusqu'à l'arrivée d'un nouvel événement correspondant au modèle donné.

#### Modèles d'évènements (patterns)

Les modèles d'évènements acceptés (ou *patterns*) sont une version simplifiée de ceux utilisés pour l'[event-filter](moteur-che-event_filter.md).

Les valeurs acceptées sont uniquement de type `string` et la seule condition acceptée est de type `equal` (pas d'expression régulière, de `not equal`, etc.).

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

## Collection MongoDB associée

Les heartbeats sont stockés dans la collection MongoDB `heartbeat`.

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
