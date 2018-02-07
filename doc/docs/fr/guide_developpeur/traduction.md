# Documentation de Canopsis pour les traductions

## Pour commencer

Les traductions dans Canopsis s'implémentent côté front et sont toujours locales à une brique.

Pour traduire un élément "non traduit" dans Canopsis, localiser la brique dans laquelle il se trouve puis effectuer la traduction dans celle-ci.

Un mot ou une expression à traduire doit être modifié dans le code (dans le fichier javascript ou dans le template) afin d'être traduit et sa traduction doit être spécifiée dans un dictionnaire prévu à cet effet.

## Modifications pour traduire une expression

### L'expression se trouve dans le template

par exemple, il faut traduire 'More info' présent dans ce template:

```html
    <a href="#">
        More info
    </a>
```

Ici, la modification sera d'ajouter un renderer appelé 'translator' comme ceci:

```html
    <a href="#">
        {{tr 'More info'}}
    </a>
```

Dans le cas où la chaîne est un attribut comme ici pour un placeholder:

```html
    {{input insert-newline='search' value=value class="form-control inlineInput" placeholder="Search" type="text"}}
```

La syntaxe pour le renderer sera celle-ci:

```html
    {{input insert-newline='search' value=value class="form-control inlineInput" placeholder=(tr 'Search') type="text"}}
```

### L'expression se trouve dans un fichier javascript

On utilisera ici la méthode loc() d'Ember. Par convention celle-ci est assignée à `__` eb début de fichier comme ici:

```javascript
    var get = Ember.get,
             set = Ember.set,
             moment = window.moment;
             __ = Ember.String.loc;
```

Ensuite la chaîne sera complétée de `__` comme dans cet exemple:

```javascript
    return __('Canceled by')
```

## Spécifications de la traduction

Pour traduire des expressions dans une brique frontend, plusieurs éléments sont requis.

- à la racine de la brique, créer un répertoire i18n
- Dans ce répertoire, créer les fichiers JSON de langues (en.json, fr.json). Ces derniers contiendront les dictionnaires de traduction.
Voici des exemples de dictionnaires à implémenter:

en.json
```javascript
    {
   "Acknowledged by ": "Acknowledged by ",
   "Ack removed by ": "Ack removed by ",
   "Ticket association by ": "Ticket association by ",
   "Ticket declared by ": "Ticket declared by ",
   "Canceled by ": "Canceled by ",
   "Comment by ": "Comment by ",
   "Restored by ": "Restored by ",
   "Status increased": "Status increased",
   "Status decreased": "Status decreased",
   "State increased": "State increased",
   "State decreased": "State decreased",
   "State changed": "State changed",
   "Snoozed by ": "Snoozed by ",
   "Cropped states (since last change of status)": "Cropped states (since last change of status)",
   "Hard limit reached !": "Hard limit reached !"
 }
```

fr.json
```javascript
    {
   "Acknowledged by ": "Acquitté par ",
   "Ack removed by ": "Acquittement supprimé par ",
   "Ticket association by ": "Association de ticket par ",
   "Ticket declared by ": "Ticket déclaré par ",
   "Canceled by ": "Annulé par ",
   "Comment by ": "Commenté par ",
   "Restored by ": "Restauré par ",
   "Status increased": "Status incrémenté",
   "Status decreased": "Status décrémenté",
   "State increased": "Etat incrémenté",
   "State decreased": "Etat décrémenté",
   "State changed": "Etat modifié",
   "Snoozed by ": "Reporté par ",
   "Cropped states (since last change of status)": "Etats croisés (depuis le dernier changement de statut)",
   "Hard limit reached !": "Limite physique atteinte !"
 }
```

- Créer un répertoire requirejs-modules à la racine de la brique s'il n'existe pas déjà.
- Créer un fichier i18n à l'intérieur de ce répertoire, contenant ceci:
```javascript
 /*
  * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
  *
  * This file is part of Canopsis.
  *
  * Canopsis is free software: you can redistribute it and/or modify
  * it under the terms of the GNU Affero General Public License as published by
  * the Free Software Foundation, either version 3 of the License, or
  * (at your option) any later version.
  *
  * Canopsis is distributed in the hope that it will be useful,
  * but WITHOUT ANY WARRANTY; without even the implied warranty of
  * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  * GNU Affero General Public License for more details.
  *
  * You should have received a copy of the GNU Affero General Public License
  * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
  */
 
 define([
     'text!canopsis/brick-nomdelabrique/i18n/' + i18n.lang + '.json'
 ], function (langFile) {
     var langFile = JSON.parse(langFile);
     var langKeys = Em.keys(langFile);
 
     for (var i = 0; i < langKeys.length; i++)
         Em.STRINGS[langKeys[i]] = langFile[langKeys[i]];
 });
```

Pensez bien à adapter le nomdelabrique avec la brique dans laquelle vous travaillez.

- Pour finir, `npm run full-compile` dans la brique pour prendre en compte toutes les modifications.