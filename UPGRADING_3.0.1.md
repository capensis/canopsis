# Mise à jour Canopsis 3.0.1

## Mongo

Avec le retrait de l'ancien code lié au linklist, il faut supprimer toutes les tâches de ce type qui ont paramétrées.

Dans la console mongo, fait un :
```bash
db.getCollection('object').remove({'task':'tasklinklist'})
```
