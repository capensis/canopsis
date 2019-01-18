# Tableau de statistiques

// INSERER IMAGE TABLEAU DE STATS

## Sommaire
### Guide utilisateur

1. [Présentation du widget](#presentation-du-widget)

### Guide exploitant

1. [Paramètres du widget](#parametres-du-widget)

## Guide utilisateur

### Présentation du widget

Le widget tableau de statistiques se présente sous la formce d'un tableeau. Chaque ligne de ce tableau correspond à une entité, et aux statistiques associées à cette entité. La première colonne est non configurable. Celle-ci présente, pour chaque ligne, le nom de l'entité. Chacune des colonnes qui suit correpond à une statistique.

Le tableau présente entre 1 et n statistiques. Les statistiques affichées sont configurées depuis le panneau de paramètres du widget (voir [paramètres du widget](#parametres-du-widget)).

Une tendance par rapport à la période précédente peut être affichée à droite de la valeur de la statistique (voir [paramètres du widget](#parametres-du-widget)).

## Guide exploitant

### Paramètres du widget

1. Taille du widget (*requis*)
2. Titre (*optionnel*)
3. Durée (*requis*)
4. Date de fin (*requis*)
5. Sélecteur de statistique (*requis*)
6. Paramètres avancés
    1. Editeur de filtre (*optionnel*)

#### Taille du widget (*requis*)

Ce paramètre permet de régler la taille du widget.

![Paramètre Taille du widget](../../img/settings/widget-size.png "Paramètre Taille du widget")

La première information à renseigner est la ligne dans laquelle le widget doit apparaître. Ce champ permet de rechercher parmi les lignes disponibles. Si aucune ligne n'est disponible, ou pour en créer une nouvelle, entrez son nom, puis appuyez sur la touche Entrée.

Ensuite, les 3 champs en dessous permettent de définir respectivement la largeur occupée par le widget sur mobile, tablette, de ordinateur de bureau.
La largeur maximale est de 12 colonnes pour un widget, la largeur minimale est de 3 colonnes.

#### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

#### Durée (*requis*)

Ce paramètre permet de définir la durée à prendre en compte pour le calcul des statistiques.

Exemple: Une durée de 2 mois permet de calculer les statistiques concernant les 2 mois précédents la [date de fin](#date-de-fin-requis).

Valeur par défaut: 1 jour.

#### Date de fin (*requis*)

Ce paramètre permet de définir la date de fin de calcul des statistiques pour ce widget.

Exemple: Si la [durée](#duree-requis) a été paramètrée à 2 mois, et la date de fin au 01 Septembre 2018, les statistiques affichées correspondront aux valeurs pour la période allant du 01 Juillet 2018 au 01 Septembre 2018.

Les dates disponibles dépendent de l'unité de [durée](#duree-requis) sélectionnée.

- Pour une durée en mois, la date de fin ne peut être que le premier jours d'un mois, à 00:00 GMT.
- Pour toutes les autres unités de durée, la date de fin doit correspondre à une heure pleine. Ex: 17/01/2019 18:00.

#### Sélecteur de statistiques (*requis*)

Ce paramètre permet de définir les statistiques à afficher dans le tableau. Chaque statistique se verra affecter une colonne du tableau.

**Il est obligatoire d'ajouter au moins une statistique**

Pour ajouter une statistique, cliquez sur le bouton ```Ajouter une statistique```.

Une fenêtre s'ouvre.

// INSERER IMAGE MODAL AJOUT DE STAT

Cette fenêtre vous permet de définir la statistique souhaitée.

- Statistique à afficher (voir [liste des statstiques disponibles](../index.md#les-statistiques-disponibles)).
- Titre associé à cette statistique.
- Tendance: Permet de définir si vous souhaitez récupérer et afficher la tendance par rapport à la période précédente, pour chaque valeur.
- Options: Liste d'options concernant la statistique sélectionnée. Les options varient selon la type de statistique voulue :
    - ```Recursif```: Si l'option est activée, permet de calculer la statistique sur l'entité, ainsi que sur ses dépendances, et les dépendances de ses dépendances, etc...
    - ```Etats```: Permet de ne prendre en compte que les alarmes avec le/les état(s) (ok, mineure, majeure ou critique) sélectionné(s).
    - ```Auteurs```: Permet de ne prendre en compte que les alarmes dont le/les auteur(s) fait parti de la liste précisée ici. Pour ajouter un auteur à la liste, entrez son nom, puis appuyer sur la touche "Entrée".
    - ```Sla```: Permet de préciser le temps définit dans le SLA. **Attention: Ce paramètre est requis pour le calcul des statistiques "Taux d'Ack conforme SLA" et "Taux de résolution conforme Sla"**.

Cliquez sur le bouton ```Envoyer``` pour ajouter cette statistique.

La liste des statistiques ajoutées au widget est visible depuis le panneau de paramètres du widget. Un bouton vous permet ici d'éditer la statistique, ou de la supprimer de la liste.

// INSERER IMAGE LISTE STATS

#### Paramètres avancés

##### Editeur de filtres (*optionnel*)

Ce paramètre vous permet de définir le filtre à appliquer à la sélection d'entité. Il permet de ne sélectionner que les entités pour lesquels on souhaite afficher les statistiques.

Pour créer un filtre, ou éditer le filtre deja présent, cliquez sur le bouton ```Créer/Editer```.
Pour supprimer le filtre deja existant, cliquez sur le bouton situé à droite du bouton d'édition/création.

Au clic sur le bouton ```Créer/Editer```, une fenêtre d'édition de filtre s'ouvre. Une fois le nom du filtre, et le filtre lui-même renseignés, cliquez sur le bouton ```Envoyer``` pour le sauvegarder.

Pour plus de détails sur les filtres, et l'édition de filtres, cliquez [ici](../../../filtres).
