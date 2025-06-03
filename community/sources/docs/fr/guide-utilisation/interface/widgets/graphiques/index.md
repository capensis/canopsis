# Widgets graphiques

Les widgets graphiques permettent de visualiser des métriques sous forme de graphiques pour analyser l’état ou l’évolution du système d’information.  
Ils supportent plusieurs types de rendus : histogramme, graphique en ligne, diagramme circulaire (camembert) et valeurs chiffrées.  
Les données affichées peuvent provenir de métriques internes à Canopsis (nombre d’alarmes, taux d’ack, etc.) ou de métriques externes issues d’événements (au format perf_data).

Chaque widget peut être personnalisé (période, méthode de calcul, comparaison temporelle, couleurs…) pour répondre à des besoins d’analyse ponctuelle ou continue.

![Graphiques](./img/graphiques.png)

## Types de widgets graphiques

Canopsis propose 4 types principaux de widgets graphiques, chacun adapté à une visualisation spécifique des métriques :

- **Histogramme (Bar chart)** : pour visualiser des valeurs sur une période, empilées ou séparées.
- **Graphique en ligne (Line chart)** : pour suivre l’évolution de valeurs dans le temps.
- **Diagramme circulaire (Pie chart)** : pour montrer des proportions entre différentes métriques.
- **Valeurs numériques (Numbers)** : pour afficher directement les valeurs agrégées de certaines métriques.

## Métriques internes disponibles

Voici la liste complète des **métriques internes** disponibles dans Canopsis. Elles sont issues des alarmes, entités et actions utilisateur.


| Nom de la métrique                                      | Unité    | Méthodes de calcul disponibles | Méthode par défaut |
|---------------------------------------------------------|----------|--------------------------------|--------------------|
| Nombre d'alarmes créées                                 | nombre   | Sum, Average, Min, Max         | **Sum**            |
| Nombre d'alarmes actives                                | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes non affichées                          | nombre   | Sum, Average, Min, Max         | **Average**        |
| Nombre d'alarmes en cours de remédiation automatique    | nombre   | Sum, Average, Min, Max         | **Average**        |
| Nombre d'alarmes avec remédiation manuelle              | nombre   | Sum, Average, Min, Max         | **Average**        |
| Nombre d'alarmes avec comportement périodique           | nombre   | Sum, Average, Min, Max         | **Average**        |
| Nombre d'alarmes corrélées                              | nombre   | Sum, Average, Min, Max         | **Average**        |
| Nombre d'alarmes avec acquittement                      | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes actives avec acquittement              | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes avec acquittement annulé               | nombre   | Sum, Average, Min, Max         | **Sum**            |
| Nombre d'alarmes actives avec tickets                   | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes actives sans tickets                   | nombre   | Average, Min, Max              | **Average**        |
| % d'alarmes corrélées                                   | %        | Sum                            | **Sum**            |
| % d'alarmes avec remédiation automatique                | %        | Sum                            | **Sum**            |
| % d'alarmes avec tickets créés                          | %        | Average                        | **Average**        |
| % d'alarmes non affichées                               | %        | Sum                            | **Sum**            |
| % d'alarmes remédiées manuellement                      | %        | Sum                            | **Sum**            |
| Temps moyen pour acquitter les alarmes                  | secondes | Average, Min, Max              | **Average**        |
| Temps moyen pour résoudre les alarmes                   | secondes | Average, Min, Max              | **Average**        |
| Délai d'acquittement des alarmes                        | secondes | Average, Min, Max              | **Average**        |
| Délai de résolution des alarmes                         | secondes | Average, Min, Max              | **Average**        |
| Temps minimum pour acquitter les alarmes                | secondes | Min                            | **Min**            |
| Temps maximum pour acquitter les alarmes                | secondes | Max                            | **Max**            |
| Temps minimum pour résoudre les alarmes                 | secondes | Min                            | **Min**            |
| Temps max pour résoudre les alarmes                     | secondes | Max                            | **Max**            |
| Nombre d'alarmes non acquittées                         | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes non acquittées avec durée 1-4h         | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes non acquittées avec durée 4-24h        | nombre   | Average, Min, Max              | **Average**        |
| Nombre d'alarmes non acquittées de plus de 24h          | nombre   | Average, Min, Max              | **Average**        |


