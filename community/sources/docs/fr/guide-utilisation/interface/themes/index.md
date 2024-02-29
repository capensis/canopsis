# Thèmes graphique

## Gestion des thèmes

Dans Canopsis, le menu `Thèmes` est dédié à la gestion des thèmes.  
![Menu thèmes](./img/themes_menu.png)

La liste des thèmes s'affiche alors et présente les thèmes par défaut qui sont livrés nativement dans Canopsis

![liste](./img/themes_liste.png)

## Création d'un thème

Le bouton "+" vous permet d'afficher le formulaire de création des thèmes. 

![bouton création](./img/themes_bouton_creation.png)

Les éléments graphiques personnalisables dans Canopsis sont regroupés dans 4 catégories

1. Principaux éléments de l'interface graphique
1. Paramètres de taille de police
1. Paramètres du bac à alarmes
1. Couleurs de criticités

### Principaux éléments de l'interface graphique

![parametres1](./img/themes_parametres1.png)

### Paramètres de taille de police

![parametres2](./img/themes_parametres2.png)

### Paramètres du bac à alarmes

![parametres3](./img/themes_parametres3.png)

### Couleurs de criticités

![parametres4](./img/themes_parametres4.png)


## Variables Handlebars associées

Il est possible de conditionner certains affichages en fonction du thème appliqué.  
Pour cela, vous disposez de la variable `theme` et notamment `theme.name` que vous pouvez utiliser comme suit : 

```js
{{#if theme.dark}}
<p style="background:black;color:white">Thème  dark activé</p>
{{else}}
<p style="background:white;color:black">Thème light activé</p>
{{/if}}

{{#compare theme.name '==' 'Canopsis'}}<p>Le thème actif est <span style="background-color:white;color:black">Canopsis</span></p>{{/compare}}
{{#compare theme.name '==' 'Canopsis dark'}}<p>Le thème actif est <span style="background-color:black;color:white">Canopsis dark</span></p>{{/compare}}
```
