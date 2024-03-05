# Présentation de l'interface web de Canopsis

L'interface web de Canopsis est basée sur [VueJS](https://vuejs.org) et [Material Design](https://material.io).

Elle n'est officiellement compatible qu'avec une version raisonnablement récente des navigateurs Mozilla Firefox et Google Chrome. Voyez le document décrivant les [limitations de compatibilité avec les anciens navigateurs](../limitations/index.md) pour plus de détails.

## Les vues

Canopsis fonctionne avec un système de vues, personnalisables, regroupées dans des groupes de vues.  
Si vous avez les droits nécessaires, il vous est possible de créer/modifier des vues et des groupes de vues.  
Par ailleurs, il exite également un système de vues privées qui permet à un utilisateur de créer des vues qui ne seront visibles que de lui.  
Pour plus de détails concernant les vues, rendez-vous dans la section "Les vues", en cliquant [ici](./vues/index.md).

## Les widgets

Plus types de widgets sont disponibles dans Canopsis, pour en connaitre le fonctionnement, rendez-vous dans les sections "Les widgets", en cliquant [ici](./widgets/index.md).

## Thèmes

L'interface de Canopsis est livrée avec plusieurs thèmes de base qui ne sont pas modifiables par l'utilisateur : 

* `Canopsis` : Thème par defaut
* `Canopsis dark` : Mode sombre du thème par défaut
* `Color blind` : Thème adapté pour les daltoniens
* `Color blind dark` : Mode sombre du thème adapaté pour les daltoniens

Par ailleurs, Canopsis offre la possibilité de créer ses propres thèmes graphiques.  
Cette [page de documentation](./themes/index.md) vous décrit les différents paramètres.

## Les playlists

Les vues peuvent être affichées sous forme de playlists c'est-à-dire les unes après les autres avec un délai associé.  
Pour savoir comment utiliser cette fonctionnalité, rendez-vous sur cette [page explicative](../menu-administration/listes-de-lecture.md).

## Options de l'interface

Plusieurs options sont disponibles pour personnaliser l'interface de Canopsis.  
Pour les découvrir, cliquez [ici](../menu-administration/parametres.md)

## Helpers Handlebars

À différents endroits de l'interface de Canopsis, le texte affiché est personnalisable grâce au langage Handlebars.  
Différents *helpers* sont disponibles pour ajouter un peu de logique dans ces templates.  
Pour plus de détails, consultez la [documentation des helpers Handlebars de Canopsis](helpers/index.md).

## La diffusion de messages

Il est possible de programmer l'affichage d'un bandeau contenant des messages d'information qui apparaîtront dans l'interface de Canopsis.  
Reportez vous à cette [page](../menu-administration/diffusion-de-messages.md) pour en savoir plus.
