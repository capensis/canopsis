# Moteur de recherche

Actullement, il y a deux types de recherches disponibles dans le bac à alarmes :
  * la recherche naturelle
  * la recherche avancée

![Champ de recherche](./img/champ-recherche.png "Champ de recherche")  

## La recherche naturelle

La recherche naturelle consiste en la recherche d'une chaine de caractères,
saisie dans le champ de recherche, sur toutes les colonnes affichées dans le
bac à alarmes.

## La recherche avancée

Vous retrouverez un résumé de cette partie en cliquant sur le point d'interrogation.

La recherche avancée est une recherche qui permet de rechercher des alarmes en
fonction de la valeur d'une ou plusieurs de ces colonnes défini explicitement
contrairement à la recherche naturelle.

### Description da la grammaire

La grammaire est constituée d'une ou plusieurs conditions séparées par des
opérateurs logiques.

Pour faire une recherche avec grammaire, il faut absolument la préfixer par
 "- ".

### Les conditions

Une condition est constituée de deux opérandes séparées par un opérateur de
comparaison. Cette condition peut être précédée d'inverseur qui va inverser le
résultat de la condition.

L'opérande de gauche correspond au nom de la colonne dans laquelle
rechercher la valeur contenue dans l'opérande de droite.

#### Le nom de colonne

Le nom de colonne est une chaine de caractères alpha-numériques correspondant
aux noms des colonnes affichées dans le bac à alarmes.

#### Les types de valeur

La valeur peut prendre plusieurs formes :
  - une chaine de caractères alpha-numérique entre guillemets
  - un booléen (**"TRUE"**, **"FALSE"**)
  - un entier
  - un nombre flottant
  - ou **"NULL"**

#### Les opérateurs de comparaison

Il existe 8 opérateurs de comparaison :

  * "**<=**" pour sélectionner des alarmes dont la valeur numérique est inférieure
  ou égal à l'opérande de droite ;
  * "**<**" pour sélectionner des alarmes dont la valeur numérique est strictement
  inférieure à l'opérande de droite ;
  * "**=**" pour sélectionner des alarmes dont la valeur est égale à l'opérande de
  droite ;
  * "**!=**" pour sélectionner des alarmes dont la valeur est différente de
  l'opérande de droite ;
  * "**>=**" pour sélectionner des alarmes dont la valeur numérique est supérieure
  ou égal à l'opérande de droite ;
  * "**>**" pour sélectionner des alarmes dont la valeur numérique est strictement
  supérieure à l'opérande de droite ;
  * "**LIKE**" pour rechercher des alarmes dont la chaine de caractères correspond
  à l'expression régulière MongoDB.


#### Les opérateurs logiques

Il existe 3 opérateurs booléens :

  * **AND** qui permet de réaliser un ET logique entre deux conditions ;
  * **OR** qui permet de réaliser un OU logique entre deux conditions ;
  * **NOT** qui permet d'inverser le résultat d'une condition.


### Exemples d'utilisation

  * ```- Connector = "connector_1"``` : pour rechercher toutes les alarmes
  dont le connecteur est "connector_1" ;
  * ```- Connector="connector_1" AND Resource="resource_3"``` : pour rechercher
  toutes les alarmes dont le connecteur est "connector_1" et la ressource est
  "resource_3" ;
  * ```- Connector="connector_1" OR Resource="resource_3"``` : pour rechercher
  toutes les alarmes dont le connecteur est "connector_1" ou la ressource est
  "resource_3" ;
  * ```- Connector LIKE 1 OR Connector LIKE 2``` : pour rechercher toutes les
  alarmes dont le connector contient un 1 ou toutes les alarmes dont le
  connector contient un 2 ;
  * ```- NOT Connector = "connector_1"``` : pour rechercher toutes les alarmes
  dont le connecteur n'est pas "connector_1".

### Liste de requêtes utiles

  * Recherche d'alarmes à partir d'un numéro de ticket :
  ```- ticket.val = "123456"```