## Métriques externes (perf_data)

Les métriques externes peuvent être envoyées via des événements Canopsis contenant un champ `perf_data` :

```json
"perf_data": "cpu=20%;80;90;0;100"
```

Elles sont traitées et enregistrées dans TimescaleDB avec les informations suivantes :  

- `metric_label`
- `entity_id`
- `timestamp`
- `value`
- `unit`

Les unités supportées sont :  

- **%** : pourcentage
- **s, ms, us** : secondes, millisecondes, microsecondes
- **B, KB, MB, GB, TB** : tailles
- **c** : compteur continu
- *(ou vide pour des nombres simples)*


## Méthodes de calcul possibles

Pour toutes les métriques (internes et externes), différentes méthodes de calcul peuvent être appliquées dans les widgets :

| Méthode            | Description                                                                  |
|--------------------|------------------------------------------------------------------------------|
| **Sum**            | Somme des valeurs sur la période                                             |
| **Average**        | Moyenne des valeurs                                                          |
| **Min / Max**      | Valeur minimale ou maximale                                                  |
| **Last**           | Dernière valeur reçue (utile pour les compteurs ou l’état courant)           |


## Affichages des graphiques dans les widgets Bac à alarmes et Explorateur de contexte

Les graphiques peuvent être affichés directement dans :

- Le [**bac à alarmes**](../bac-a-alarmes/)
- L’[**explorateur de contexte**](../contexte)

Ils s’affichent dans un onglet `Graphiques` si des métriques externes ont été collectées pour l’alarme ou l’entité.

Les types de graphes supportés ici sont :

- Histogramme
- Graphique en ligne
- Nombre
_(Pas de diagramme circulaire pour ces vues)_


## Bonnes pratiques

- Utiliser les **presets** quand ils sont proposés
- Préférer des **métriques homogènes** (unités identiques) par widget
- Nommer clairement les `metric_label` pour éviter les confusions


## Nettoyage des données

- Par défaut, les **métriques internes** sont conservées 1 an
- Les **métriques externes** sont conservées 6 mois

La politique de rétention peut être ajustée dans les [paramètres de stockage](../../../menu-administration/parametres-de-stockage/)




## Cas d’usage supportés

- Une entité avec **une seule métrique** → un graphe unique
- Une entité avec **plusieurs métriques sur un même graphe**
- Une entité avec **plusieurs métriques sur plusieurs graphes distincts**
- Les métriques peuvent être ajoutées :  
  - une par une
  - ou automatiquement par **masque (regexp)** avec la fonction *auto add*


## Intégration avec les événements `perf_data`

Pour exploiter des métriques dans Canopsis via des événements (ex : venant des [connecteurs](../../../../interconnexions/)), le champ `perf_data` doit suivre ce format :

```text
'label'=value[UOM];[warn];[crit];[min];[max]
```

Exemple :

```json
"perf_data": "cpu=20%;80;90;0;100"
```

- Seules la **valeur** et l’**unité** sont prises en compte
- Les valeurs `warn`, `crit`, `min`, `max` sont **ignorées**
- Les métriques sont associées aux entités Canopsis via leur `entity_id`


##  Paramètres des widgets

Voici la liste des paramètres de widgets graphiques. Ces paramètrent varient en fonction du type de widget.

### Titre

Ce paramètre permet de définir le titre du widget, qui sera affiché au-dessus de celui-ci.

Un champ de texte vous permet de définir ce titre.

### Presets fournis pour chaque type de widget

Les presets facilitent la création de widgets en proposant des regroupements de métriques prêtes à l’emploi avec des paramètres préconfigurés (titre, période, méthode de calcul, etc).

### Histogramme

