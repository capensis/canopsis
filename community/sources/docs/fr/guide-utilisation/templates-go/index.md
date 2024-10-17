# Templates (Go)

Dans bon nombre de fonctionnalités de Canopsis, vous avez la possibilité d'utiliser des `Templates Go`. Ces templates permettent de générer du texte en fonction des variables associées aux événements, aux alarmes, ou encore aux entités.

Ces templates Go s'appuient sur le [package officiel GOLang](https://pkg.go.dev/text/template).

Canopsis définit en plus de cela des [fonctions](#fonctions-propres-a-canopsis) qui lui sont propres.

## Déclaration de variables

Les variables stockent des informations sur les alarmes, les événements et les entités. 

* `{{ .Alarm }}`  : propriétés d'une alarme
* `{{ .Event }}`  : propriétés d'un évènement 
* `{{ .Entity }}` : propriétés d'une entité

Ces trois éléments contiennent plusieurs propriétés que l'on peut utiliser pour créer des chaînes dynamiques.  
Par exemple, `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}` va créer la chaîne de caractères `http://mon-api.xyz/nom-du-component/nom-de-la-ressource`.

Voici une liste des principales variables mises à disposition. De façon générale, le champ original `alarm.v.nom_du_champ` sera transposé sous la forme `{{ .Alarm.Value.NomDuChamp }}`.

=== "Alarme"

    | Champs principaux                                  | Valeur                                                                |
    |:---------------------------------------------------|:--------------------------------------------------------------------- |
    | Display Name                                       | `{{ .Alarm.Value.DisplayName }}`                                      |
    | Composant                                          | `{{ .Alarm.Value.Component }}`                                        |
    | Ressource                                          | `{{ .Alarm.Value.Resource }}`                                         |
    | Type de connecteur                                 | `{{ .Alarm.Value.Connector }}`                                        |
    | Nom du Connecteur                                  | `{{ .Alarm.Value.ConnectorName }}`                                    |
    | Message initial de l'alarme                        | `{{ .Alarm.Value.InitialOutput }}`                                    |
    | Message de l'alarme                                | `{{ .Alarm.Value.Output }}`                                           |
    | Message initial Long de l'alarme                   | `{{ .Alarm.Value.InitialLongOutput }}`                                |
    | Message long de l'alarme                           | `{{ .Alarm.Value.LongOutput }}`                                       |
    | Message associé au dernier changement de criticité | `{{ .Alarm.Value.State.Message }}`                                    |
    | Criticité                                          | `{{ .Alarm.Value.State.Value }}`                                      |
    | Statut                                             | `{{ .Alarm.Value.Status.Value }}`                                     |
    | Auteur du ticket                                   | `{{ .Alarm.Value.Ticket.Author }}`                                    |
    | ID du ticket                                       | `{{ .Alarm.Value.Ticket.Ticket }}`                                    |
    | Message du ticket                                  | `{{ .Alarm.Value.Ticket.Message }}`                                   |
    | Auteur de l'acquittement                           | `{{ .Alarm.Value.ACK.Author }}`                                       |
    | Message de l'acquittement                          | `{{ .Alarm.Value.ACK.Message }}`                                      |
    | Auteur du dernier commentaire                      | `{{ .Alarm.Value.LastComment.Author }}`                               |
    | Message du dernier commentaire                     | `{{ .Alarm.Value.LastComment.Message }}`                              |
    | Informations enrichies depuis dynamic-infos        | `{{ (index (index .Alarm.Value.Infos "%rule_id%") "%infos_name%") }}` |
    | `abc` dans l'entité                                | `{{ (index .Entity.Infos "abc").Value }}`                             |


    | Champs de dates                             | Valeur                              |
    |:--------------------------------------------|:----------------------------------- |
    | Création de l'alarme                        | `{{ .Alarm.Value.CreationDate }}`   |
    | Activation de l'alarme                      | `{{ .Alarm.Value.ActivationDate }}` |
    | Dernier changement de sévérité de l'alarme  | `{{ .Alarm.Value.LastUpdateDate }}` |
    | Dernier événement reçu (Voir Option [EnableLastEventDate](../../../guide-administration/administration-avancee/modification-canopsis-toml/#section-canopsisalarm)                                     | `{{ .Alarm.Value.LastEventDate }}`  |
    | Résolution                                  | `{{ .Alarm.Value.Resolved }}`       |
 

    | Champs de durées                              | Valeur                                         |
    |:--------------------------------------------- |:---------------------------------------------- |
    | Inactivité                                    | `{{ .Alarm.Value.InactiveDuration }}`          |
    | Inactivité liée aux comportements périodiques | `{{ .Alarm.Value.PbehaviorInactiveDuration }}` |
    | Mise en veille (Mise à jour en fin de snooze) | `{{ .Alarm.Value.SnoozeDuration }}`            |

 
    | Champs de compteurs                                                        | Valeur                                             |
    |:-------------------------------------------------------------------------- |:-------------------------------------------------- |
    | Nombre total de changements de sévérité                                    | `{{ .Alarm.Value.TotalStateChanges }}`             |
    | Nombre de changements de sévérité depuis la dernière mise à jour de statut | `{{ .Alarm.Value.StateChangesSinceStatusUpdate }}` |
    | Nombre d'événements reçus                                                  | `{{ .Alarm.Value.EventsCount }}`                   |
 

=== "Entité"

    | Nom du champ                              | Valeur                                           |
    |:----------------------------------------- |:------------------------------------------------ |
    | ID                                        | `{{ .Entity.ID }}`                               |
    | Nom                                       | `{{ .Entity.Name }}`                             |
    | Composant                                 | `{{ .Alarm.Value.Component }}`                   |
    | Connector                                 | `{{ .Alarm.Value.Connector }}`                   |
    | `abc` dans les informations du composant  | `{{ (index .Entity.ComponentInfos "abc").Value }}` |
    | `abc` dans les informations de l'entité   | `{{ (index .Entity.Infos "abc").Value }}`          |
 

=== "Evénement"

    Pour les champs de date, comme par exemple `{{ .Event.Timestamp }}`, il est possible de récupérer l'information de différents manières.
    
    | Champ                             | Résultat                                                    |
    |:--------------------------------- |:----------------------------------------------------------- |
    | Type du connecteur                                | `{{ .Event.Connector }}`                    |
    | Nom du connecteur                                 | `{{ .Event.ConnectorName }}`                |
    | Composant                                         | `{{ .Event.Component }}`                    |
    | Ressource                                         | `{{ .Event.Resource }}`                     |
    | Type d'événement                                  | `{{ .Event.EventType }}`                    |
    | Message de l'événement                            | `{{ .Event.Output }}`                       |
    | Auteur de l'événement                             | `{{ .Event.Author }}`                       |
    | `abc` dans les extra informations de l'événement  | `{{ index .Event.ExtraInfos "abc" }}` |

=== "Environnement"

    Vous pouvez déclarer des variables d'environnement dans le fichier [Canopsis.toml](../../guide-administration/administration-avancee/modification-canopsis-toml#section-canopsistemplatevars). Ces variables sont accessibles dans les templates Go de la manière suivante.

    | Champ                    | Résultat            |
    |:------------------------ |:------------------- |
    | Variable `projet`        | `{{ .Env.projet }}` |
 

Dans les fonctionnalités [Scénario](../menu-exploitation/scenarios/) et [Déclaration de tickets](../menu-exploitation/regles-declaration-tickets/), vous pouvez accéder à certaines variables supplémentaires.

=== "Additional Data"

    Vous avez également la possibilité de récupérer des informations propres à l'action qui s'exécute.
    
    | Champ                                   | Résultat                                                                                                  |
    |:--------------------------------------- |:--------------------------------------------------------------------------------------------------------- |
    | `{{ .AdditionalData.RuleName }}`        | Nom de la règle                                                                                           |
    | `{{ .AdditionalData.AlarmChangeType }}` | Nom du trigger (sous forme de chaîne : ack, stateinc, etc.) |                                             |
    | `{{ .AdditionalData.Author }}`          | Auteur de l'action                                                                                        |
    | `{{ .AdditionalData.Initiator }}`       | Initiateur de l'action (**user** pour une action exécutée depuis l'interface graphique, **system** pour une action exécutée par un moteur) |
    | `{{ .AdditionalData.Output }}`          | Message de l'événement                                                                                    |


=== "Children"

    Lorsque l'alarme à laquelle le webhook est confronté est une [méta alarme](../menu-exploitation/regles-metaalarme/), il est possible de parcourir les alarmes conséquences pour en extraire le contenu.  
    Pour cela, un opérateur `range` permet d'itérer sur la variable `.Children` qui contient l'ensemble des alarmes conséquences de la méta alarme.
    
    La syntaxe à utiliser est la suivante :
    
    ```
    {{ range $variable := .Children }} ... {{ end }}
    ```
    
    Voici un exemple concret d'utilisation de cette variable dans un payload de Webhook :
    
    ```
    {
      "message": "{{ range $children := .Children }}{{ $children.ID }} - {{ $children.Value.State.Message }}\n{{ end }}"
    }
    ```
    
    Le payload de ce webhook sera donc constitué d'un attribut `message` dont la valeur sera une succession de lignes contenant l'identifiant et le message des alarmes conséquences séparés par un "-".
    
    ```
    {
      "message": "23818029-b80d-416e-9d12-5963c76bcbfa - message alarme conséquence 1\n6594ddea-9fd7-4423-a2db-ba10b72c67aa - message alarme conséquence 2\n"
    }
    ```

## Fonctions propres à Canopsis

En plus des [fonctions de base](https://pkg.go.dev/text/template) pour tester la valeur des variables, il existe plusieurs fonctions pour transformer le contenu de la variable.

Pour les utiliser, il faut appeler la fonction après la variable comme ceci : `{{ .LaVariable | fonction }}` ou `{{ .LaVariable | fonction param }}` si la fonction a besoin d'autres paramètres.

On peut aussi enchaîner différentes fonctions à la suite si on veut transformer plusieurs fois les variables `{{ .LaVariable | fonction1 | fonction2 paramA paramB | fonction3 paramC }}`.

| Fonction          | Description                                                                                                              | Exemples                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | 
|-------------------|--------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **json**          | Encode une variable en JSON                                                                                              | Si `.Entity.Infos.info1.Value = ["a","b","c"]` et `.Entity.Infos.info2.Value = "d \"e\""`<br><ul><li>`{{ .Entity.Infos.info1.Value                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | json }}` -> `["a","b","c"]`</li><li>`{{ .Entity.Infos.info2.Value | json }}` -> `"d \"e\""`</li><li>`{{ .Entity.Infos | json }}` -> `{"info1":{"name":"info1","description":"","value":["a","b","c"]},"info2":{"name":"info2","description":"","value":"d \"e\""}}`</li></ul><br>    En complément :<br><ul><li>`{{ .Entity.Infos.info1.Value }}` -> `[a b c]`</li><li>`{{ .Entity.Infos.info2.Value }}` -> `d "e"`<</li><li>`{{ .Entity.Infos }}` -> `map[info1:{info1  [a b c]} info2:{info2  d "e"}]`</li></ul> |
| **json_unquote**  | Encode une variable en JSON et supprime les guillemets                                                                   | Si `.Entity.Infos.info1.Value = ["a","b","c"]` et `.Entity.Infos.info2.Value = "d \"e\""`<br><ul><li>`{{ .Entity.Infos.info1.Value                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | json_unquote }}` -> `["a","b","c"]`</li><li>`{{ .Entity.Infos.info2.Value | json_unquote }}` -> `d \"e\"`</li><li>`{{ .Entity.Infos | json_unquote }}` -> `{"info1":{"name":"info1","description":"","value":["a","b","c"]},"info2":{"name":"info2","description":"","value":"d \"e\""}}`</li></ul> |
| **split**         | Split une chaîne grâce à un séparateur                                                                                   | Si `.Alarm.Value.Output = "a/b c"`<br><ul><li>`{{ .Alarm.Value.Output                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           | split "/" 0 }}` -> `a`</li><li>`{{ .Alarm.Value.Output | split "/" 1 }}` -> `b c`</li><li>`{{ .Alarm.Value.Output | split " " 1 }}` -> `c`</li></ul> |
| **trim**          | Supprime les espaces en début et fin de chaîne                                                                           | Si `.Alarm.Value.Output = "  a b c      "`<br><ul><li>`{{ .Alarm.Value.Output                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | trim }}` -> `a b c`</li></ul> |
| **replace**       | Remplace les correspondances d'une expression régulière dans une chaîne par une chaîne                                   | Si `.Alarm.Value.Output = "abc b 10"`<br><ul><li>`{{ .Alarm.Value.Output                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | replace "b" "d" }}` -> `adc d 10`</li><li>`{{ .Alarm.Value.Output | replace "\\d+" "20" }}` -> `abc b 20`</li><li>`{{ .Alarm.Value.Output | replace "(\\d+)" "$1 out of 100" }}` -> `abc b 10 out 100`</li></ul> |
| **uppercase**     | Transforme toutes les lettres en majuscule                                                                               | Si `.Alarm.Value.Output = "a b c ô é"`<br><ul><li>`{{ .Alarm.Value.Output                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       | uppercase }}` -> `A B C Ô É`</li></ul> |
| **lowercase**     | Transforme toutes les lettres en minuscule                                                                               | Si `.Alarm.Value.Output = "A B C Ô É"`<br><ul><li>`{{ .Alarm.Value.Output                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       | lowercase }}` -> `a b c ô é`</li></ul> |
| **formattedDate** | Formatte la date en UTC (déprécié, utilisez `localtime` à la place)                                                      | Si `.Alarm.Timestamp = 1635404700`<br><ul><li>`{{ .Alarm.Timestamp                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | formattedDate "Mon, 02 Jan 2006 15:04:05 MST" }}` -> `Thu, 28 Oct 2021 07:05:00 UTC`</li></ul> |
| **localtime**     | Formatte la date dans la timezone locale ou dans une timezone définie                                                    | Si `.Alarm.Timestamp = 1635404700`<br><ul><li>`{{ .Alarm.Timestamp                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | localtime "Mon, 02 Jan 2006 15:04:05 MST" }}` -> `Thu, 28 Oct 2021 09:05:00 CEST`</li><li>`{{ .Alarm.Timestamp | localtime "Mon, 02 Jan 2006 15:04:05 MST" "Australia/Queensland" }}` -> `Thu, 28 Oct 2021 17:05:00 AEST`</li></ul><br>Cette fonction prend en paramètre une chaîne qui est le format attendu de la date. La chaîne doit correspondre à la syntaxe des dates en Go. Cette syntaxe se base sur une date de référence, le `01/02 03:04:05PM '06 -0700` qui correspond au lundi 2 janvier 2006 à 22:04:05 UTC. Quand la chaîne n'arrive pas à être analysée par le langage, elle est renvoyée telle quelle. |
| **regex_map_key** | Extrait un item d'une map via une expression régulière                                                                   | Si `.Event.ExtraInfos = {"info":"a","anotherinfo":"b"}`<br><ul><li>`{{ regex_map_key .Event.ExtraInfos "info" }}` -> `a` ou `b` car l'ordre dans une map est indéterminé</li><li>`{{ regex_map_key .Event.ExtraInfos "^info" }}` -> `a`</li></ul>                                                                                                                                                                                                                                                                                                                                                               |
| **map_has_key**   | Vérifie la présence d'un item par sa clé dans une map                                                                    | Si `.Event.ExtraInfos = {"info":"a"}`<br><ul><li>`{{ if map_has_key .Event.ExtraInfos "info" }}{{ .Event.ExtraInfos.info }}{{ else }}default{{ end }}` -> `a`</li><li>`{{ if map_has_key .Event.ExtraInfos "anotherinfo" }}{{ .Event.ExtraInfos.anotherinfo }}{{ else }}default{{ end }}` -> `default`</li></ul><br>En complément :<br><ul><li>`{{ index .Event.ExtraInfos "info" }}` -> `a`</li><li>`{{ index .Event.ExtraInfos "anotherinfo" }}` -> `<no value>`</li><li>`{{ .Event.ExtraInfos.info }}` -> `a`</li><li> `{{ .Event.ExtraInfos.anotherinfo }}` -> Erreur de compilation du template.</li></ul> |
| **tag_has_key**   | Vérifie si un tag est présent dans une liste de tags                                                                     | Si `.Alarm.Tags = ["Tag1: Value1", "Tag2"]`<br><ul><li>`{{ tag_has_key .Alarm.Tags "Tag1" }}` -> `true`</li><li>`{{ tag_has_key .Alarm.Tags "Tag2" }}` -> `true`</li><li>`{{ tag_has_key .Alarm.Tags "Tag3" }}` -> `false`</li></ul>                                                                                                                                                                                                                                                                                                                                                                            |
| **get_tag**       | Renvoie la valeur d'un tag. Si le tag n'existe pas ou n'a pas de valeur, une chaine vide est renvoyée                    | Si `.Alarm.Tags = ["Tag1: Value1", "Tag2"]`<br><ul><li>`{{ get_tag .Alarm.Tags "Tag1" }}` -> `Value1`</li><li>`{{ get_tag .Alarm.Tags "Tag2" }}` -> `""`</li><li>`{{ get_tag .Alarm.Tags "Tag3" }}` -> `""`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                           |
| **substrLeft**    | Format: `{{ substrLeft .Var N }}` <br> Renvoie N caractères depuis le début de la chaine.                                | Si `.Alarm.Value.Output = "test-string"` <br><ul><li>`{{ substrLeft .Alarm.Value.Output 4 }}` -> `test`</li><li>`{{ substrLeft .Alarm.Value.Output 100 }}` -> `"test-string"`</li><li>`{{ substrLeft .Alarm.Value.Output -1 }}` -> `""`</li></ul>                                                                                                                                                                                                                                                                                                                                                               |
| **substrRight**   | Format: `{{ substrRight .Var N }}` <br> Renvoie N caractères depuis la fin de la chaine.                                 | Si `.Alarm.Value.Output = "test-string"` <br><ul><li>`{{ substrRight .Alarm.Value.Output 6 }}` -> `string`</li><li>`{{ substrRight .Alarm.Value.Output 100 }}` -> `"test-string"`</li><li>`{{ substrRight .Alarm.Value.Output -1 }}` -> `""`</li></ul>                                                                                                                                                                                                                                                                                                                                                          |
| **substr**        | Format: `{{ substr .Var M N }}` <br> Renvoie N caractères depuis la position M.                                          | Si `.Alarm.Value.Output = "test-string"` <br><ul><li>`{{ substr .Alarm.Value.Output 2 5 }}` -> `st-st`</li><li>`{{ substr .Alarm.Value.Output 2 100 }}` -> `"st-string"`</li><li>`{{ substrRight .Alarm.Value.Output 2 -1 }}` -> `""`</li></ul>                                                                                                                                                                                                                                                                                                                                                                 |
| **strlen**        | Format: `{{ strlen .Var }}` <br> Renvoie la longueur de la chaine.                                                       | Si `.Alarm.Value.Output = "test-string"` <br><ul><li>`{{ strlen .Alarm.Value.Output }}` -> `11`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| **strpos**        | Format: `{{ strpos .Str .Substr }}` <br> Renvoie l'index de la première occurrence de Substr dans Str. Renvoie -1 sinon. | Si `.Alarm.Value.Output = "test-string"` <br><ul><li>`{{ strpos .Alarm.Value.Output "-" }}` -> `4`</li><li>`{{ strpos .Alarm.Value.Output "+" }}` -> `-1`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| **add**           | Format: `{{ add .X .Y }}` <br> Renvoie la somme de X et Y.                                                               | <ul><li>`{{ add 2 2 }}` -> `4`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| **sub**           | Format: `{{ sub .X .Y }}` <br> Renvoie la différence entre X et Y.                                                       | <ul><li>`{{ sub 4 2 }}` -> `2`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| **mult**          | Format: `{{ mult .X .Y }}` <br> Renvoie le produit de X et Y.                                                            | <ul><li>`{{ mult 2 2 }}` -> `4`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| **div**           | Format: `{{ div .X .Y }}` <br> Renvoie le résultat de la division de X par Y.                                            | <ul><li>`{{ div 6 2 }}` -> `3`</li><li>`{{ div 7 2 }}` -> `3`</li></ul>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |

??? Note "Le tableau ci-dessous montre quelques directives qui sont reconnues, ainsi que leur correspondance avec la fonction `date` dans les systèmes UNIX."

    | Directive pour les templates | Correspondance UNIX ([date](http://www.linux-france.org/article/man-fr/man1/date-1.html)) | Définition                        | Exemples            |
    |:---------------------------- |:--------------- |:--------------------------------- |:------------------- |
    | `Mon`                        | `%a`            | Abréviation du jour de la semaine | Mon..Sun            |
    | `Monday`                     | `%A`            | Nom du jour de la semaine         | Monday..Sunday      |
    | `Jan`                        | `%b`            | Abréviation du nom du mois        | Jan..Dec            |
    | `January`                    | `%B`            | Nom du mois                       | January..December   |
    | `02`                         | `%d`            | Jour du mois                      | 01..31              |
    | `15`                         | `%k`            | Heure (sur 24 heures)             | 0..23               |
    | `01`                         | `%m`            | Mois                              | 01..12              |
    | `04`                         | `%M`            | Minute                            | 01..59              |
    | `05`                         | `%S`            | Seconde                           | 01..61              |
    | `2006`                       | `%Y`            | Année                             | 1970, 1984, 2019… |
    | `MST`                        | `%Z`            | Fuseau horaire                    | CEST, EDT, JST…   |
    
    Ainsi, pour transformer un champ en une date au format `heure:minute:seconde`, il faudra utiliser `localtime \"15:04:05\"` (même si le champ dans l'alarme ou l'événement ne correspondent pas à cette heure).
    
    La [documentation officielle de Go](https://golang.org/pkg/time/#pkg-constants) fournit par ailleurs les valeurs à utiliser pour des formats de dates standards. Pour obtenir une date suivant le RFC3339, il faudra utiliser `localtime \"2006-01-02T15:04:05Z07:00\"`. De même, `localtime \"02 Jan 06 15:04 MST\"` sera appelé pour générer une date au format RFC822.

## Fonctions incluses dans GO


| Fonction | Description | Exemples | 
| -------- | ----------- | -------- |
| **range** | Permet d'itérer sur une variable | Si `.Entity.Infos.info1.Value = ["a","b","c"]` et `.Entity.Infos.info2.Value = "d"`<br><ul><li>`{{ range (index .Entity.Infos "info1").Value }}{{ . }};{{ end }}` -> `a;b;c;`</li><li>`{{ range (index .Entity.Infos "not-exist").Value }}{{ . }};{{ end }}` -> `empty string`</li><li>`{{ range .Entity.Infos }}{{ .Name }}:{{ .Value }};{{ end }}` -> `info1:[a b c];info2:d;`</li><li>`{{ range $index, $item := (index .Entity.Infos "info1").Value }}{{ $index }}:{{ $item }};{{ end }}` -> `0:a;1:b;2:c;`</li></ul> |
| **urlquery** | Transforme le contenu de la variable en une chaîne de caractères compatible avec le format des URL | Si `.Alarm.Value.ticket.Ticket = 50`<br><ul><li>`http://une-api.org/edit/{{ .Alarm.Value.Ticket.Value | urlquery }}` -> `http://une-api.org/edit/50`</li></ul> |
| **printf** | Affiche un contenu à partir d'un template | Si `.Entity.Infos.info1.Value = "value1"` et `.Entity.Infos.info2.Value = "value2"`<br><ul><li>`{{ $description := printf "%s -- %s" .Entity.Infos.info1.Value .EntityInfos.info2.Value }}` -> `"value1 -- value2"`</li></ul> |


## Exemples

Une fois qu'on possède les variables nécessaires, la seconde étape est la génération du texte.

Pour mieux comprendre comment fonctionne la génération et comment utiliser les variables, plusieurs exemples seront présentés.

### Génération simple

Dans un premier cas, on peut utiliser directement les variables pour générer le texte. En reprenant le template `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}`, on peut voir :

*  la première variable `$comp` définie comme `.Alarm.Value.Component` ;
*  la seconde variable `$res` qui a pour valeur `.Alarm.Value.Resource` ;
*  enfin le texte lui-même `http://mon-api.xyz/{{$comp}}/{{$res}}` qui va donner `http://mon-api.xyz/nom-du-component/nom-de-la-ressource` après transformation.

!!! note
    Dans cet exemple simple, on aurait pu se passer de variables et utiliser directement `http://mon-api.xyz/{{ .Alarm.Value.Component }}/{{ .Alarm.Value.Resource }}`.

### Génération selon les variables

Dans une utilisation plus avancée, on peut  générer du texte en fonction de l'état de la variable. Dans le cas suivant `"{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}"`, on trouve :

*  la déclaration de variable `$val`, `{{ $val := .Alarm.Value.Status.Value }}` ;
*  le texte `http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}` qui contient une condition sur la variable `$val`.

La condition `{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}` va regarder si `$val` vaut soit 0, soit 2, soit 4. Si c'est le cas, on obtiendra `http://127.0.0.1:5000/even`, `http://127.0.0.1:5000/odd` sinon.

Ici, nous avons utilisé le `or` et le `eq`, mais il est possible de tester les conditions avec `and`, `not`, `ne` (not equal), `lt` (less than), `le` (less than or equal), `gt` (greater than) ou `ge` (greater than or equal).

### Avec variables et fonctions

```json
{
    "payload" : "{{ $comp := .Alarm.Value.Component }}{{ $reso := .Alarm.Value.Resource }}{{ $val := .Alarm.Value.Status.Value }}{\"component\": \"{{$comp}}\",\"resource\": \"{{$reso}}\", \"parity\": {{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}\"even\"{{else}}\"odd\"{{end}},  \"value\": {{$val}} }"
}
```

On définit trois variables que sont `$comp`, `$reso` et `$val` puis on complète le champ `parity` du payload en regardant la valeur de `$val`. Dans le cas où `$val` vaut 2, alors le payload envoyé sera :

```json
{
    "component":"127.0.0.1",
    "resource":"HDD",
    "parity": "even",
    "value": 2
}
```
