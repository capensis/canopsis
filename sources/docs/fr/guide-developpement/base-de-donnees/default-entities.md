# Entités

## Structure

Voici la structure d'une [entité](../../guide-utilisation/vocabulaire/index.md#entite).

```javascript
{
    "_id" :                // entity ID
    "name" :               // entity name
    "impact" : [
        ...                // entities impacted
    ],
    "depends" : [
        ...                // entities which depends
    ],
    "enable_history" : [
        ...
    ],
    "measurements" :
    "enabled" :            // entity is enabled or not
    "infos" : {
        ...                // enriched infos
    },
    "type" :               // entity type
}
```

## Context-graph

À la création et la mise-à-jour d'une entité, son contexte est graphé via le [context-graph](../../guide-utilisation/vocabulaire/index.md#context-graph).

Les champs `impact` et `depends` contiennent chacun un tableau constitué des `_id` d'entités auxquelles l'entité courante est liée.

### Évènement

Voici comment sont graphées les entités d'un [évènement](../../guide-utilisation/vocabulaire/index.md#evenement) :

### `connector`

- `impact` : `_id` de la `resource` (si l'évènement est de type `ressource`)
- `depends` : `_id` du `component`

### `component`

- `impact` : `_id` du `connector`
- `depends` : `_id` de la `resource` (si l'évènement est de type `ressource`)

### `resource`

- `impact` : `_id` du `component`
- `depends` : `_id` du `connector`

### Observateur

Un [Observateur](../../guide-utilisation/vocabulaire/index.md#observateur) (ou « watcher » ) est conçu pour que les `depends` de son `context-graph` incluent les entités correspondant à son ou ses patterns.

En retour, il est automatiquement ajouté aux `impact` des entités en question.

## `_id` des entités

### Composant

La valeur du champ `_id` du composant est directement celui du champ `component` présent dans l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement).

### Connecteur

La valeur du champ `_id` du connecteur est la concaténation des champs `connector` et `connector_name` présents dans l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement).

### Ressource

La valeur du champ `_id` de la ressource est la concaténation des champs `resource` et `component` présents dans l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement).

### Watcher

La valeur du champ `_id` du watcher est celle indiquée à l'ajout via l'[API watcherng](../api/api-v2-watcherng.md). En l'absence d'`_id` (comme lors de l'ajout via l'[explorateur de contexte](../../guide-utilisation/interface/widgets/contexte/index.md), un `_id` unique est généré automatiquement.

## Collection MongoDB

Les [entités](../../guide-utilisation/vocabulaire/index.md#entite) sont stockées dans la collection `default_entities` du MongoDB de Canopsis.

Elles ne sont pas purgées automatiquement.

Leur purge manuelle peut amener à l'existence d'alarmes sans entité. Il faut donc procéder avec précaution.
