# Histogramme de statistiques

// INSERER IMAGE HISTOGRAMME

## Sommaire
### Guide utilisateur

1. [Présentation du widget](#presentation-du-widget)

### Guide exploitant

1. [Paramètres du widget](#parametres-du-widget)

## Guide utilisateur

### Présentation du widget

Le widget Histogramme de statistiques vous permet de représenter un ensemble de statistiques, sous forme d'histogramme.

Il vous permet d'afficher les statistiques souhaitées pour un ou plusieurs groupes, représentants tout ou partie des entités présentes dans Canopsis.

## Guide exploitant

### Paramètres du widget

1. Taille du widget (*requis*)
2. Titre (*optionnel*)
3. Date de fin (*requis*)
4. Durée (*requis*)
5. Groupes de statistiques (*requis*)
6. Sélecteur de statistique (*requis*)
7. Paramètres avancés
    1. Couleurs des statistiques (*optionnel*)

#### Taille du widget (*requis*)

Ce paramètre permet de régler la taille du widget.

![Paramètre Taille du widget](../../img/settings/widget-size.png "Paramètre Taille du widget")

La première information à renseigner est la ligne dans laquelle le widget doit apparaître. Ce champ permet de rechercher parmi les lignes disponibles. Si aucune ligne n'est disponible, ou pour en créer une nouvelle, entrez son nom, puis appuyez sur la touche Entrée.

Ensuite, les 3 champs en dessous permettent de définir respectivement la largeur occupée par le widget sur mobile, tablette, de ordinateur de bureau.
La largeur maximale est de 12 colonnes pour un widget, la largeur minimale est de 3 colonnes.

#### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

#### Date de fin (*requis*)

Ce paramètre permet de définir la date de fin de calcul des statistiques pour ce widget.

Exemple: Si la [durée](#duree-requis) a été paramètrée à 2 mois, et la date de fin au 01 Septembre 2018, les statistiques affichées correspondront aux valeurs pour la période allant du 01 Juillet 2018 au 01 Septembre 2018.

Les dates disponibles dépendent de l'unité de [durée](#duree-requis) sélectionnée.

- Pour une durée en mois, la date de fin ne peut être que le premier jours d'un mois, à 00:00 GMT.
- Pour toutes les autres unités de durée, la date de fin doit correspondre à une heure pleine. Ex: 17/01/2019 18:00.

#### Durée (*requis*)

Ce paramètre permet de définir la durée à prendre en compte pour le calcul des statistiques.

Exemple: Une durée de 2 mois permet de calculer les statistiques concernant les 2 mois précédents la [date de fin](#date-de-fin-requis).

Valeur par défaut: 1 jour.

#### Groupes de statistiques (*requis*)

Ce paramètre vous permet de créer des groupes. Les statistiques pourront alors êtres organisées dans l'histogramme. Chaque statistique sera affiché pour chacun des groupes.

Ces groupes sont constitués d'un nom, et d'un filtre.

Ex: Affichage de 2 statistiques, pour 2 groupes différents.

// INSERER IMAGE

Les groupes sont utiles nottament pour afficher les statistiques souhaitées, en les regroupant par sous-parties de votre SI.

**Il est obligatoire d'ajouter au moins un groupe**

Pour ajouter un groupe, cliquez sur le boutton ```Ajouter un groupe```. Une fenêtre s'ouvre alors vous demandant de renseigner le nom du groupe. Un bouton ```Editeur de filtre``` est également présent, vous permettant de définir le filtre à appliquer afin de sélectionner les entités appartenants à ce groupe. Pour plus de détails sur les filtres, et l'édition de filtres, cliquez [ici](../../../filtres).
Cliquez sur ```Sauvegarder``` pour terminer l'ajout du groupe.

La liste des groupes ajoutés apparaît alors dans le panneau de paramètres du widget. Chaque groupe dispose d'un bouton d'édition, ainsi que d'un bouton de suppression.

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

##### Couleurs des statistiques (*optionnel*)

Ce paramètre vous permet de définir la couleur que vous souhaitée affiché pour chacune des statistiques sélectionnées.

La liste des statistique est affichée, avec un bouton ```Sélectionner une couleur```, ainsi que la couleur déjà selectionnée (s'il y en a une).

Pour sélectionner une couleur, cliquez sur le bouton ```Sélectionner une couleur```. Une fenêtre s'affiche. Plusieurs modes de sélection de couleur sont accessibles.

Sélectionnez la couleur souhaitée, puis cliquez sur le boutton ```Envoyer```. La couleur a été sauvegardée.
