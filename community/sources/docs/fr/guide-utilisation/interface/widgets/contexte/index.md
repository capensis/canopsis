# Explorateur de contexte

L'Explorateur de contexte est un widget central dans Canopsis. Il permet d'afficher, filtrer, rechercher et manipuler l'ensemble des [entités](../../..//vocabulaire/#entite) du système : composants, connecteurs, ressources et services. C'est, en quelque sorte, le référentiel interne de toutes les entités connues de la plateforme.

Ce widget est composé de plusieurs zones clés :

* Un tableau personnalisable listant les entités
* Une barre de recherche avancée
* Un système de filtres (simples ou combinés)
* Des actions de gestion d'entités (création, modification, suppression, duplication)
* Un accès direct aux comportements périodiques

![Explorateur de contexte](./img/context-explorer.png  "Explorateur de contexte")

## Guide utilisateur

### Entités

Le tableau d'entités présente la liste de toutes les entités. Une ligne correspond à une entité.  
En plus de détails de l'entité, chaque ligne expose une liste d'actions opérables sur l'entité.

Des détails supplémentaires concernant l'entité s'affichent en la "dépliant" (comportements périodiques, informations, impacts et dépendances, etc.).

![Plus d'infos entités](./img/context-entity-more-infos.png "Plus d'infos entités").

### Recherche

Le champ de recherche permet de réaliser une recherche parmi les entités.

![Champ de recherche](./img/champ-recherche.png "Champ de recherche")

Pour faire une recherche 'simple', il suffit d'entrer les termes de la recherche dans le champ de texte, puis d'appuyer sur la touche Entrée, ou de cliquer sur l'icone :material-magnify:

Dans l'explorateur de contexte, il est possible d'effectuer des recherches plus avancées à l'aide de l'icône :material-tune:.

Pour supprimer la recherche, cliquez sur l'icone :material-close:.

### Catégorie

Une entité peut être positionnée dans une "catégorie" au moment de sa création ou de son import.  
L'explorateur de contexte permet de filtrer les entités selon ces catégories.

### Filtres

Le sélecteur de filtre permet d'appliquer un [filtre](../../patterns) sur l'Explorateur de contexte. Seuls les entités correspondant aux critères du filtres seront affichées.

![Sélecteur de filtre](./img/filter-selector.png "Sélecteur de filtre")

Pour sélectionner un filtre, il suffit de cliquer sur le champ 'Sélectionner un filtre'. Une liste des filtres disponibles apparaît.
Cliquez sur un filtre. Celui-ci est sélectionné, et directement appliqué.
Pour ne plus appliquer de filtre, il suffit de cliquer sur l'icône présente au bout du champ de sélection de filtre. L'explorateur de contexte se rafraichit, le champ de sélection revient dans état initial, le filtre n'est plus appliqué !

L'option "Mix filters", présente dans l'entête du sélecteur de filtre permet de cumuler plusieurs filtres avec un **ET logique**.

### Filtre "Aucun événement"

L'explorateur de contexte permet de filtrer les entités pour lesquelles aucun événement n'a été reçu depuis un certain temps.  
Ces entotés sont régies par les [règles d'inactivité](../../../menu-exploitation/regles-inactivite/).

### Création d'entités "services"

Depuis l'Explorateur de contexte, il vous est possible de créer des entités de type "Service".

Pour accéder aux fenêtres de création, cliquer sur le bouton ![Icône Création Entité](./img/add-entity-button.png "Icône Création Entité").

Tous les paramètres concernant les services sont [documentés ici](../../../services)

### Actions

Pour chaque entité de l'explorateur de contexte, trois actions sont disponibles :

- **Éditer** : Au clic sur l'icône d'édition ![Icône Editer entité](./img/edit-entity-icon.png "Icône Editer entité"), une fenêtre s'ouvre. Celle-ci reprend les informations de l'entité ou de le service (*Cf: [Création d'entités et de services](#creation-dentites-et-de-services)*). Après avoir modifié les informations souhaitées, cliquez sur 'Envoyer'. Une fenêtre vous informe que l'édition a été effectuée avec succès.
- **Dupliquer**: Au clic sur l'icône ![Icône Dupliquer entité](./img/duplicate-entity-icon.png "Icône Dupliquer entité"), une fenêtre s'ouvre. Celle-ci reprend les informations de l'entité ou de le service que vous souhaitez dupliquer. Après avoir entré les informations souhaitées, cliquez sur 'Envoyer'. Une fenêtre vous informe qu'une nouvelle entité a été créée avec succès !
- **Supprimer** : Permet de supprimer une entité/un service. Au clic sur l'icône de suppression ![icône Supprimer entité](./img/delete-entity-icon.png "icône Supprimer entité"), une fenêtre de confirmation s'ouvre. Cliquez sur 'Oui' pour confirmer la suppression de l'entité/de le service.
- **Ajouter un comportement périodique** : Permet d'ajouter un comportement périodique à l'entité/à le service. Au clic sur l'icône ![icône Ajouter Pbehavior](./img/add-pbehavior-icon.png "icône Ajouter Pbehavior"), une fenêtre de création de comportement périodique s'affiche. Pour plus d'information, voir : [Les comportement périodiques - Pbehaviors](../../pbehaviors/index.md).

### Comportements périodiques

Depuis l'explorateur de contexte, il est possible d'ajouter un comportement périodique directement sur une entité, ou sur un sélection d'entités.

Pour plus de détails sur l'ajout de comportements périodiques, voir : [Les comportement périodiques - Pbehaviors](../../pbehaviors/index.md).

Pour ajouter un comportement périodique sur une entité, cliquez sur l'icône ![icône Ajouter Pbehavior](./img/add-pbehavior-icon.png "icône Ajouter Pbehavior") sur la ligne de l'entité sur laquelle vous souhaitez ajouter le comportement.

Pour ajouter un comportement périodique sur une sélection d'entités, sélectionnez les entités en cochant les cases (présentes en début de ligne de chaque entités). Une fois une entité ou plus sélectionnées, deux icônes sont apparues en haut de l'explorateur de contexte. La première icône permet de supprimer toutes les entités sélectionnées, le deuxième permet d'ajouter un comportement périodique à ces entités. Cliquez sur le bouton ![icône Ajouter Pbehavior](./img/add-pbehavior-icon.png "icône Ajouter Pbehavior"). Une fenêtre de création de comportement périodique apparaît alors.

## Guide exploitant

Vous pouvez configurer les widgets (taille, remplacement, nom, etc.) directement dans une vue via le mode édition (*Cf: [Vues - Documentation de la grille d'edition](../../vues/edition-grille.md)*).

### Aide - Variables

Durant la configuration de votre widget Exporateur de contexte, notamment la liste des colonnes, il vous sera possible d'accéder à des variables concernant les entités.

Afin de connaitre les variables disponibles, une modale d'aide est disponible.

Pour y accéder, entrez dans le mode d'édition (*Cf: [Vues - Mode d'édition](../../vues/index.md#mode-édition)*).

Un bouton d'action supplémentaire "Liste des variables disponibles" apparaît alors pour chaque entité du tableau.

Au clic sur ce bouton, une fenêtre s'ouvre. Celle-ci liste toutes les variables disponibles dans vos différents paramètres. Un bouton, à droite de chacune des variables, vous permet de copier directement dans le Presse-papier le chemin de cette variable.

### Paramètres du widget

1. Titre
2. Paramètres avancés
  1. Colonne de tri par défaut
  2. Nom des colonnes
  3. Filtres
  4. Types d'entités

#### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

#### Paramètres avancés

##### Colonne de tri par défaut

Ce paramètre permet de définir la colonne par laquelle trier les alarmes.

![Paramètre colonne de tri par défaut](../img/settings/default-column-sort.png "Paramètre colonne de tri par défaut")

Un champ de texte vous permet d'abord de définir la colonne à utiliser. Il faut ici entrer la **valeur** de la colonne, et non son nom (*Cf: [Paramètre "Nom des colonnes"](#nom-des-colonnes)*).

Un sélecteur vous permet ensuite de définir le sens de tri :

*  "ASC" = Ascendant
*  "DESC" = Descendant

##### Nom des colonnes

Ce paramètre permet de définir quels colonnes seront affichées dans l'explorateur de contexte.

![Paramètre Nom des colonnes](../img/settings/column-names.png "Paramètre Nom des colonnes")

Afin d'**ajouter une colonne**, cliquez sur le bouton 'Ajouter'.
Une colonne vide est alors ajoutée. Afin de finaliser l'ajout, il est nécessaire de remplir les champs demandés.
Le champ "Label" définit le nom de la colonne, qui sera affiché en haut de tableau. Le champ "Valeur" définit la valeur que doit prendre ce champ. Tous les champs de l'entité sont directement disponibles.

Exemple : "name", qui contient le nom de l'entité, ou encore "type", qui contient le type de l'entité.

Pour supprimer une colonne, cliquez dans la liste des colonnes sur la croix rouge présente en haut à droite de la case de la colonne que vous souhaitez effacer.

Dans la liste des colonnes sont égalements présentes, pour chaque colonne, des flèches permettant de modifier l'ordre des colonnes. Les colonnes sont présentées dans l'ordre de haut en bas. Pour modifier la place d'une colonne, cliquez sur une des flèches. Pour faire monter/descendre une colonne dans la liste.

##### Filtres

Ce paramètre permet de sélectionner un filtre à appliquer à l'explorateur de contexte, et d'en créer de nouveaux.

Un champ de sélection permet d'abord de choisir un filtre à appliquer à l'explorateur de contexte parmi les filtres existants. Sélectionnez le filtre que vous souhaitez appliquer parmi les filtres disponibles. Une fois les paramètres sauvegardés, le filtre sera appliquer à l'explorateur de contexte (*Cf: [filtres](#filtres)*).

Pour créer un nouveau filtre, cliquez sur le bouton 'Ajouter'. Une fenêtre de création de filtre s'ouvre. Pour plus de détails sur les filtres et leur création, cliquez [ici](../../filtres/index.md).
Une fois votre filtre créé, celui-ci apparaît dans la liste disponible en dessous du sélecteur de filtre. Cette liste vous permet d'éditer, ou de supprimer les filtres.

L'option "Mix filters" est également disponible depuis ce menu. Pour plus de détails concernant cette option, voir  [Mix filters](#mix-filters).

#### Types d'entités

Ce paramètre permet de sélectionner les différents types d'entités que vous souhaitez voir apparaître dans l'explorateur de contexte.

![Paramètre Types d'entités](../img/settings/entities-types.png "Paramètre Types d'entités")

Les types d'entités sont : Composant, Connecteur, Ressource et Service.

Il vous suffit de cocher les cases correspondantes aux types d'entités que vous souhaitez voir apparaître.
