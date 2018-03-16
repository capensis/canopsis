# Typage

Vérifications sur les types des données en base.

## MongoDB

Documentation : https://docs.mongodb.com/v3.4/reference/operator/query/type/

Cas général de la requête permettant de récupérer les documents ne correspondant pas au typage attendu sur un champs :

```json
{
    $and: [
        {"<field>": {$exists: true}},
        {"<field>": {
            $not: {$type: <typeNumber or "alias">}
        }}
    ]
}
```

Exemple avec le champs `enabled` devant être de type `bool` :

```json
{
    $and: [
        {"enabled": {$exists: true}},
        {"enabled": {
            $not: {$type: 8}
        }}
    ]
}
```
