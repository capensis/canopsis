# Calendrier

Le widget Calendrier permet de visualiser les volumes d'alarmes sous forme de compteurs par jour, directement dans une interface calendaire.  
Il est particulièrement utile pour repérer des périodes critiques, des tendances récurrentes ou des pics d’activité.

Les périodes d'affichage possibles sont :

* Vue mensuelle
* Vue hebdomadaire
* Vue journalière

Il se base sur le système de [patterns](../../patterns) et chaque filtre donne lieu à l'affichage d'un compteur.

Il est utile pour :

* Identifier rapidement les journées les plus chargées
* Mettre en évidence les effets d’un incident ou d’un déploiement
* Fournir une vision temporelle synthétique des événements à traiter

![Calendrier](./img/calendrier.png  "Calendrier")

## Utilisation courante

Le widget se présente sous la forme d'un calendrier, l'utilisateur peut sélectionner la période (mensuelle, hebdomadaire, ou jour).  
Les compteurs représentant les filtres sont visibles et sont cliquables.  
Le clic permet de consulter les alarmes dans un bac à alarmes.


## Paramètres du widget

### Titre (*optionnel*)

Ce paramètre permet de définir le titre du widget, qui sera affiché au-dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

### Paramètres du bac à alarmes

Vous trouverez ici tous les paramètres relatifs au bac à alarme qui s'ouvre sous forme de modale lorsque l'utilisateur clique sur un compteur.

### Paramètres avancés

#### Filtres (*requis*)

Ce paramètre permet de définir les filtres pour lesquels vous souhaitez des compteurs.  
Pour plus de détails sur les filtres et leur création, voir la partie sur [Les filtres](../../patterns/).

Pour créer un filtre, cliquez sur le bouton 'Ajouter'. Une fenêtre de création de filtre apparaît alors.
Vous avez la possibilité d'éditer ou de supprimer des filtres existants.  

#### Filtre sur Ouverte/Résolue (*requis*)

Ce paramètre permet de choisir le contexte de calcul des filtres.

* Alarmes ouvertes : les filtres s'appliquent uniquement aux alarmes ouvertes.
* Alarmes ouvertes et récemment résolues : les filtres s'appliquent aux alarmes ouvertes ainsi que les alarmes résolues depuis moins de `TimeToKeepResolvedAlarms` (Voir [la documentation du fichier canopsis.toml](../../../../guide-administration/administration-avancee/modification-canopsis-toml/#section-canopsisalarm)).
* Alarmes résolues : les filtres s'appliquent sur les alarmes résolues uniquement.


#### Niveaux de criticté

Cet ensemble de paramètres permet de définir la couleur qui sera utilisée pour les compteurs.

* Niveaux de criticité : seuils au-delà desquels les couleurs des cases du calendrier sont appliquées
* Sélecteur de couleur : vous pouvez personnaliser les couleurs des différents seuils

