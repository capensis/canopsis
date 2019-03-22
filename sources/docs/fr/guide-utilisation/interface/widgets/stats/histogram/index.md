# Histogramme de statistiques

![Histogramme de statistiques](./img/histogram.png "Histogramme de statistiques")

## Sommaire
### Guide utilisateur

1. [Présentation du widget](#presentation-du-widget)

### Guide exploitant

1. [Paramètres du widget](#parametres-du-widget)

## Guide utilisateur

### Présentation du widget

Le widget Histogramme de statistiques vous permet de représenter un ensemble de statistiques, sous forme d'histogramme.

Il vous permet d'afficher les statistiques souhaitées pour une ou plusieurs statistiques, sur une période définie.

## Guide exploitant

### Paramètres du widget

1. Taille du widget (*requis*)
2. Titre (*optionnel*)
3. Interval de date (*requis*)
4. Sélecteur de statistique (*requis*)
5. Filtre (*optionnel*)
6. Couleurs des statistiques (*optionnel*)

#### Taille du widget (*requis*)

Ce paramètre permet de régler la taille du widget.

![Paramètre Taille du widget](../../img/settings/widget-size.png "Paramètre Taille du widget")

La première information à renseigner est la ligne dans laquelle le widget doit apparaître. Ce champ permet de rechercher parmi les lignes disponibles. Si aucune ligne n'est disponible, ou pour en créer une nouvelle, entrez son nom, puis appuyez sur la touche Entrée.

Ensuite, les 3 champs en dessous permettent de définir respectivement la largeur occupée par le widget sur mobile, tablette, de ordinateur de bureau.
La largeur maximale est de 12 colonnes pour un widget, la largeur minimale est de 3 colonnes.

#### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

#### Interval de date (*requis*)

Ce paramètre permet de définir l'interval de dates pour lequel les statistiques doivent être affichées.

Par défaut l'interval correspond aux statistiques du jour.

##### Période

Les deux champs de période correspondent à l'interval entre deux valeurs des statistiques. Sur l'histogramme, cela se traduira par le temps écoulé entre deux points de calcul des statistiques.

**Pour éviter un temps d'affichage des statistiques trop long, il convient de sélectionner une période suffisament grande, comparée à l'interval choisie**

Exemple: Pour une interval correspondant à une année entière, le temps d'affichage des statistiques heure par heure se verra considérablement augmenté. Alors que l'affichage des statistiques mois par mois prendra, lui, un temps beaucoup plus raisonnable.

##### Interval

Deux permettent ici de sélectionner une date de début, ainsi qu'une date de fin de calcul des statistiques. Le troisième champ (à droite) permet, lui, de sélectionner un interval parmis ceux prédéfinis.

A l'intérieur des champs de sélection de date (gauche), il est possible :

- De sélectionner une date fixe, en cliquant sur l'icone de calendrier, puis en sélectionnant la date voulue
- De sélectionner une date 'dynamique'

###### Langage de sélection de date dynamique

- Le champ doit toujours commencer par le mot clé 'now', faisant référence à la date actuelle
- Ce mot clé now peut être suivi directement de modificateurs. Ces modificateurs sont de la forme :

    * Opérateur: '+' ou '-'
    * Valeur: Nombre d'unités à ajouter/soustraire
    * Unité: 'h' pour 'heures, 'd' pour 'jours', 'm' pour 'mois' et 'y' pour 'année

- A la suite du modificateur peut s'ajouter un opérateur permettant d'arrondir la valeur au début/à la fin de l'unité voulu. Cet opérateur se présente sous la forme: '/unité'. Cet arrondis se fera à la valeur inférieur pour la date de début, à la valeur supérieur pour la date de fin.

Exemples: 

- 'now-7d' -> 'La date d'aujourd'hui moins 7 jours'
- 'now-2m' -> 'La date d'aujourd'hui moins 2 mois'
- 'now-7d/d'

    * Si cette valeur correspond à une date de début -> 'La date d'aujourd'hui moins 7 jours, arrondie au début de la journée'
    * Si cette valeur correspond à une date de fin -> 'La date d'aujourd'hui moins 7 jours, arrondie à la fin de la journée'

##### Intervals prédéfinis

Le champ de droite vous permet de sélectionner parmis un panel d'intervals de dates prédéfinis, afin de ne pas avoir à entrer manuellement l'interval voulu dans les champs de gauche.

#### Sélecteur de statistiques (*requis*)

Ce paramètre permet de définir les statistiques à afficher dans l'histogramme.

**Il est obligatoire d'ajouter au moins une statistique**

Pour ajouter une statistique, cliquez sur le bouton ```Ajouter une statistique```.

Une fenêtre s'ouvre.

![Modale ajout de statistique](../img/add-stat.png "Modale ajout de statistique")


Cette fenêtre vous permet de définir la statistique souhaitée.

- Statistique à afficher (voir [liste des statstiques disponibles](../index.md#les-statistiques-disponibles)).
- Titre associé à cette statistique.
- Options: Liste d'options concernant la statistique sélectionnée. Les options varient selon la type de statistique voulue :
    - ```Recursif```: Si l'option est activée, permet de calculer la statistique sur l'entité, ainsi que sur ses dépendances, et les dépendances de ses dépendances, etc...
    - ```Etats```: Permet de ne prendre en compte que les alarmes avec le/les état(s) (ok, mineure, majeure ou critique) sélectionné(s).
    - ```Auteurs```: Permet de ne prendre en compte que les alarmes dont le/les auteur(s) fait parti de la liste précisée ici. Pour ajouter un auteur à la liste, entrez son nom, puis appuyer sur la touche "Entrée".
    - ```Sla```: Permet de préciser le temps définit dans le SLA. **Attention: Ce paramètre est requis pour le calcul des statistiques "Taux d'Ack conforme SLA" et "Taux de résolution conforme Sla"**.

Cliquez sur le bouton ```Envoyer``` pour ajouter cette statistique.

La liste des statistiques ajoutées au widget est visible depuis le panneau de paramètres du widget. Un bouton vous permet ici d'éditer la statistique, ou de la supprimer de la liste.

![Liste de statistiques](../img/stats-list.png "Liste de statistiques")

#### Filtre (*optionnel*)

Ce paramètre vous permet de définir le filtre à appliquer à la sélection d'entité. Il permet de ne sélectionner que les entités pour lesquels on souhaite afficher les statistiques.

Pour créer un filtre, ou éditer le filtre deja présent, cliquez sur le bouton ```Créer/Editer```.
Pour supprimer le filtre deja existant, cliquez sur le bouton situé à droite du bouton d'édition/création.

Au clic sur le bouton ```Créer/Editer```, une fenêtre d'édition de filtre s'ouvre. Une fois le nom du filtre, et le filtre lui-même renseignés, cliquez sur le bouton ```Envoyer``` pour le sauvegarder.

Pour plus de détails sur les filtres, et l'édition de filtres, cliquez [ici](../../../filtres).

#### Couleurs des statistiques (*optionnel*)

Ce paramètre vous permet de définir la couleur que vous souhaitée affiché pour chacune des statistiques sélectionnées.

La liste des statistique est affichée, avec un bouton ```Sélectionner une couleur```, ainsi que la couleur déjà selectionnée (s'il y en a une).

Pour sélectionner une couleur, cliquez sur le bouton ```Sélectionner une couleur```. Une fenêtre s'affiche. Plusieurs modes de sélection de couleur sont accessibles.

Sélectionnez la couleur souhaitée, puis cliquez sur le boutton ```Envoyer```. La couleur a été sauvegardée.
