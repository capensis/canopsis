# Consignes
## Mise en œuvre de la remédiation

!!! Abstract "Pages associées"
    cette page fait partie d'une série de trois documentations au sujet de la
    remédiation dans canopsis :
    
    1. la page [la remédiation dans canopsis][remed-index] présente la
    fonctionnalité et montre l'exécution d'une consigne pour l'utilisateur
    final ;
    2. la page en cours décrit la configuration de la remédiation dans canopsis
    pour la gestion des consignes et des jobs ;
    3. le guide d'administration
    [la remédiation et les jobs dans canopsis][admin-remed] traite de
    l'architecture technique et de la configuration à réaliser pour exécuter des
    jobs dans les ordonnanceurs supportés.

[TOC]

### Introduction

Cette documentation vous permet de mettre en œuvre une remédiation de bout en bout.

Comme expliqué en [préambule][remed-index], l'anatomie d'une remédiation suit le
schéma ci-après :

```mermaid
graph TD
    C[Consigne] -->|1| E1(Étape 1)
    C[Consigne] -->|2| E2(Étape 2)
    E1 -->|1| O1(Opération 1)
    O1 -->|1| J1{Job 1}
    E1 -->|2| O2(Opération 2)
    E1 -->|3| O3(Opération 3)
    E2 -->|1| O4(Opération 1)
    E2 -->|2| O5(Opération 2)
    E2 -->|3| O6(Opération 3)
```

### Accès

Le menu `Administratration -> Consignes` vous permet d'accéder à la liste des consignes.  
![Consignes menu](./img/consignes_menu.png)

Sur cette page, vous trouverez la liste des consignes, leur état et quelques boutons d'actions rapides. Le processus de création de consigne est détaillé dans le chapitre [Créer une consigne](#creer-une-consigne).

![Consigne vue liste](./img/consignes_vue_liste.png)

### Les droits

Le module de remédiation est soumis au [système de droits](./droits.md) sur l'interface et sur
les API.

**Interface**

