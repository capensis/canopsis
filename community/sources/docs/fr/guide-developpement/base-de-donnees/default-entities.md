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

À la création et la mise à jour d'une entité, son contexte est graphé via le [context-graph](../../guide-utilisation/vocabulaire/index.md#context-graph).

Les champs `impact` et `depends` contiennent chacun un tableau constitué des `_id` d'entités auxquelles l'entité courante est liée.

### Entité

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

## `_id` des entités

### Composant

La valeur du champ `_id` du composant est directement celui du champ `component` présent dans l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement).

### Connecteur

La valeur du champ `_id` du connecteur est la concaténation des champs `connector` et `connector_name` présents dans l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement).

### Ressource

La valeur du champ `_id` de la ressource est la concaténation des champs `resource` et `component` présents dans l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement).

## Collection MongoDB

Les [entités](../../guide-utilisation/vocabulaire/index.md#entite) sont stockées dans la collection `default_entities` du MongoDB de Canopsis.

Elles ne sont pas purgées automatiquement.

Leur purge manuelle peut amener à l'existence d'alarmes sans entité. Il faut donc procéder avec précaution.
