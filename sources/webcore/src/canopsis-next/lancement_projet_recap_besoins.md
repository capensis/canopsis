# Refonte du Front-office de Canopsis - récapitulatif des besoins



## Fonctionnalités de Canopsis à conserver :


### Vues principales

- gestion d'alarmes sur le SI (visualisation et traitement)
	- Bac à Alarmes
	- Calendrier (reporting)


- surveillance de l'état du SI
	- pannes
		- Bac à alarmes
		- Calendrier
		- Topo/carto

	- métriques
		- Listes
		- Graphs
	- planification de maintenance 
		- Calendrier

- détection de panne //Root cause analysis
		- Topologie
		- Bac à alarmes
		- Service weather


- gestion de Canopsis
	- configuration des moteurs
	- gestion des utilisateurs/roles/droits
	- configuration de l'interface(gestion des vues)


### Widgets

- widgets : 
	-graphes : 
		- circulaires
		- diagramme en baton
		- courbes
	- service weather
	- list alarms
	- topologie (à remplacer)
	- cartographie (à créer)


- mixins : 

- editors :



## Fonctionnalités de canopsis à supprimer/non requises

- météo v1

## Contraintes


### Performances

- chargement de n'importe quelle page : < 3 secondes
- chargement initial : < 5 sec



### Interface 

- interface responsive :

	- mobile
	- tablette
	- écran de PC

- l'utilisateur n'a pas à configurer les détails techniques de ses vues
	- exemple : il n'a pas a configurer le type d'éditeur pour une vue
	- les vues ont une configuration par défaut

- layout : simplifier la création de pages et leur mise en forme :
	- faire des grilles de taille prédéfinie 
	- glisser/déposer des widgets pour les réorganiser

- rafraichissement automatique : 
	- l'utilisateur ne doit jamais avoir à rafraichir la page
	- toutes les vues doivent se rafraichir automatiquement (push/long polling)



## Framework de Front 

### Features

- rapide à l'exécution
- routage
- réutilisation de composants
- gestion des données 
- i18n
- templating
	- pouvoir exécuter du js raw
 - mixins/composants réutilisables
- support websockets/push http/2

### Build / processes

- système de build moderne (webpack ?)
- norme de développement
- process de développement
- tests unitaires/fonctionnels : 
	- du framework 
	- intégrables dans l'app


### Exigences non fonctionnelles

- framework connu : facile de trouver des développeurs
	 - Ember 2 (sous conditions)
	 - Angular 2
	 - React
	 - Vue

 - documentation :
 	- qualité
 	- quantité

 - communauté : 
 	- stack overflow dédié? 
 	- reddit dédié ?
 		- voir les nouveaux contenus, la taille de la communauté
 		- Soutenu/dévelopé par une compagnie vs 1 Communauté ?


## Todo : 

- maquette de l'UI
- specs des widgets / maquettes 
- comparaison des frameworks
- choix du framework JS
- choix du framework CSS
- définition du style de l'UI
- planning
