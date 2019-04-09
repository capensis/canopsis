# Structure du projet

## Général

![Structure dossier 1er niveau](./img/folder_structure_1.png  "Structure dossier 1er niveau")

### Les dossiers

#### /dist

Ce dossier contient les fichiers compilés.
Il peut ne pas exister lors d'une nouvelle installation.

Ce dossier n'est pas utilisé au moment du développement.

#### /docker

Ce dossier contient le DockerFile permettant de construire une image Docker de l'interface.

Ce dossier n'est pas utilisé au moment du développement.

#### /node_modules

Ce dossier contient toutes les dépendances du projet.

Celui-ci n'est jamais modifié directement durant le développement.

#### /public

Ce dossier contient le fichier html 'racine' du projet, ainsi que les polices utilisées.

#### /src

Ce dossier contient la quasi totalité du code du projet.

C'est le dossier dans lequel se fera le développement.

#### /tests

Ce dossier contient les tests unitaires ainsi que les tests fonctionnels.


### Les fichiers de configuration

#### Générale

#### Environnement

#### Gestion des dépendances

#### Docker-compose

## Le dossier 'src'

Comme vu précedemment, le dossier ```src``` contient le code du projet.

### main.js

Ce fichier créé une nouvelle instance ```Vue``` en lui intégrant les différentes options et librairies, puis lance l'application en montant le composant racine (```app.vue```).

### app.vue

Ce fichier contient le composant racine de l'application.

Il intégre les barres latérales, les modales, les popups, ainsi que le composant nécessaire, en fonction du routage.

### bootstrap.js

Ce fichier tire parti de la librairie ```webfontloader``` pour charger les polices utilisées dans le projet.

### config.js

Ce fichier contient les variables de configuration de l'application.

Exemple: 

- ```API_ROUTES```: Contient les routes utilisées pour effectuer des requêtes à Canopsis.
- ```DEFAULT_LOCALE```: Langage par défaut à utiliser dans l'interface

### constants.js

Ce fichier contient toutes les constantes utilisées dans le projet.

### event-bus.js

Ce fichier n'est plus utilisé et devra être supprimé prochainement.

### router.js

Ce fichier contient la configuration du routage du projet.

Une variable ```routes``` est d'abord définie. Elle contient la liste des routes disponibles.

Puis une nouvelle instance de ```Router``` (de la librairie ```vue-router```) est créée. Celle-ci utilise les routes définies au dessus.

### /assets

Ce dossier contient les feuilles de style globales du projet, ainsi que les images utilisées.

### /components

Ce dossier contient les composants du projet.

Il est divisé en sous-dossier permettant de trier les composants en fonction de leur utilisation.

### /filters

Ce dossier contient les différents filtres utilisés dans les templates.

Ces filtres permettent d'appliquer une fonction Javascript sur une/des valeur(s), dans un template.

### /helpers

Ce dossier contient les différents helpers utilisés dans les templates.

Les helpers sont des fonctions Javascript, destinées à êtres utilisées à plusieurs reprises dans le code.

### /i18n

Ce dossier contient les traductions.

Le fichier ```index.js``` ajoute ces traductions à l'instance courante de ```Vue```.

Les traductions elles-mêmes se situent dans le dossier ```messages```. Chaque fichier correspond ici à une langue.

### /mixins

Ce dossier contient les différents mixins.

Les mixins permettent de réutiliser des fonctionnalités dans plusieurs composants.

### /services

### /store

### /views