| Nom du preset                          | Affichage       | Métriques incluses                                                            | Détails                |
|----------------------------------------|-----------------|-------------------------------------------------------------------------------|------------------------|
| Nombre d’alarmes actives               | Barres séparées | Nombre d'alarmes actives                                                      | Comparaison activée    |
| Statistiques des acquittements         | Empilé          | Nombre d'alarmes avec acquittement, Nombre d'alarmes non acquittées           | Comparaison désactivée |
| Statistiques des tickets               | Empilé          | Nombre d'alarmes actives avec tickets, Nombre d'alarmes actives sans tickets  | Comparaison désactivée |
| Statistiques des acquittements annulés | Empilé          | Nombre d'alarmes avec acquittement, Nombre d'alarmes avec acquittement annulé | Comparaison désactivée |

### Diagramme circulaire

| Nom du preset                 | Métriques incluses                                                            | Méthode de calcul | Détails                  |
|-------------------------------|-------------------------------------------------------------------------------|-------------------|--------------------------|
| Répartition des acquittements | Nombre d'alarmes avec acquittement, Nombre d'alarmes non acquittées           | Sum               |                          |
| Répartition des tickets       | Nombre d'alarmes actives avec tickets, Nombre d'alarmes actives sans tickets  | Sum               |                          |
| Corrélation vs non-corrélées  | % d'alarmes corrélées, % d'alarmes non corrélées                              | Sum               | Si % non-corrélées dispo |

### Nombres

| Nom du preset                        | Métriques incluses                                                           | Méthode de calcul par défaut |
|--------------------------------------|------------------------------------------------------------------------------|------------------------------|
| Nombre d’alarmes créées              | Nombre d'alarmes créées                                                      | Sum                          |
| Nombre d’alarmes actives             | Nombre d'alarmes actives                                                     | Average                      |
| Alarmes acquittées vs non acquittées | Nombre d'alarmes avec acquittement, Nombre d'alarmes non acquittées          | Average                      |
| Alarmes avec vs sans tickets         | Nombre d'alarmes actives avec tickets, Nombre d'alarmes actives sans tickets | Average                      |
| Moyenne des temps d’acquittement     | Délai moyen d'acquittement des alarmes                                       | Average                      |
| Moyenne des temps de résolution      | Délai moyen de résolution des alarmes                                        | Average                      |


### Affichage des métriques

Il s'agit de la sélection des métriques que vous souhaitez afficher sur le graphique.

* Bouton "AJOUTER"  : Sélection de [métriques internes](#metriques-internes-disponibles)
* Bouton "AJOUTER EXTERNE" : Sélection de métriques externes
* Bouton "AJOUT AUTOMATIQUE" : Sélection de métriques externes par l'intermédiaire d'un masque

Vous pouvez ensuite définir : 

* **Etiquette affichée** : Label pour la métrique
* **Couleur personnalisée** : Couleur de la courbe, de la barre
* **Méthode de calcul** : Fonction d'agrégation appliquée (moyenne, somme, min, max, dernier)

### Paramètres avancées

#### Plage horaire par défaut

Il s'agit d'un sélecteur de plages horaires avec des valeurs usuelles.

#### Échantillonnage par défaut

Une fonction d'agrégation est appliquée sur les valeurs des métriques.  
L'échantillonage permet d'appliquer cette fonction sur une période données : heure, jour, semaine, mois

#### Filtres

Ce paramètre permet de définir les filtres pour lesquels vous souhaitez des métriques.
Pour plus de détails sur les filtres et leur création, voir la partie sur [Les filtres](../../patterns/).

Pour créer un filtre, cliquez sur le bouton 'Ajouter'. Une fenêtre de création de filtre s'ouvre alors.
Vous avez la possibilité d'éditer ou de supprimer des filtres existants.

#### Afficher la comparaison

Cette option permet d'afficher un graphique fantôme pour comparer les valeurs avec la période période précédente.

#### Tendance (pour les Nombres)

Une flêche ↑ ou ↓ est affcihée en fonction de la tendance.


