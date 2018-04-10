# Comment migrer un bac à alarme existant chez un client vers une version permettant de customiser le nom des colonnes

## Contrainte

La principale contrainte réside dans le fait que pour l'ajout de cette nouvelle fonctionnalité, le schéma de description du widget a dû être modifié.

La fonctionnalité ajoute un nouvel attribut dans le schema du widget qui ressemble à ceci:

```javascript
"widget_columns": [
    {
        "value": "resource",
        "label": "resourcessss"
    },
    {
        "value": "component",
        "label": "componentssss"
    },
    {
        "value": "extra_details",
        "label": "suivi des actions"
    }
],
```

## Procédure

### Repérer la vue

- Se connecter à l'UI de Canopsis
- Repérer l'ID de la vue qui contient le bac à alarmes à migrer (visible dans l'URL), si la vue en question se trouve être la page d'accueil, retrouver le nom de cette view dans 'Settings -> User Interfaces' dans la barre de gauche dans l'interface, ici se trouvera le nom de la vue d'accueil.
- Rechercher cette vue dans MongoDB: ``` db.object.find({'_id': 'id trouvé dans lurl'}) ou bien db.object.find({'crecord_name': 'nom trouvé dans linterface'})```

### Editer le widget courant

- Copier la vue que nous avons trouvé précédemment dans un éditeur
- Y ajouter le champ cité au début avec les valeurs citées, **le champ doit être insérer dans l'attribut widget, au même niveau que l'attribut "columns"**
- Supprimer la view actuelle dans Mongo et y insérer à la place celle que nous avons modifié dans l'éditeur

### Mettre à jour le code

- Depuis les sources de la brick-listalarm: ``` git checkout pull && git checkout develop ```
- ``` rm -rf /opt/canopsis/var/www/src/canopsis/brick-listalarm && cp /path/to/sources/brick-listalarm /opt/canopsis/var/www/src/canopsis```


## Résultat

Pour observer les changements, opérer un F5 sur l'UI avec le cache désactivé.