# Mise en oeuvre

## Introduction

!!! Note
    Disponible à partir de Canopsis 4.0.0

Cette documentation vous permet de mettre en oeuvre une remédiation de bout en bout.  
Comme expliqué en [préambule](./index.md), l'anatomie d'une remédiation suit le schéma ci-après : 

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

## Les droits

Le module de remédiation est soumis au système de droits sur l'interface et sur les APIs.

**Interface**

Ces droits sont configurés dans le panneau de droits sous l'onglet `technical`

| Droits sur l'interface       | Définition                                         |
|:---------------------------- |:-------------------------------------------------- |
| `Remediation`                | Accès au panneau d'administration des remédiations |
| `Remediation configuration`  | Accès aux configurations d'ordonnanceurs           |
| `Remediation instruction`    | Accès aux consignes                                |
| `Remediation job`            | Accès aux jobs d'ordonnanceurs                     |

**API**

Ces droits sont configurés dans le panneau de droits sous l'onglet `API`

| Droits sur les APIs      | Définition                                          |
|:--------------------- -- |:--------------------------------------------------- |
| `Instructions`           | Manipulation des consignes                          |
| `Runs instructions`      | Exécuter une consigne                               |
| `File`                   | Manipulation des fichiers inclus dans les consignes |
| `Job configs`            | Manipulation des configurations des ordonnanceurs   |
| `Jobs`                   | Manipulation des jobs d'ordonnanceurs               |

## Gestions des consignes

### Créer une consigne

Pour créer une consigne, rendez vous dans le menu d'administration de la remédiation, onglet `CONSIGNES`.  

![Ajout consigne](./img/remediation_instruction_ajout.png)

Saisissez à présent les différentes étapes et opérations de votre consigne. Voici un exemple 

**Nom et description de la consigne**

![Ajout consigne details1](./img/remediation_instruction_ajout_details1.png)

**Etape 1**

![Ajout consigne details2](./img/remediation_instruction_ajout_details2.png)

**Etape 2**

![Ajout consigne details3](./img/remediation_instruction_ajout_details3.png)

!!! Note
    Veuillez noter que les templates des opérations peuvent utiliser des [variables de payload](../../../guide-administration/architecture-interne/templates-golang/#templates-pour-payload)

    Ainsi vous disposez principalement des variables `.Alarm` et `.Entity`.

Une fois créée, votre consigne sera affichée dans la liste des consignes.  

![Liste consignes](./img/remediation_instruction_liste.png)

### Assigner une consigne à des alarmes

Lorsque votre consigne a été créée, vous devez l'assigner à une ou des alarmes. 
Pour cela, vous allez pouvoir sélectionner ces alarmes grâce à des patterns spécifiques de l'alarme ou de l'entité associée à l'alarme.
Utilisez pour cela le bouton d'action situé à droite de votre consigne.


![Assignation_consigne1](./img/remediation_instruction_assignation1.png)

Puis associer vos alarmes en saisissant les patterns souhaités, dans notre cas, les alarmes dont la ressource contient `ping`.

![Assignation_consigne2](./img/remediation_instruction_assignation2.png)

A ce stade, vous pouvez vérifier que les alarmes sélectionnées par les patterns sont bien éligibles à votre consigne.

![Assignation_consigne3](./img/remediation_instruction_assignation3.png)

### Associer un job à une opération

!!! Warning
    Pour qu'un job soit disponible, une configuration spécifique est nécessaire. [Consultez ce paragraphe pour cela](#jobs-associes-a-un-ordonnanceur)

Vous avez la possibilité d'associer un job d'ordonnanceur à une opération dans une consigne.

Pour cela, dans votre consigne, il vous suffit de sélectionner le ou les jobs qui seront présentés à l'utilisateur de la consigne pour exécution.

![Job1](./img/remediation_instruction_job1.png)

Au moment de l'exécution de la consigne, les jobs associés pourront être exécuter.

![Job2](./img/remediation_instruction_job2.png)

## Jobs associés à un ordonnanceur

Pour être en mesure de relier un job à une opération, il est nécessaire de définir une configuration d'ordonnanceur ainsi que le job en lui-même.  
Pour cela, RDV dans le menu `CONFIGURATIONS` du panneau d'administration des remédiations.  

En cliquant sur le "+" en bas à droite, vous accéderez au formulaire suivant : 

![Ajout configuration](./img/remediation_configuration_ajout.png)

Explications sur les champs demandés :

* Nom : Nom de la configuration qui sera utilisée dans la définition du job
* Type : `Rundeck` ou `Awx` fonction de votre ordonnanceur
* Hôte : Adresse http de votre ordonnanceur
* Jeton d'autorisation : Jeton lié à votre utilisateur déclaré dans l'ordonnanceur

!!! Note
    Vous pouvez consulter [cette page](../../../guide-administration/remediation/) qui concerne les configurations de Rundeck et Awx 

Lorsque la configuration d'ordonnanceur est prête, vous pouvez déclarer un `job`

RDV dans le menu `JOBS` du panneau d'administration des remédiations.

En cliquant sur le "+" en bas à droite, vous accéderez au formulaire suivant : 

![Ajout jobs](./img/remediation_job_ajout.png)

Explications sur les champs demandés :

* Nom : Nom du job qui sera utilisé dans une consigne
* Configuration : Sélection d'une configuration précedemment créée
* Job ID : Identifiant du job donné par l'ordonnanceur
* Payload : Corps du message qui sera transmis à l'ordonnaceur au moment de l'exécution du job

!!! Note
    Vous pouvez consulter [cette page](../../../guide-administration/remediation/) qui concerne les configurations de Rundeck et Awx 

### Payload

Le `payload` ou corps de message associé à un job permet de variabiliser son exécution et ainsi passer des paramètres à l'ordonnanceur de tâches.  
2 objets sont disponibles dans ce payload : 

1. `.Alarm`
2. `.Entity`
