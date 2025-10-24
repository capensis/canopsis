# Météo des services

![Météo de services](./img/weather.png  "Météo de services")

## Sommaire

### Guide utilisateur

1. [Présentation générale](#presentation-generale)
2. [Les tuiles](#les-tuiles)
3. [La modale](#la-modale)

### Guide exploitant

1. [Aide sur les variables](#aide-variables)
2. [Paramètres du widget](#parametres-du-widget)

## Guide utilisateur

### Présentation générale

### Les tuiles

La météo de services est composée de tuiles.

Exemple d'une tuile :

![Exemple d'une tuile - Météo de services](./img/tuile-weather.png  "Exemple d'une tuile - Météo de services")

Chaque tuile correspond à un service.

Le contenu de texte de cette tuile est personnalisable (*Cf: [Guide exploitant](#guide-exploitant_1)*). Il permet de présenter des informations sur le service.

La couleur de la tuile et une icône présente sur celle-ci permettent d'obtenir des informations sur **la criticité** du service.

#### Bac à alarmes

Sur le partie inférieuré des tuiles, une indication "See alarms" est cliquable et permet un accès direct à la liste des alarmes liées au service de la tuile.

#### La couleur

Il est possible de définir la couleur des tuiles par criticité d'alarme (par défaut) ou par priorité de service.

![](./img/color_indicator.png)

##### Par criticité

La couleur de la tuile correspond à la criticité du service. Elle est calculée en prenant en compte la pire criticité parmi les entités surveillées par ce service.

> **Exemple**
>
> Un service surveille deux entités, A et B:
> 
> - A a une criticité de 1
> - B a une criticité de 3
>
> La criticité du service sera alors égale à 3.

- Vert: Criticité = 0 (Ok)
- Jaune: Criticité = 1 (Mineur)
- Orange: Criticité = 2 (Majeur)
- Rouge: Criticité = 3 (Critique)
- Gris: Le service (ou toutes les entités du service) possède un comportement périodique actif (pause, maintenance…).

##### Par priorité

La couleur de la tuile correspond à l'état d'impact du service. Cette valeur est le produit de l'alarme la plus critique parmi les entités surveillées par ce service et du niveau d'impact défini de ce dernier.

> **Exemple**
>
> Un service surveille deux entités, A et B:
> 
> - A a un criticité de 1 
> - B a une criticité de 3
> - Et le niveau d'impact du service est 5
>
> L'état d'impact du service est donc de `3 * 5 = 15`

Voici la palette de couleurs correspondant à l'état d'impact :

![](./img/color_prio_table.png)

#### L'échelle

Chaque tuile comprend une petite échelle d'état du service.

![](./img/tuile_scale.png)

Cette échelle permet situer l'état du service sur la palette de couleurs ainsi que sa valeur numérique (0 à 30).

#### Le point d'interrogation

Cette icône cliquable permet l'affichage de la valeur de toutes les variables du service. 

Ces variables peuvent être utilisées dans l'édition des [templates](#template-tuile).

#### L'icône

- **Soleil**: Le service possède une criticité "Ok" (égale à 0)
- **Personne**: Le service possède une criticité "Mineure" (égale à 1) ou "Majeure" (égale à 2)
- **Nuage**: Le service possède une criticité "Critique" (égale à 3)
- **Clé**: Le service possède un comportement périodique actif, de type "Maintenance"
- **Lune**: Le service possède un comportement périodique actif, de type "Hors plage de surveillance"
- **Pause**: Le service ne possède pas de comportement périodique, mais toutes les entités liées à ce service possèdent un comportement périodique actif.

#### Le clignotement

Une tuile de la météo de service clignotera si une des entités lui appartenant possède une alarme non acquittée, et que celui-ci n'est pas en pause ou ne possède pas d'entité en pause.

### La modale

Au clic sur une tuile de la météo de services, une fenêtre s'ouvre.

Le contenu de cette fenêtre est configurable depuis les paramètres du widget.

Celle-ci contient la liste des entités surveillées. Au clic sur l'une de ces entités, deux onglets apparaissent :

- "Info" : Affiche les informations configurées dans le template des entités qui se trouve dans les paramètres avancés du widget. Ainsi que la listes des actions possible.

![](./img/modale_info.png)

- "Tree of dependencies" : Affiche l'arbre de dépendences de l'entité sélectionnée.

![](./img/modale_tree.png)

### Les actions

Dans la liste des entités affichée, des actions sont disponibles sur chacune d'entre elles. Les actions disponibles dépendent de la criticité de l'entité.

Au clic sur les icônes d'actions, celles-ci sont mises en attente. Elles ne sont exécutées qu'au clic sur le bouton `Save change`.

- ![Action: Déclarer un incident](./img/action_declareTicket.png "Action: Déclarer un incident") *Déclarer un incident*: Cette action vous permet de déclarer un numéro de ticket, associé à un incident. Au clic sur cette action, une fenêtre s'ouvre, vous permettant d'indiquer un numéro de ticket. Cette action déclenche également automatiquement une action d'acquittement.
- ![Action: Pause](./img/action_pause.png "Action: Pause") *Pause*: Cette action vous permet de mettre une entité en pause. Au clic, une fenêtre s'ouvre. Celle-ci vous permet de renseigner un commentaire, ainsi que la raison de la pause. Cette action n'est disponible que pour les entités qui ne sont pas déjà en pause.
- ![Action: Play](./img/action_play.png "Action: Play") *Play*: Cette action vous permet de modifier tout les comportements périodiques de type `Pause`. La date de fin de ces comportements périodiques est modifiée pour passer à la date actuelle, ce qui met de fait fin à la pause. Cette action n'est disponible que pour les entités en pause.
- ![Action: acquittement](./img/action_ack.png "Action: acquittement") *Acquittement*: Cette action vous permet d'acquitter une alarme présente sur une entité. Cette action n'est disponible que pour les entités ayant une criticité différente de "Ok" (0), et ayant une alarme non acquittée.
- ![Action: Validate](./img/action_validate.png "Action: Validate") *Valider*: Cette action déclenche un changement de criticité de l'alarme, de majeure (2) à critique (3). Elle entraîne également automatiquement une action d'acquittement. Celle-ci n'est disponible que pour les entités ayant une criticité majeure (2).
- ![Action: Invalidate](./img/action_invalidate.png "Action: Invalidate") *Invalider*: Cette action déclenche une action d'annulation de l'alarme. Elle entraîne également automatiquement une action d'acquittement. Celle-ci n'est disponible que pour les entités ayant un criticité majeure (2).

## Guide exploitant

### Aide - Variables

Durant la configuration de votre widget Météo de services, notamment les Templates, il vous sera possible d'accéder à des variables concernant les services.

> **Exemple**
>
> Il vous sera possible d'afficher, pour chacune des tuiles de la météo de services, le nom du service, ou son identifiant, etc.

Afin de connaitre les variables disponibles, cliquer sur [le point d'interrogation](#le-point-dinterrogation) d'une tuile.

### Paramètres du widget

1. [Taille du widget](#taille-du-widget-requis)
2. [Titre](#titre-optionnel)
3. [Éditeur de filtre](#editeur-de-filtre-optionnel)
4. [Paramètres avancés](#parametres-avances)
    1. [Colonne de tri par défaut](#colonne-de-tri-par-defaut)
    2. [Template - Tuiles](#template-tuile)
    3. [Template - modale](#template-modale)
    4. [Template - Entités](#template-entites)
    5. [Colonnes - Petit](#colonnes-petit)
    6. [Colonnes - Moyen](#colonnes-moyen)
    7. [Colonnes - Large](#colonnes-large)
    8. [Marges](#marges)
    9. [Hauteur](#hauteur)
    10. [Type de modale](#type-de-modale)

#### Taille du widget (*requis*)

Ce paramètre permet de régler la taille du widget.

![Paramètre Taille du widget](../img/settings/widget-size.png "Paramètre Taille du widget")

La première information à renseigner est la ligne dans laquelle le widget doit apparaitre. Ce champ permet de rechercher parmi les lignes disponibles. Si aucune ligne n'est disponible, ou pour en créer une nouvelle, entrez son nom, puis appuyez sur la touche Entrée.

Ensuite, les 3 champs en dessous permettent de définir respectivement la largeur occupée par le widget sur mobile, tablette, de ordinateur de bureau.
La largeur maximale est de 12 colonnes pour un widget, la largeur minimale est de 3 colonnes.

#### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

#### Editeur de filtre (*optionnel*)

Ce paramètre permet de définir le filtre à appliquer à la météo de services.
Ce filtre permet de n'afficher qu'une partie des services.
Pour plus de détails sur les filtres et leur création, voir la partie sur [Les filtres](../../filtres/index.md).

Pour créer un filtre, ou éditer celui actuellement actif, cliquez sur le bouton 'Créer/Editer'. Une fenêtre de création de filtre s'ouvre alors.

Pour supprimer le filtre actuellement actif, cliquez sur l'icone de suppression se trouvant à droite du bouton 'Créer/Editer'. Une fenêtre vous demande alors de confirmer la suppression.

!!! warning "Champs utilisables dans le filtre"
    Le filtre utilise les champs des entités (qui sont différents des champs utilisables dans les templates). Par exemple, pour filtrer sur le nom d'un service, il faut utiliser `name`, et non `display_name`.

#### Paramètres avancés

##### Colonne de tri par défaut

Ce paramètre permet de trier les tuiles selon un attribut pré-défini par ordre alphabétique.  

!!! attention
    Le tri implémenté est sensible à la casse et fait que les majuscules sont traitées avant les minuscules.

Par défaut, les attributs disponibles pour le tri sont :

* `name` 
* `state`

Vous avez la possibilité d'utiliser le critère de votre choix en écrivant directement dans la configuration l'attribut de tri souhaité.  

Exemple : pour faire un tri selon la valeur du champ enrichi `mon_attribut` ajouté depuis l'[explorateur de contexte](../contexte/index.md), remplissez le formulaire comme suit : 

![Tri par défaut](./img/tri.png)

##### Template - Tuile

Ce paramètre permet de personaliser les informations affichées à l'intérieur des tuiles de la météo de service.

Le langage utilisé ici est le Handlebars.

Cliquez sur le bouton 'Afficher/Editer'. Une fenêtre s'ouvre avec un éditeur de texte. Entre le texte souhaité pour le template des tuiles, puis cliquez sur 'Envoyer'.

Une variable est disponible ici pour vous permettre d'afficher les détails du service : `entity`.
Exemple : Pour afficher le champ `display_name` du service (qui correspond au nom du service), il vous faut écrire dans le template : `{{ entity.display_name }}`.
Tous les champs disponibles dans le service sont disponibles ici.

##### Template - Modale

Ce paramètre permet de personnaliser les informations affichées dans la fenêtre 'Plus d'infos' (ouverte au clic sur 'Plus d'infos', sur une des tuiles de la météo de services).

Il vous est possible ici d'afficher, à n'importe quel endroit de la modale, la liste des entités concernées par le service sur lequel vous avez cliqué. Pour ce faire, insérez dans le template:

```
{{ entities }}
```

Cela aura pour effet d'insérer dans la modale la liste des entités. Par défaut, le nom de l'entité sera affiché pour chacune d'entre elles. Il vous est possible de modifier la valeur affichée ici. Tous les champs de l'entité sont disponibles. Pour ce faire, ajoutez un argument ```name``` à la balise précédemment ajoutée. Il vous est donc possible d'écrire, par exemple :

```
{{ entities name="entity.entity_id" }}
```

Pour chaque entité de la liste, l'id de l'entité sera affiché, à la place de son nom.

On peut également entrer :

```
{{ entities name="entity.infos.customer.value" }}
```

Pour chaque entité de la liste, la valeur de leur champ enrichi customer sera affiché, à la place du nom.

Celui-ci fonctionne de la même manière que le paramètre Template - Tuile présenté ci-dessus. Cliquez [ici](#template-tuile) pour vous rendre à cette partie.

##### Template - Entités

Ce paramètre permet de personnaliser les informations affichées pour chaque entités dans la fenêtre 'Plus d'infos' (ouverte au clic sur 'Plus d'infos', sur une des tuiles de la météo de services).

**Attention: La liste des entités n'est affichée que si cela a été précisé dans le [Template - Modale](#template-modale).**

Le langage utilisé ici est le Handlebars.

Cliquez sur le bouton 'Afficher/Editer'. Une fenêtre s'ouvre avec un éditeur de texte. Entre le texte souhaité pour le template des tuiles, puis cliquez sur 'Envoyer'.

Une variable est disponible ici pour vous permettre d'afficher les détails de l'entité : `entity`.
Exemple : Pour afficher le champ 'name' de l'entité (qui correspond au nom de l'entité), il vous faut écrire dans le template : `{{ entity.name }}`.
Tous les champs disponibles dans l'entité sont disponibles ici.

##### Colonnes - Petit

Ce paramètre permet de définir la proportion de l'écran, en largeur, prise par chaque tuile de la météo de services. Ce paramètre concerne les écrans de mobiles (largeur < 450px). Une tuile occupe au minimum une colonne (1/12 de la largeur de la page), et au maximum 12 colonnes (100 % de la largeur de la page).

Il suffit de faire glisser le curseur pour sélectionner le nombre de colonne par tuile souhaité.

##### Colonnes - Moyen

Ce paramètre permet de définir la proportion de l'écran, en largeur, prise par chaque tuile de la météo de services. Ce paramètre concerne les écrans de tablettes (largeur < 900px). Une tuile occupe au minimum une colonne (1/12 de la largeur de la page), et au maximum 12 colonnes (100 % de la largeur de la page).

Il suffit de faire glisser le curseur pour sélectionner le nombre de colonne par tuile souhaité.

##### Colonnes - Large

Ce paramètre permet de définir la proportion de l'écran, en largeur, prise par chaque tuile de la météo de services. Ce paramètre concerne les écrans d'ordinateurs (largeur > 900px). Une tuile occupe au minimum une colonne (1/12 de la largeur de la page), et au maximum 12 colonnes (100 % de la largeur de la page).

Il suffit de faire glisser le curseur pour sélectionner le nombre de colonne par tuile souhaité.

##### Marges

Ce paramètre permet de régler les espaces séparant les tuiles de la Météo de services.

Celui-ci est séparé en quatre, vous permettant de régler l'espace que vous souhaitez pour chaque côté des tuiles (haut, bas, droite et gauche).

Pour modifier ce paramètre, faites glisser le sélecteur, afin de choisir une valeur entre 0 et 5 (0 correspondant à l'absence de marge, 5 le maximum de marge).

Par défaut, ce paramètre est réglé sur une valeur de 1 pour chacuns des côtés des tuiles.

##### Hauteur

Ce paramètre permet de régler la hauteur des tuiles de la Météo de services.

Pour le modifier, faites glisser le sélecteur, afin de choisir une valeur entre 1 (hauteur minimale) et 20 (hauteur maximale).

Par défaut, ce paramètre est réglé sur une valeur de 1.