Ces droits sont configurés dans le panneau de droits sous l'onglet
« [Technical](./droits.md#technical) ».

| Droits sur l'interface       | Définition                                         |
|:---------------------------- |:-------------------------------------------------- |
| `Remediation`                | Accès au panneau d'administration des remédiations |
| `Remediation configuration`  | Accès aux configurations d'ordonnanceurs           |
| `Remediation instruction`    | Accès aux consignes                                |
| `Remediation job`            | Accès aux jobs d'ordonnanceurs                     |

**API**

Ces droits sont configurés dans le panneau de droits sous l'onglet « [API](./droits.md#api) ».

| Droits sur les API       | Définition                                          |
|:------------------------ |:--------------------------------------------------- |
| `Instructions`           | Manipulation des consignes                          |
| `Runs instructions`      | Exécuter une consigne                               |
| `File`                   | Manipulation des fichiers inclus dans les consignes |
| `Job configs`            | Manipulation des configurations des ordonnanceurs   |
| `Jobs`                   | Manipulation des jobs d'ordonnanceurs               |


### Créer une consigne

Pour créer une consigne, rendez-vous dans le menu d'administration de la
remédiation, onglet « CONSIGNES ».

![Ajout consigne](./img/consignes_creer_01.png)

Saisissez à présent les différentes étapes et opérations de votre consigne.
Voici un exemple :

**Nom et description de la consigne**

![Ajout consigne details1](./img/consignes_creer_02.png)

**Étape 1**

![Ajout consigne details2](./img/consignes_creer_03.png)

Puis, ajouter l'**Étape 2**

![Ajout consigne ajout étape](./img/consignes_creer_04.png)

![Ajout consigne details3](./img/consignes_creer_05.png)

!!! Note
    Veuillez noter que les textes associés aux opérations peuvent utiliser des
    [variables][templates-payload] au travers de template.

    Ainsi vous disposez principalement des variables `.Alarm` et `.Entity`.

Une fois créée, votre consigne sera affichée dans la liste des consignes.

![Liste consignes](./img/consignes_creer_06.png)

[templates-payload]: ../../templates-go/


| Type de consigne | Description |
| ---------------- | ----------- |
| Manuel           | L'exécution de la remédiation est à l'initiative du pilote à partir d'un bac à alarmes. Le système lui présente toutes les opérations à effectuer. Ces opérations peuvent inclure des [jobs](#taches) |
| Automatique      | L'exécution de la remédiation est déclenchée par un [trigger](#declenchement-dune-consigne-et-activation-dune-alarme). Le pilote ne peut que constater le résultat de la remédiation |
| Manuel simplifié | L'exécution de la remédiation est à l'initiative du pilote à partir d'un bac à alarmes. Ces remédiations sont uniquement une succession de jobs, sans opération manuelle à exécuter |

### Déclenchement d'une consigne et activation d'une alarme

L'option [ActivateAlarmAfterAutoRemediation](../../../guide-administration/administration-avancee/modification-canopsis-toml/#section-canopsisalarm) permet de décaler l'activation de l'alarme une fois la remédiation automatique terminée.

??? note "Schémas de fonctionnement"

    **"Option ActivateAlarmAfterAutoRemediation désactivée"**
    
    ```mermaid
    flowchart TD
        A[Création d'une alarme] --> B[snooze]
        A --> C[Début de comportement périodique]
        D --> G[unsnooze]
        D --> H[Fin de comportement périodique]
        A --> I[Remédiation automatique au moment de la création]
        B --> D[Ne pas activer l'alarme]
        C --> D
        G --> F{L'alarme est-elle snoozée ??\nOU\nen comportement périodique?\n}
        F -->|Oui| D
        F -->|Non| E[Activer l'alarme]
        H --> F
        I --> E
    ```

    **"Option ActivateAlarmAfterAutoRemediation activée"**
    
    ```mermaid
    flowchart TD
        A[Création d'une alarme] --> B[snooze]
        A --> C[Début de comportement périodique]
        D --> G[unsnooze]
        D --> H[Fin de comportement périodique]
        D --> J[Remédiation automatique au moment\nde la création terminée]
        A --> I[Remédiation automatique au moment de la créatio]
        B --> D[Ne pas activer l'alarme]
        C --> D
        G --> F{L'alarme est-elle snoozée ?\nOU\nen comportement périodique ?\nOU\nRemediation automatique au moment\nde la création en cours}
        F -->|Oui| D
        F -->|Non| E[Activer l'alarme]
        H --> F
        I --> D
        J --> F
    ```

Par ailleurs, les remédiations automatiques peuvent à présent [être déclenchées](../../../guide-administration/architecture-interne/triggers/#triggers-go) sur

* Création d'une alarme
* Activation d'une alarme
* Augmentation/Diminution de la sévérité d'une alarme
* Changement et verrouillage de sévérité
* Entrée ou sortie de comportement périodique


### Assigner une consigne à des alarmes

Lorsque votre consigne a été créée, vous devez l'assigner à une ou des alarmes.
Pour cela, vous allez pouvoir sélectionner ces alarmes grâce à des patterns
spécifiques de l'alarme ou de l'entité associée à l'alarme.
Utilisez pour cela le bouton d'action situé à droite de votre consigne.


![Assignation_consigne1](./img/consignes_tache_association_01.png)

Puis, associez vos alarmes en saisissant les patterns souhaités. Dans notre
exemple, il s'agit d'assigner la consigne aux alarmes dont la ressource contient
`ping`.



![Assignation_consigne2](./img/consignes_tache_association_02.png)

À ce stade, vous pouvez vérifier dans un bac à alarmes que les alarmes
sélectionnées par les patterns remplissent bien les conditions de votre consigne.

![Assignation_consigne3](./img/consignes_tache_association_03.png)

## Configuration

Cet onglet permet de gérer vos configurations d'ordonnanceurs.  
![consignes vue configuration](./img/consignes_vue_configuration.png)


### Ajouter un ordonnanceur

Pour être en mesure de relier un job à une opération, il est nécessaire de
définir une configuration d'ordonnanceur, ainsi que le job en lui-même.  
Pour cela, rendez-vous dans le menu « CONFIGURATIONS » du panneau
d'administration des remédiations.

En cliquant sur le « + » en bas à droite, vous accéderez au formulaire suivant :

![Ajout configuration](./img/consignes_configuration_ajout.png)

Explications sur les champs demandés :

* Nom : nom de la configuration qui sera utilisée dans la définition du job
* Type : `Rundeck`, `Awx`, `Jenkins`, ou encore `Visual Tom` fonction de votre ordonnanceur
* Hôte : adresse HTTP de votre ordonnanceur
* Jeton d'autorisation : jeton lié à votre utilisateur déclaré dans
l'ordonnanceur

!!! Note
    La configuration de cette liaison entre Canopsis et l'ordonnanceur est
    expliquée plus en détails dans le
    [guide d'administration de la remédiation][admin-remed].

## Tâches

### Créer une tâche
Lorsque la configuration d'ordonnanceur est prête, vous pouvez déclarer un
*job*.

Rendez-vous dans le menu « JOBS » du panneau d'administration des remédiations.

En cliquant sur le « + » en bas à droite, vous accéderez au formulaire suivant :

![Ajout jobs](./img/consignes_tache_modale_ajout.png)

Explications sur les champs demandés :

* Nom : nom du job qui sera utilisé dans une consigne
* Configuration : sélection d'une configuration précédemment créée
* Job ID : identifiant du job donné par l'ordonnanceur
* Payload : corps du message qui sera transmis à l'ordonnanceur au moment de
l'exécution du job

!!! Note
    L'association de job est illustrée dans le
    [guide d'administration de la remédiation][admin-remed].

#### Payload

Le *payload* (corps de message) associé à un job permet de variabiliser son
exécution et ainsi passer des paramètres utiles à l'ordonnanceur de tâches.

2 objets sont disponibles dans ce payload :

1. `.Alarm`
2. `.Entity`

L'utilisation des payloads dans les ordonnanceurs est développée dans le
[guide d'administration de la remédiation][admin-remed-payloads].


[remed-index]: ../remediation/index.md
[admin-remed]: ../../guide-administration/remediation/index.md
[admin-remed-payloads]: ../../guide-administration/remediation/index.md#utilisation-des-payloads

### Associer un job à une opération

!!! attention
    Pour qu'un job soit disponible, une configuration spécifique est nécessaire.
    [Consultez ce paragraphe pour cela](#ajouter-un-ordonnanceur).

Vous avez la possibilité d'associer un job d'ordonnanceur à une opération dans
une consigne.

Pour cela, dans votre consigne, il vous suffit de sélectionner les jobs qui
seront présentés à l'utilisateur de la consigne pour exécution.  
![Job1](./img/consignes_tache_association_01.png)

Pour confirmer que le job est bien associé, la petite pastille s'affiche (**1**).  
![Job2](./img/consignes_tache_association_02.png)

Au moment de l'exécution de la consigne, les jobs associés pourront être
exécutés.  
![Job3](./img/consignes_tache_association_03.png)


## Statistiques de remédiation

Cette vue simple, permet d'afficher les statistiques de traitement de vos consignes.  
![Statistiques vue graphique](./img/consignes_statistiques_liste.png)
