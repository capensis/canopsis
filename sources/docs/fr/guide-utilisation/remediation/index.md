# La remédiation dans Canopsis

## Introduction

!!! quote "Définition"
    Plan d’actions mis en œuvre pour corriger une situation.

La principe de `remédiation` a été mis en oeuvre dans Canopsis répondre à différents objectifs :

* Faire office de référentiel de consignes
* Mettre à disposition d'une alarme toutes les consignes adaptées
* Identifier les alarmes orphelines, c'est-à-dire sans consignes associées

Les bénéfices de l'utilisation de ce modules sont multiples 

| Bénéficiaires                 | Bénéfices                                |
|:----------------------------- |:---------------------------------------- |
| `Equipes de pilotage`         | Gain de temps                            |
|                               | Diminution du risque d’erreur            |
|                               | Observance des remédiations              |
| `Equipes de management`       | Amélioration continue du service         |
|                               | Données objectives de suivi              |
| `Le SI de manière générale`   | Référentiel pour d’autres outils         |
|                               | Communication / diffusion des résultats  |

## Anatomie d'une remédiation

La remédiation est représentée par une `Consigne` composée elle-même d'`Etapes` composées d'`Opérations`.  
Une `opération` peut être liée à un job de remédiation qui sera exécuté par un ordonnanceur de tâches.

```mermaid
graph TD
    C[Consigne] -->|1| E1(Etape 1)
    C[Consigne] -->|2| E2(Etape 2)
    E1 -->|1| O1(Opération 1)
    O1 -->|1| J1{Job 1}
    E1 -->|2| O2(Opération 2)
    E1 -->|3| O3(Opération 3)
    E2 -->|1| O4(Opération 1)
    E2 -->|2| O5(Opération 2)
    E2 -->|3| O6(Opération 3)
```            

## Prérequis

Les ordonnanceurs de tâches compatibles avec les jobs de remédiation sont :

1. [Rundeck](https://www.rundeck.com/)
2. [Awx](https://www.ansible.com/products/awx-project)

## Exécution d'une remédiation

Du point de vue "pilotage", voici comment exécuter une consigne.  

* Vérifier la présence d'une ou plusieurs consignes associées

![Présence](./img/remediation_consigne_existe.png)

* Exécuter la consigne à partir du menu d'actions

![Exécuter](./img/remediation_consigne_executer.png)

* Evaluer la consigne

Vous avez la possibilité d'évaluer la consigner que vous venez d'exécuter.
Ces évaluations seront comptabilisées et transmises aux administrateurs.

![Evaluer](./img/remediation_consigne_evaluation.png)

* Filtrer les alarmes avec ou sans consignes

![Filtrer](./img/remediation_consigne_filtres.png)

## La suite

Pour paramétrer le module de `Remédiation` dans Canopsis, vous pouvez consulter la [documentation de mise en œuvre de la remédiation](./mise-en-oeuvre.md).
