# Bac à alarmes

![Bac à alarmes](./img/listalarm.png  "Bac à alarmes")

## Sommaire

### Guide utilisateur

1. [Alarmes](#alarmes)
2. [Recherche](#recherche)
3. [Filtres](#filtres)
4. [Actions](#actions)
5. [Éléments par page](#elements-par-page)
6. [Suivi personnalisé](#suivi-personnalise)
7. [Lien direct vers une alarme](#lien-direct-vers-une-alarme)

### Guide exploitant

1. [Aide sur les variables](#aide-variables)
2. [Paramètres du widget](#parametres-du-widget)

## Guide utilisateur

### Alarmes

Le tableau d'alarmes présente la liste des alarmes. Une ligne correspond à une alarme.
Les colonnes affichées sont personnalisables (*Cf: [Guide exploitant](#guide-exploitant)*).
En plus de détails de l'alarme, chaque ligne expose une liste d'actions opérables sur l'alarme (*Cf: [Actions](#actions)*).

Il est possible d'attacher à chaque colonne une Info popup, qui s'ouvrira au clic sur le texte de la colonne, présentant dans une fenêtre un texte personnalisable (*Cf: [Guide exploitant](#guide-exploitant)*).

Au clic sur une alarme (en dehors du texte des colonnes), la chronologie de l'alarme s'affiche.

![Chronologie de l'alarme](./img/timeline.png "Chronologie de l'alarme")

Cette chronologie reprend certains éléments du cycle de vie de l'alarme (notamment les actions effectuées sur celle-ci).

### Recherche

Le champ de recherche permet de réaliser une recherche parmi les alarmes.

![Champ de recherche](../../recherche/img/champ-recherche.png "Champ de recherche")

Pour faire une recherche 'simple', il suffit d'entrer les termes de la recherche dans le champ de texte, puis d'appuyer sur la touche Entrée, ou de cliquer sur l'icone ![Icone recherche](../../recherche/img/search-icon.png "Icone recherche")

Dans le bac à alarmes, il est possible d'effectuer des recherches plus avancées. Une aide concernant la syntaxe à utiliser est disponible en survolant avec la souris l'icone d'aide ![Icone aide recherche avancée](./img/advanced-search-icon.png "Icone aide recherche avancée"). Une documentation est également disponible pour cette aspect [ici](../../recherche/index.md) !

Pour supprimer la recherche, cliquez sur l'icone ![Icone suppression recherche](../../recherche/img/delete-search-icon.png "Icone suppression recherche")

### Filtres

Le sélecteur de filtre permet d'appliquer un filtre sur le Bac à alarmes. Seules les alarmes correspondant aux critères du filtres seront affichées.

![Sélecteur de filtre](../../filtres/img/filter-selector.png "Sélecteur de filtre")

Pour sélectionner un filtre, il suffit de cliquer sur le champ 'Sélectionner un filtre'. Une liste des filtres disponibles apparaît.
Cliquez sur un filtre. Celui-ci est sélectionné et directement appliqué.
Pour ne plus appliquer de filtre, il suffit de cliquer sur l'icone présent au bout du champ de sélection de filtre. Le bac à alarmes se rafraichit, le champ de sélection revient dans état initial, le filtre n'est plus appliqué !

#### Mix filters

L'option "Mix filters", présente à gauche du sélecteur de filtre permet de cumuler plusieurs filtres.

Pour activer cette option, cliquez sur le bouton ![Mix filters](../../filtres/img/mix-filters.png "Mix filters").
Une fois l'options activée, un sélecteur apparaît à droite du bouton d'activation ![Mix filters operator](../../filtres/img/mix-filters-operator.png "Mix filters operator"). Ce sélecteur permet de choisir l'opérateur utilisé pour réunir les filtres.

- "AND": Les critères présents dans tout les filtres doivent êtres vérifiés
- "OR": Les critères présents dans un ou plusieurs des filtres doivent êtres vérifiés.

Une fois l'opérateur sélectionné, il ne vous reste plus qu'à sélectionner les filtres à appliquer dans le menu déroulant de sélection de filtres.

#### Suivi personnalisé

Le Suivi personnalisé sert à paramétrer des filtres par période. Ils permet de filtrer les alarmes en ne conservant que les alarmes d'une période donnée.

Ce filtre est disponible en cliquant sur l'icone ![Filtre par période](./img/period-filter.png "Filtre par période") présente à droite du sélecteur de filtre. Une fenêtre apparaît.

![modale filtre par période](./img/modal-filtre-periode.png "modale filtre par période")

Il suffit alors de sélectionner la période souhaitée parmi les périodes prédéfinies, ou d'en créer une personalisée en sélectionnant 'Custom', puis en renseignant les dates de début et de fin.

Dans un bac à alarmes en cours, le filtre est appliqué sur la date de création.

Dans un bac à alarmes résolues, le filtre est appliqué sur la date de résolution.

Cliquez ensuite sur 'Appliquer'.

La fenêtre se ferme, le bac à alarmes se rafraîchit. Votre filtre par période est appliqué.
Celui-ci est visible en haut du Bac à alarmes.

![Filtre par période selectionné](./img/filter-current-period.png "Filtre par période selectionné")

Afin de supprimer ce filtre, cliquez sur le bouton de fermeture présent sur le filtre (*Cf Image ci-dessus*)

#### Lien direct vers une alarme

Vous pouvez accéder à une alarme en particulier grâce à une URL directe.  
Cette URL est de la forme : `http(s)://URL_CANOPSIS/alarms/<alarmID>[?widgetId=<widgetID>]`.  

* `<alarmID>` **(requis)** : correspond à l'attribut `_id` de l'alarme.
* `<widgetID>` **(optionnel)** : correspond à l'identifiant d'un widget. Lorsque cet identifiant est précisé, la configuration du widget s'applique (colonnes, plus d'infos, etc.)

L'identifiant d'un widget est disponible pour copie dans le mode **édition** d'une vue en bas d'un widget.

### Actions

Pour chaque alarme, des actions sont disponibles.

Pour le détail de chacune des actions, voir la [liste des actions](actions.md) du Bac à alarmes.

### Éléments par page

Le champ 'Eléments par page' permet de sélectionner le nombre d'alarmes à afficher sur chaque page.

Le choix par défaut est réglable dans les paramètres du bac à alarmes (*Cf: [Guide exploitant](#guide-exploitant)*)

### Suivi personnalisé

Le champ 'Eléments par page' permet de sélectionner le nombre d'alarmes à afficher sur chaque page.

Le choix par défaut est réglable dans les paramètres du bac à alarmes (*Cf: [Guide exploitant](#guide-exploitant)*)

## Guide exploitant

Vous pouvez configurer les widgets (taille, remplacement, nom, etc.) directement dans une vue via le mode édition (*Cf: [Vues - Documentation de la grille d'edition](../../vues/edition-grille.md)*).

### Aide - Variables

Durant la configuration de votre widget Bac à alarmes, notamment paramètres "Info popup", et "Fenêtre Plus d'infos", il vous sera possible d'accéder à des variables concernant les alarmes et les entités.

Exemple : Il vous sera possible d'afficher, dans la fenêtre "Plus d'infos", la criticité de l'alarme.

Afin de connaitre les variables disponibles, une modale d'aide est disponible.

Pour y accèder, entrez dans le mode d'édition (*Cf: [Vues - Mode d'édition](../../vues/index.md#mode-édition)*).

Un bouton d'action supplémentaire "Liste des variables disponibles" apparaît alors pour chaque alarme.

Au clic sur ce bouton, une fenêtre s'ouvre. Celle-ci liste toutes les variables disponibles dans vos différents paramètres. Un bouton, à droite de chacune des variables, vous permet de copier directement dans le Presse-papier le chemin de cette variable.

### Paramètres du widget

1. Titre
2. Paramètres avancés
    1. Colonne de tri par défaut
    2. Nom des colonnes
    3. Nombre d'éléments par page par défaut
    4. Filtre sur open/resolved
    5. Filtres

#### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

#### Paramètres avancés

##### Colonne de tri par défaut

Ce paramètre permet de définir la colonne par laquelle trier les alarmes.

À noter : par défaut, le tri sur un bac à alarmes se base sur la date de création.

![Paramètre colonne de tri par défaut](../img/settings/default-column-sort.png "Paramètre colonne de tri par défaut")

Un champ de texte vous permet d'abord de définir la colonne à utiliser. Il faut ici entrer la **valeur** de la colonne, et non son nom.

Exemple : pour trier sur la base de la colonne que vous avez nommée "Connecteur", avec comme valeur "v.connector" (*Cf: [Paramètre "Nom des colonnes"](#nom-des-colonnes)*), il faut entrer ici "v.connector" et non "Connecteur".

Un sélecteur vous permet ensuite de définir le sens de tri :

*  "ASC" = Ascendant
*  "DESC" = Descendant

##### Nom des colonnes

Ce paramètre permet de définir quels colonnes seront affichées dans le bac à alarmes.

![Paramètre Nom des colonnes](../img/settings/column-names.png "Paramètre Nom des colonnes")

Afin d'**ajouter une colonne**, cliquez sur le bouton 'Ajouter'.
Une colonne vide est alors ajoutée. Afin de finaliser l'ajout, il est nécessaire de remplir les champs demandés.

Le champ "Label" définit le nom de la colonne, qui sera affiché en haut de tableau. Le champ "Valeur" définit la valeur que doit prendre ce champ. Tous les champs de l'alarme et de l'entité concernée par l'alarme peuvent être utilisés.

Voici quelques exemples pratiques de colonnes :

###### Champs basiques

Label  | Valeur
--|--
Type de connecteur | `alarm.v.connector `
Nom du connecteur | `alarm.v.connector_name`
Composant | `alarm.v.component`
Ressource | `alarm.v.resource`
Message | `alarm.v.output`
Criticité | `alarm.v.state.val`
Statut | `alarm.v.status.val`

###### Champs enrichis

Label  | Valeur
--|--
Nom du champ enrichi	| `infos.NOM_DU_CHAMP_ENRICHI`

###### Dates

Label  | Valeur
--|--
Date de création | `alarm.v.creation_date`
Date du dernier changement de criticité | `alarm.v.state.t`
Date de fin | `alarm.v.resolved`
Durée de l'alarme | `alarm.v.duration`

###### Acquittement

Label  | Valeur
--|--
Auteur de l'acquittement | `alarm.v.ack.a`
Message de l'acquittement | `alarm.v.ack.m`

###### Ticket

Label  | Valeur
--|--
Auteur du ticket | `alarm.v.ticket.a`
Numéro du ticket | `alarm.v.ticket.val`
Message du ticket | `alarm.v.ticket.m`
Type du ticket | `alarm.v.ticket._t`

###### Mise en veille

Label  | Valeur
--|--
Auteur de la mise en veille | `alarm.v.snooze.a`

Pour supprimer une colonne, cliquez dans la liste des colonnes sur la croix rouge présente en haut à droite de la case de la colonne que vous souhaitez effacer.

Dans la liste des colonnes sont également présentes, pour chaque colonne, des flèches permettant de modifier l'ordre des colonnes. Les colonnes sont présentées dans l'ordre de haut en bas. Pour modifier la place d'une colonne, cliquez sur une des flèches. Pour faire monter/descendre une colonne dans la liste.

Enfin, une option est présente pour chaque colonne, permettant d'activer (ou non) l'interprétation HTML de la valeur présente dans cette colonne.

Exemple: Vous souhaitez afficher la valeur du champ ```output``` des alarmes. Vous ajoutez donc une colonne ayant pour valeur ```alarm.v.output```. Ce champ a pour valeur ```<p style="color: red;">Exemple d'output</p>```. Si l'option ```HTML``` est desactivée, la valeur du champ sera affichée telle quelle. Si elle est activée, le code HTML sera alors interprété.

Il est à noter que seuls certaines balises et attributs sont autorisés dans les colonnes du Bac à alarmes.

- Balises autorisées: ```h3, h4, h5, h6, blockquote, p, a, ul, ol, nl, li, b, i, strong, em, strike, code, hr, br, div, table, thead, caption, tbody, tr, th, td, pre, iframe, span, font, u```

- Attributs autorisés :
    - Pour toutes les balises: ```style```
    - Pour les balises ```a```: ```href, name, target```
    - Pour les balises ```img```: ```src, alt```
    - Pour les balises ```font```: ```color, size, face```

##### Nombre d'éléments par page par défaut

Ce paramètre permet de définir combien d'éléments seront affichés, par défaut, pour chaque page du bac à alarmes.

Pour modifier ce paramètre, sélectionnez simplement la valeur souhaitée.

Les valeurs disponibles sont : 5, 10, 20, 50 et 100.

##### Filtre sur Open/Resolved

Ce paramètre permet de filtrer les alarmes en fonction de leur état de résolution.

*  Open : Alarmes "Ouvertes"
*  Resolved : Alarmes "Résolues"

Pour modifier ce paramètre, sélectionnez les types d'alarmes que vous souhaitez afficher en cochant la case correspondante.

Il est possible de ne cocher aucune des cases (aucune alarme ne sera affichée), une des deux cases, ou les deux cases (les alarmes ouvertes ET résolues seront alors affichées).

Lorsqu'une alarme est résolue, elle reste entre 1 et 2 minutes dans le bac à Alarmes "Ouvertes" avant de basculer dans le bac à Alarmes "Résolues". 

Lorsqu'une alarme est annulée, elle reste pendant 1 heure dans le bac à Alarmes "Ouvertes" avant de passer dans le bac à Alarmes "Résolues".

##### Filtres

Ce paramètre permet de sélectionner un filtre à appliquer au bac à alarmes, et d'en créer de nouveaux.

Un champ de sélection permet d'abord de choisir un filtre à appliquer au bac à alarmes parmi les filtres existants. Sélectionnez le filtre que vous souhaitez appliquer parmi les filtres disponibles. Une fois les paramètres sauvegardés, le filtre sera appliqué au bac à alarmes (*Cf: [filtres](#filtres)*).

Pour créer un nouveau filtre, cliquez sur le bouton 'Ajouter'. Une fenêtre de création de filtre s'ouvre. Pour plus de détails sur les filtres et leur création, cliquez [ici](../../filtres/index.md).

Une fois votre filtre créé, celui-ci apparaît dans la liste disponible en dessous du sélecteur de filtre. Cette liste vous permet d'éditer ou de supprimer les filtres.

L'option "Mix filters" est également disponible depuis ce menu. Pour plus de détails concernant cette option, voir  [Mix filters](#mix-filters).

Voici quelques exemples pratiques de filtres :

###### Champs basiques

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Composant  | `component`  | `equal`  | *VALEUR_DU_COMPOSANT*
Ressource  | `resource`  | `equal`  | *VALEUR_DE_LA_RESSOURCE*
Connecteur	| `connector` | `equal` | *VALEUR_DU_CONNECTEUR*
Connecteur	| `connector_name` | `equal` | *VALEUR_DU_NOM_DU_CONNECTEUR*
Message	| `v.output` | `equal` | *VALEUR_DU_MESSAGE*

###### Selon la criticité

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Uniquement les alarmes Mineures  | `v.state.val`  | `equal`  | `1` (valeur de type number)
Uniquement les alarmes Majeures  | `v.state.val`  | `equal`  | `2` (valeur de type number)
Uniquement les alarmes Critiques  | `v.state.val`  | `equal`  | `3` (valeur de type number)

###### Champs enrichis

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Champ enrichi	| `entity.infos.NOM_DU_CHAMP_ENRICHI.value` | equal | *VALEUR_DU_CHAMP_ENRICHI*

###### En fonction des informations dynamiques

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Les alarmes qui contiennent une information dynamique de type `consignes`	| `v.infos.*.type` | equal | *consigne*

###### Acquittement

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Exclure les alarmes avec acquittement  | `v.ack._t`  | `not equal` | `ack` (valeur string)
Uniquement les alarmes avec acquittement  | `v.ack._t`  | `equal` | `ack` (valeur string)
Uniquement les alarmes avec acquittement sans champ `Note` (fast-ack basique)  | `v.ack.m`  | `is empty` | *PAS_DE_VALEUR*
Exclure les alarmes avec acquittement avec un champ `Note`  | `v.ack.m`  | `is not empty` | *PAS_DE_VALEUR*
Auteur de l'acquittement  | `v.ack.a`  |  `equal`  | *NOM_DE_L_AUTEUR_DE_L_ACQUITTEMENT*
Message de l'acquittement | `v.ack.m`  |  `equal`  | *CONTENU_DU_MESSAGE_DE_L_ACQUITTEMENT*

###### Ticket

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Exlure les alarmes avec Ticket (quel que soit le type)  | `v.ticket._t`  | `is null` | *PAS_DE_VALEUR*
Exclure les alarmes avec Ticket de type `assocticket`  | `v.ticket._t`  | `not equal` | `assocticket` (valeur string)
Exclure les alarmes avec Ticket de type `declareticket`  | `v.ticket._t`  | `not equal` | `declareticket` (valeur string)
Uniquement les alarmes avec Ticket  | `v.ticket._t`  | `is not null` | *PAS_DE_VALEUR*
Auteur du Ticket  | `v.ticket.a`  |  `equal`  | *NOM_DE_L_AUTEUR_DU_TICKET*

###### Mise en veille

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Exclure les alarmes mises en veille | `v.snooze._t`  | `not equal` | `snooze` (valeur string)
Uniquement les alarmes mises en veille | `v.snooze._t`  | `equal` | `snooze` (valeur string)
Auteur de la mise en veille | `v.snooze.a`  |  `equal`  | *NOM_DE_L_AUTEUR_MISE_EN_VEILLE*

###### Comportements périodiques

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Uniquement les alarmes qui possèdent un comportement périodique actif (anciennement `has_active_pb`) | `pbehavior` | `exists` | `true` (valeur booléenne)
Uniquement les alarmes qui ne possèdent pas de comportement périodique actif  (anciennement `has_active_pb`) | `pbehavior` | `exists` | `false` (valeur booléenne)
Uniquement les alarmes avec un comportement périodique particulier actif | `pbehavior.type.type` | `equal` | `maintenance` (valeur string)

Tous les attributs des comportements périodiques peuvent être utilisés à des fin de filtrage : `author`, `name`, `type`, `priority`, `icon_name`.

###### Changement de criticité

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Exclure les alarmes dont on a manuellement changé la criticité  | `v.state._t`  | `not equal` | `changestate` (valeur string)
Uniquement les alarmes dont on a manuellement changé la criticité  | `v.state._t`  | `equal` | `changestate` (valeur string)

###### Services

Description  | 1<sup>re</sup> colonne  | 2<sup>e</sup> colonne | 3<sup>e</sup> colonne
--|---|--|--
Exclure les alarmes liées à des services | `entity.type`  | `not equal` | `service` (valeur string)
Uniquement les alarmes des services | `entity.type`  | `equal` | `service` (valeur string)

##### Info popup

Ce paramètre permet d'ajouter un info popup sur une des colonnes du Bac à alarmes (*Cf: [Infos popup - Guide utilisateur Bac à alarmes](#alarmes)*).

Pour ajouter une info popup, cliquez sur le bouton 'Ajouter'.

Une case info popup vide apparaît.
Cette case comporte deux champs :

*  Colonne : Ce champ permet de définir sur quelle colonne l'info popup sera disponible. Il faut ici entrer la **valeur** de la colonne, et non son nom.
Exemple : pour ajouter une info popup sur la colonne que vous avez nommée "Connecteur", avec comme valeur "alarm.v.connector" (*Cf: [Paramètre "Nom des colonnes"](#nom-des-colonnes)*), il faut entrer ici "alarm.v.connector" et non "Connecteur".
*  Texte : Ce champ, qui a la forme d'un éditeur de texte, permet de définir le contenu de l'info popup. Le langage utilisé ici pour le template de la popup est l'Handlebar. Deux variables sont disponibles : "alarm" et "entity". Exemple : Pour ajouter au template la criticité de l'alarme, ajoutez au template `{{ alarm.v.state.val }}`.

Vous pouvez ajouter autant d'info popup que vous le souhaitez.

Pour supprimer une info popup, cliquez sur la croix rouge, en haut à droite de la case de l'info popup que vous souhaitez supprimer.

##### Fenêtre Plus d'infos

Ce paramètre permet de définir le contenu de la fenêtre plus d'infos. Le bouton permettant d'ouvrir cette fenêtre se trouve dans les actions de chaque alarme du bac à alarmes.

Ce champ se présente sous forme d'un éditeur de texte.
Le langage utilisé dans cet éditeur est le Handlebars.
Deux variables sont disponibles ici, 'alarm' et 'entity'.

En plus du texte que vous souhaitez afficher, il vous est donc possible d'intégrer des informations de l'alarme ou de l'entité concernée par cette alarme.

Exemple : Pour afficher la criticité de l'alarme, ajoutez `{{ alarm.v.state.val }}`.
