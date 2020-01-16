# Personnalisation des typages

## Affichage de type de données formatées dans une colonne

Depuis la version `3.35.0`, il est maintenant possible de choisir un formatage de donnée particulier pour les champs présents et affichés dans une colonne du bac à alarmes de Canopsis.

Pour cela, il faut au préalable ajouter des règles de conversion via l'API `associativetable` qui sauvegardera dans un champ `alarm-column-filters`  de la collection MongoDB `default_associativetable`

### Envoi via l'API

Voici un exemple de configuration avec ajout de `2 règles`

```json
$ curl -H "Content-Type: application/json" -X POST -u <user>:<passwor> -d '{
    "filters" : 
    [
      {
        "column": "v.extra.activation_date",
        "filter": "date",
        "attributes": ["long"]
      },
      {
        "column": "v.duration",
        "filter": "duration"
      }
    ]
}' http://<ip_canopsis>:8082/api/v2/associativetable/alarm-column-filters
```

Dans cet exemple, nous allons configurer `2 valeurs de colonne` avec un rendu particulier.

Toute colonne qui afficherait :

* l'attribut `v.extra.activation_date` initialement de type `timestamp Unix` qui sera affiché en format `date` de type `long` ( voir [matrice de correspondance](#Les attributs liés aux dates peuvent contenir plusieurs type de formats) des types d'attributs liés au filtre date )
* l'attribut `v.duration` initialement en secondes qui sera affiché en format `duration` de type durée par défaut

### Les filtres utilisés ( attribut `filter` ) peuvent contenir plusieurs valeurs

* `date`
    1) Date
    2) Format
    3) Ignore la vérification du jour
    4) Valeur par défaut

* `duration`
    1) Valeur
    2) Locale (default - moment locale or constant)
    3) Durée (par défaut - `'D __ H _ m _ s _'`)

* `json` - Convertit des champs JSON en champs JSON avec indentation
    1) Chaîne Json
    2) Indentations (par défaut - `4`)
    3) Valeur par défaut (par défaut - `'{}'`)

* `percentage`
    1) Valeur
    2) Nombre de chiffres à virgule à garder (par défaut - `3`)

#### Les attributs liés aux dates peuvent contenir plusieurs type de formats

* `long` - DD/MM/YYYY H:mm:ss
* `medium` - DD/MM H:mm
* `short` - DD/MM/YYYY
* `time` - H:mm:ss
* `dateTimePicker`  -  DD/MM/YYYY HH:mm
* `dateTimePickerWithSeconds` - DD/MM/YYYY HH:mm:ss
* `datePicker` - DD/MM/YYYY
* `timePicker` - HH:mm
* `timePickerWithSeconds` - HH:mm:ss
* `veeValidateDateTimeFormat` - dd/MM/yyyy HH:mm



Il est à noter que ces champs utilisés ne peuvent être que des sous éléments d'une `alarme` et pas d'une `entité`

Par défaut, la date du jour ne sera pas affichée ( uniquement l'heure )

### Le fichier peut ensuite être envoyé via l'API pour charger la configuration

Attention à bien utiliser la méthode `POST` pour la première création de règle et `PUT` par la suite

```shell
curl -H "Content-Type: application/json" -X POST -u <user>:<passwor> -d @configuration.json http://<ip_canopsis>:8082/api/v2/associativetable/alarm-column-filters
```



### Accès à la configuration via l'UI

Nous allons configurer la personnalisation du champ `alarm.v.extra.activation_date` présent dans l'explorateur de contexte et correspondant au champ `v.extra.activation_date` envoyé à l'API [ci-dessus](#Envoi via l'API)

![](img/alarm-list-setting-3.png)



Pour accéder au paramétrage de cette fonctionnalité de rendu côté interface Web, il faut :

* Passer en mode édition sur le widget de bac à alarmes
* Se rendre dans le menu de configuration du widget du bac à alarmes
* Aller dans le sous menu `Column names`

![](img/alarm-list-setting-1.png)

* Ajouter une colonne correspondant à un des types personnalisés dont l'affichage doit être rendu grâce aux filtres. Puis sauvegarder tout en bas du menu.

![](img/alarm-list-setting-2.png)

* L'affichage du rendu se fait dans la colonne paramétrée ( par défaut la date du jour n'est pas affichée avec ce type de données )

![](img/alarm-list-setting-5.png)

* Elle s'affiche uniquement avec une date différente de la date du jour



![](img/alarm-list-setting-6.png)

* Voici le rendu de l'autre champ personnalisé envoyé à l'API ( `alarm.v.duration` utilisé au niveau de la configuration du widget )

![](img/alarm-list-setting-7.png)