# Guide Utilisateur

**TODO (DWU) :** ajouter d'autres exemples.

## Section : Gestion des alarmes / Les filtres

Lors de la configuration d'une vue, il est possible d'appliquer des filtres à notre liste.

Pour cela, dans l'onglet "Filters"

![img1](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/img1.png)  

Cliquez sur "add" pour créer votre premier filtre, une fenêtre apparaît :

![create_filter](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/create_filter.png)  

Donnez un titre à votre filtre, deux méthodes de création, une pour les néophytes et une pour les utilisateurs expérimentés. Nous allons nous concentrer sur la première méthode, 
vous verrez dans un second temps que la méthode avancée évolue en même temps que vos actions sur l'autre méthode.  

## AND / OR

Il faut maintenant choisir quel filtres prendre. Deux choix principaux s'offrent à vous sous forme d'opérateurs booléens **ET** et **OU** (AND et OR). En choisir un, puis appuyer sur "Add a rule".  

![adarule](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/adarule.png)  

Dans la première colonne, choisir de quel type sera l'objet filtré. Quatre choix :

* component_name
* connector_name
* connector
* ressource

Dans la seconde, le filtre qui lui sera alloué

* equal
* not equal
* in
* not in 
* begins with
* ....

Puis, dans la dernière, il vous est possible de remplir un champs qui sera "matché" avec le filtre.

Il vous est possible d'ajouter autant de filtre que vous souhaitez en cliquant sur "Add a rule".

Vous pouvez aussi séparer vos filtres en groupes. Simplement en cliquand sur "Add a group", et le supprimer en cliquant simplement sur "Delete groupe".

## Editeur avancé 

Je souhaite créer un filtre appelé "Mon premier filtre !" qui récupère les composants dont le nom est égal à "composants" :

Cela va ce fait trés simplement comme suit 

![example](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/exmpl1.png)  

Maintenant, allons voir l'éditeur avancé, vous remarquerez que l'on retrouve certains éléments du dernier screen :

```
{
    "$and": [
        {
            "component_name": "composants"
        }
    ]
}
```

Compliquons les choses ! 

Je veux maintenant ajouter le fait qu'une ressources ne doit pas être vide. Une fois la configuration réalisée via l'interface visuelle, on retrouve plusieur choses ajoutées à l'éditeur avancé :

```
{
    "$and": [
        {
            "component_name": "composants"
        },
        {
            "resource": {
                "$ne": ""
            }
        }
    ]
}
```

Il est bien évidement possible de réaliser cette configuration via l'interface utilisateur avancée. Le bouton "Parse" va vous servir à vérifier l'exactitude de votre JSON, si celui ci est invalide ce message apparaîtra :

![invjson](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/invjson.png)  

## Résultat

Une fois votre Filtre réalisé, il apparaîtra dans le menu déroulant "select a filter". 

![select_filter](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/select_filter.png)  

## Autre filtre

Si vous souhaitez limiter votre vue dans la durée, il vous suffit de cliquer sur l'icone en forme d'horloge à côté de "select a filter", une fenêtre s'ouvrira et vous pourrez alors choisir votre interval.  
 
![reporting](/doc-ce/Guide%20Utilisateur/Gestion%20des%20alarmes/Images/reporting.gif)  