# Indicateurs statistiques et KPI

!!! Info
    Disponible uniquement en édition Pro.

Canopsis fournit des indicateurs statistiques et des indicateurs de performance (KPI).

#### Table des matières
1. [Utilisation](#utilisation)<br>
2. [Graphiques](#graphiques)<br>
 A. [Métriques d’alarmes](#metriques-dalarmes)<br>
 B. [Evaluations](#Evaluations)<br>
 C. [SLI](#sli)<br>
3. [Filtres](#filtres)<br>
4. [Paramètres d’évaluation](#parametres-devaluation)<br>
5. [Glossaire des paramètres de comparaison](#glossaire-des-parametres-de-comparaison)<br>
 A. [Compteurs](#compteurs)<br>
 B. [Taux](#taux)<br>
 C. [Temps](#temps)<br>
 D. [Criticité](#criticite)<br>

## Utilisation

Dans le menu principal de Canopsis, cliquer sur le menu administration :

![Menu Principal](./img/kpi_menu_1.png)

Dans le menu administration, cliquer sur le menu KPI :

![Menu Administration](./img/kpi_menu_2.png)

Fonctionnalités disponibles :

![Fonctionnalités](./img/kpi_fonctionnalites.png)

## Graphiques
### Métriques d’alarmes

Permet de créer un histogramme en fonction des paramètres

Valeurs configurables :

* plage de temps, personnalisée ou via les valeurs usuelles
* échantillonage : étendue de la donnée représentée (heure/jour/semaine/mois)
* [filtres](#filtres)
* [paramètres](#glossaire-des-parametres-de-comparaison), pour choisir l’indicateur souhaité

*Par ex. état du nombre d’alarmes actives et acquittées par jour sur une période de six mois :*

![Métrique d’alarme](./img/kpi_alarmes.png)

### Evaluations

Permet d'aggréger les données sous forme de tableau de colonnes.

Valeurs configurables :

* plage de temps, personnalisée ou via les valeurs usuelles
* [filtres](#filtres)
* [critères](#parametres-devaluation)
* [paramètres](#glossaire-des-parametres-de-comparaison), pour choisir l’indicateur souhaité

*Par ex. le nombre d’alarmes actives avec acquittements, rangées par rôle :*

![Notations](./img/kpi_notations.png)

### SLI

Le Service Level Indicator montre le temps passé par le SI en bon fonctionnement, en maintenance et en panne.

Valeurs configurables :

* plage de temps, personnalisée ou via les valeurs usuelles
* échantillonage : étendue de la donnée représentée (heure/jour/semaine/mois)
* afficher : présentation en pourcentage ou en valeurs absolues
* [filtres](#filtres)

*Par ex. la répartition en pourcentage du temps passé dans les divers états chaque jour :*

![SLI](./img/kpi_sli.png)

## Filtres

Les filtres permettent de restreindre le périmètre pour le SLI/SLA. (en utilisant le pattern d’entités)

Pour créer un filtre, cliquer sur le bouton "plus" :

![Création de filtre 1/5](./img/kpi_filtres_1.png)

Cliquer sur le bouton "Ajouter".

![Création de filtre 2/5](./img/kpi_filtres_2.png)

Cliquer sur le bouton "plus" pour ajouter une règle.

![Création de filtre 3/5](./img/kpi_filtres_3.png)

Passer en mode avancé pour une comparaison de texte.

![Création de filtre 4/5](./img/kpi_filtres_4.png)

Ajouter ma règle de comparaison.

*Par ex. une regex sur le nom de mon application :*

![Création de filtre 5/5](./img/kpi_filtres_5.png)

## Paramètres d’évaluation

Les paramètres d’évaluation contrôlent les critères utilisables pour les [notations](#notations).
<!--- TODO: étendre l’explication (l’actuelle renvoit aux notation, qui renvoit récursivement ici) --->

*Paramètres par défaut :*

![Paramètres](./img/kpi_parametres.png)

## Glossaire des paramètres de comparaison

### Compteurs

| Nom            | Description     |
|----------------|-----------------|
| Alarmes Créées | Nombre d’alarmes créées pendant une période donnée (hors comportements périodiques). |
| Alarmes Résolues  | Nombre d’alarmes résolues (y compris les alarmes annulées) pendant une période donnée (hors comportements périodiques). |
| Alarmes Annulées  | Nombre d’alarmes annulées pendant une période donnée (hors comportements périodiques).  |
|  Alarmes Acquittées |  Nombre d’alarmes acquittées pendant une période donnée (hors comportements périodiques). |
| Nombre d’alarmes en cours pendant la période  | Nombre d’alarmes qui étaient actives pendant une période donnée (hors comportements périodiques).  |
| Nombre d’alarmes actuellement en cours  | Nombre d’alarmes qui sont actives à un instant donné.  |
| Nombre d’alarmes acquittées actuellement en cours  | Nombre d’alarmes actives et qui ont été acquittées à un instant donné.  |
| Nombre d’alarmes non acquittées actuellement en cours  |  Nombre d’alarmes actives et qui n’ont pas été acquittées à un instant donné. |
| Alarmes Actives | Nombre d’alarmes actives pendant une période donnée. |
| Alarmes Non-affichées | Nombre d’alarmes cachées. |
| Alarmes en cours de correction automatique | Nombre d’alarmes en cours de correction par un job inclus dans une consigne. |
| Alarmes avec PBehavior | Nombre d’alarmes couvertes par un comportement périodique d’arrêt de surveillance. |
| Alarmes Corrélées | Nombre d’alarmes déclenchées par une règle de corrélation. |
| Accusés de réception annulés | Nombre d’acquittements (ack) annulés. |
| Alarmes avec acks | Nombre d’alarmes ayant été acquittées. |
| Alarmes actives sans tickets | Nombre d’alarmes pour lesquelles aucun ticket n’est associé. |

### Taux

| Nom            | Description     |
|----------------|-----------------|
| Taux d’acquittement conforme SLA  | Pourcentage d’alarmes qui ont été acquittées en un temps conforme à un SLA (relatif à une configuration par widget et hors comportements périodiques)  |
| Taux de résolution conforme SLA  |  Pourcentage d’alarmes qui ont été résolues (y compris annulées) en un temps conforme à un SLA (relatif à une configuration par widget et hors comportements périodiques). |
| Taux d’alarmes corrélées | Pourcentage d’alarmes déclenchées par une règle de corrélation. |
| Taux d’alarmes avec correction automatique | Pourcentage d’alarmes corrigées par un job inclus dans une consigne. |
| Taux d’alarmes avec tickets créés | Pourcentage d’alarmes associés à un ticket. |
| Taux d’alarmes non-affichées | Pourcentages d’alarmes cachées. |

### Temps

| Nom            | Description     |
|----------------|-----------------|
|  Proportion du temps dans une criticité |  Pourcentage de temps passé par une entité dans une ou plusieurs criticités (relatif à une configuration par widget et hors comportements périodiques).   |
|  Temps passé dans une criticité | Temps passé par une entité dans une ou plusieurs criticités (relatif à une configuration par widget et hors comportements périodiques).  |
| Temps moyen entre les pannes (MTBF)  | Temps moyen entre les pannes sur une entité (hors comportements périodiques)  |
| Délai moyen d'acquittement des alarmes | Temps moyen d'acquittement des alarmes |
| Temps moyen pour résoudre des alarmes | Temps moyen avant passage des alarmes en état résolu. |
| Durée totale de l’activité | Temps de connexion par utilisateur à l’interface web de Canopsis |

### Criticité

| Nom            | Description     |
|----------------|-----------------|
|  Criticité courante | Criticité actuelle d’une entité : 0 pour OK, 1 pour Mineur, 2 pour Majeur et 3 pour Critique (comme dans la météo de service). |
