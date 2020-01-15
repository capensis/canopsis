# Templates (Go)

Dans les [Webhooks](../moteurs/moteur-webhook.md), les champs `payload` et `url` sont personnalisables grâce aux templates Go. Les templates permettent de générer du texte en fonction de l'état de l'alarme, de l'évènement ou de l'entité.

Les templates des champs `payload` et `url` peuvent se décomposer en deux parties : la déclaration de variables et le corps du texte lui-même.

La déclaration de variables doit être positionnée avant le corps du message. Les variables se distinguent du corps du message par le fait qu'elles sont entourées de doubles accolades.

Pour plus d'informations, vous pouvez consulter la [documentation officielle de Go sur les templates](https://golang.org/pkg/text/template).

## Déclaration de variables

!!! note
    La déclaration n'est pas obligatoire mais elle est recommandée si beaucoup de variables seront utilisées pour générer du texte.

Les variables stockent des informations sur les alarmes, les événements et les entités. `{{ .Alarm }}` permet d'accéder aux propriétés d'une alarme, de même que `{{ .Event }}` pour un évènement et `{{ .Entity }}` pour une entité. Ces trois éléments contiennent plusieurs propriétés qu'on peut utiliser pour créer des chaînes dynamiques. Par exemple, `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}` va créer la chaîne de caractères `http://mon-api.xyz/nom-du-component/nom-de-la-ressource`.

Voici une liste des principales données et la manière de la récupérer.

| Nom du champ        | Valeur                              |
|:------------------- |:----------------------------------- |
| Composant           | `{{ .Alarm.Value.Component }}`      |
| Ressource           | `{{ .Alarm.Value.Resource }}`       |
| Message             | `{{ .Alarm.Value.State.Message }}`  |
| Statut              | `{{ .Alarm.Value.Status.Value }}`   |
| Etat                | `{{ .Alarm.Value.State.Value }}`    |
| Auteur du ticket    | `{{ .Alarm.Value.Ticket.Author }}`  |
| Numéro du ticket    | `{{ .Alarm.Value.Ticket.Value }}`   |
| Message du ticket   | `{{ .Alarm.Value.Ticket.Message }}` |
| Auteur de l'ACK     | `{{ .Alarm.Value.ACK.Author }}`     |
| Message de l'ACK    | `{{ .Alarm.Value.ACK.Message }}`    |
| `abc` dans l'entité | `{{ .Entity.Infos.abc.Value }}`     |

Pour les champs de date, comme par exemple `{{ .Event.Timestamp }}`, il est possible de récupérer l'information de différents manières.

| Champ                             | Résultat                                                                  |
|:--------------------------------- |:------------------------------------------------------------------------- |
| `{{ .Event.Timestamp.Day }}`      | Jour (sous forme d'entier)                                                |
| `{{ .Event.Timestamp.Minute }}`   | Minutes (sous forme d'entier)                                             |
| `{{ .Event.Timestamp.Second }}`   | Secondes (sous forme d'entier)                                            |
| `{{ .Event.Timestamp.String }}`   | Chaîne de caractères représentant la date suivant un formatage par défaut |
| `{{ .Event.Timestamp.Unix }}`     | Timestamp UNIX                                                            |
| `{{ .Event.Timestamp.UnixNano }}` | Timestamp UNIX en nanosecondes                                            |

!!! attention
    Les champs enrichis depuis un événement ou via l'event filter se retrouvent au niveau de l'entité et sont sensibles à la casse. Par exemple un champ enrichi intitulé `switch` dans l'entité sera traduit en `{{ .Entity.Infos.switch.Value }}`.

## Génération de texte

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

### Transformation des variables

En plus des fonctions de base pour tester la valeur des variables, il existe plusieurs fonctions pour transformer le contenu de la variable.

Pour les utiliser, il faut appeler la fonction après la variable comme ceci : `{{ .LaVariable | fonction }}` ou `{{ .LaVariable | fonction param }}` si la fonction a besoin d'autres paramètres.

On peut aussi enchaîner différentes fonctions à la suite si on veut transformer plusieurs fois les variables `{{ .LaVariable | fonction1 | fonction2 paramA paramB | fonction3 paramC }}`.

#### `urlquery`

`urlquery` va transformer le contenu de la variable en une chaîne de caractères compatible avec le format des URL. Cette fonction a son intérêt si l'adresse du service externe dépend de l'état de l'alarme ou du ticket et que le contenu contient des caractères spéciaux. Un exemple d'adresse serait `http://une-api.org/edit/{{ .Alarm.Value.Ticket.Value | urlquery }}` pour modifier un ticket déjà existant.

#### Fonctionnalités spécifiques à Canopsis

Certaines fonctionnalités ne sont pas présentes de base dans les templates Go. Elles ont été implémentées par l'équipe de Canopsis.

!!! note
    Les fonctions suivantes sont disponibles dans les templates des [webhooks](../moteurs/moteur-webhook.md), pas ceux de l'event-filter.

##### `json` et `json_unquote`

La fonction `json`, ainsi que sa variante `json_unquote` vont transformer un champ en un élément directement insérable dans un fichier JSON et en conservant le contenu exact. Les caractères spéciaux dans les chaînes de caractères seront envoyés en tant que tel, sans traitement pas les engines Go.

`json_unquote` va retirer les guillemets si le résultat est une chaîne de caractères, sinon il retourne la même chose que `json`.

Ainsi, `{{ .Alarm.Value.ACK.Message | json }}` va renvoyer la chaîne `"ACK by someone"`, directement avec les guillemets et donc insérable dans un JSON.

`json_unquote` va lui générer la chaîne `ACK by someone` qui n'est pas utilisable directement dans un JSON mais qui peut servir pour créer des textes plus complexes. `\"Voici un message : '{{ .Alarm.Value.ACK.Message | json_unquote }}'\"` va donner par exemple `"Voici un message : 'ACK by someone'"`.

##### `split`

`split` va diviser une chaîne de caractères en plusieurs sous-chaînes selon un séparateur et retourner une de ces sous-chaînes à partir d'un indice.

Si par exemple l'output d'un événement vaut `"SERVER#69420#DOWN"`, `{{ .Event.Output | split \"#\" 2 }}` va renvoyer la chaîne `DOWN`. Comme les indices commencent à 0, `SERVER` a pour indice 0, `69420` a pour indice 1 et `DOWN` a pour indice 2.

##### `trim`

La fonction `trim` permet de supprimer les blancs en début et fin de chaîne de caractères. Les blancs pris en compte sont ceux définis par Unicode et ils comprennent l'espace, la tabulation, l'espace insécable ainsi que les caractères de fin de ligne.

##### `replace`

`replace` prend en paramètre une expression régulière (ou regex) et une chaîne de caractères. Cette fonction va remplacer toutes les occurrences de la regex par la chaîne.

Par exemple `{{ .Event.Output | replace \"\\r?\\n\" \"\"  }}` possède pour paramètre l'expression régulière `\r?\n` et la chaîne vide. Cela va supprimer tous les caractères de fin de ligne de l'output de l'événement.

##### `formattedDate`

`formattedDate` est la fonction qui va transformer les dates en chaînes de caractères, suivant la syntaxe Golang. Elle ne fonctionne que sur les champs qui sont des `CpsTime`, comme par exemple `.Alarm.Value.CreationDate` ou `.Event.Timestamp`.

Cette fonction prend en paramètre une chaîne qui est le format attendu de la date. La chaîne doit correspondre à la syntaxe des dates en Go. Cette syntaxe se base sur une date de référence, le `01/02 03:04:05PM '06 -0700` qui correspond au lundi 2 janvier 2006 à 22:04:05 UTC. Quand la chaîne n'arrive pas à être analysée par le langage, elle est renvoyée telle quelle.

 Le tableau ci-dessous montre quelques directives qui sont reconnues, ainsi que leur correspondance avec la fonction `date` dans les systèmes UNIX.

| Directive pour les templates | Correspondance UNIX ([date](http://www.linux-france.org/article/man-fr/man1/date-1.html)) | Définition                        | Exemples            |
|:---------------------------- |:----------------------------------------------------------------------------------------- |:--------------------------------- |:------------------- |
| `Mon`                        | `%a`                                                                                      | Abréviation du jour de la semaine | Mon..Sun            |
| `Monday`                     | `%A`                                                                                      | Nom du jour de la semaine         | Monday..Sunday      |
| `Jan`                        | `%b`                                                                                      | Abréviation du nom du mois        | Jan..Dec            |
| `January`                    | `%B`                                                                                      | Nom du mois                       | January..December   |
| `01`                         | `%d`                                                                                      | Jour du mois                      | 01..31              |
| `15`                         | `%k`                                                                                      | Heure (sur 24 heures)             | 0..23               |
| `02`                         | `%m`                                                                                      | Mois                              | 01..12              |
| `04`                         | `%M`                                                                                      | Minute                            | 01..59              |
| `05`                         | `%S`                                                                                      | Seconde                           | 01..61              |
| `2006`                       | `%Y`                                                                                      | Année                             | 1970, 1984, 2019... |
| `MST`                        | `%Z`                                                                                      | Fuseau horaire                    | CEST, EDT, JST...   |

Ainsi, pour afficher transformer un champ en une date au format `heure:minute:seconde`, il faudra utiliser `formattedDate \"15:04:05\"` (même si le champ dans l'alarme ou l'événement ne correspondent pas à cette heure).

La [documentation officielle de Go](https://golang.org/pkg/time/#pkg-constants) fournit par ailleurs les valeurs à utiliser pour des formats de dates standards. Pour obtenir une date suivant le RFC3339, il faudra utiliser `formattedDate \"2006-01-02T15:04:05Z07:00\"`. De même, `formattedDate \"02 Jan 06 15:04 MST\"` sera appelé pour générer une date au format RFC822.

## Exemples

Cette section présente différents exemples de templates pour les liens et pour les payloads, accompagnés d'explications.

### Templates pour URL

#### Sans variables

```json
{
    "url" : "http://127.0.0.1:8069/post"
}
```

Pas de variables, l'adresse du service externe sera toujours la même : `http://127.0.0.1:8069/post`

#### Avec variables

```json
{
    "url" : "http://127.0.0.1:8069/edit/{{.Alarm.Value.Ticket.Value}}"
}
```

Ici l'adresse sera générée en fonction de la variable, ici le numéro de ticket.

```json
{
    "url" : "http://127.0.0.1:8069/edit/fecf03f3599769ea0"
}
```

### Templates pour payload

#### Sans variables

```json
{
    "payload" : "{\"caller_id\":1337,\"category\":\"test\",\"subcategory\":\"webhook\",\"company\":\"capensis\",\"contact_type\":\"canopsis\"}"
}
```

Ce premier exemple va envoyer toujours le même payload simple, sans utilisation de variables. Le service externe recevra :

```json
{
    "caller_id":1337,
    "category":"test",
    "subcategory":"webhook",
    "company":"capensis",
    "contact_type":"canopsis"
}
```

À noter l'absence de `\"` autour de l'id 1337, comme on souhaite l'envoyer comme un nombre et non comme une chaîne de caractères.

#### Avec variables

```json
{
    "payload" : "{{ $c := .Alarm.Value.Component }} {{ $r := .Alarm.Value.Resource }} {\"component\":\"{{$c}}\",\"resource\":\"{{$r}}\"}"
}
```

Ici, le payload sera différent en fonction du composant et de la ressource. Le payload pourra ressembler à ça

```json
{
    "component":"127.0.0.1",
    "resource":"HDD"
}
```

#### Avec variables et fonctions

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

#### Avec variables et la fonction `json`

Comme on a vu précédemment, la fonction `json` renvoie un contenu déjà compatible au format JSON. Ainsi, il est inutile besoin d'ajouter des `\"` autour des variables. Les constantes ont toujours besoin des guillemets, comme le montre l'exemple suivant.

```json
{
    "payload" : "{\"message\": {{.Alarm.Value.State.Message | json }} , \"value\": {{ .Alarm.Value.State.Value | json }}, \"type\": \"JSON\" }"
}
```

Ici, le type vaut tout le temps `"JSON"` et la présence des guillemets est obligatoire sous peine de créer un payload invalide.

```json
{
    "message": "127.0.0.1",
    "value": 2,
    "type": "JSON"
}
```

#### Formatage de la date

Cette section illustre l'utilisation de la fonction `formattedDate` pour le format des dates. Ici, l'évènement a pour timestamp `2009-11-10 23:00:00 UTC`.

La fonction utilise la syntaxe Go pour le formatage des dates. Quand la chaîne n'arrive pas à être analysée par le langage, elle est renvoyée telle quelle.

Voici un exemple avec la syntaxe Go qui va générer le résultat attendu.

```json
{
    "payload" : "{\"moment\": {{ .Event.Timestamp | formattedDate \"2006-02-01 15:04:05\" | json }} }"
}
```

```json
{
    "moment": "2009-10-11 23:00:00"
}
```

#### Fonctions en série

Enfin, on peut enchaîner plusieurs fonctions afin de transformer des variables. Dans le cas suivant, on va transformer la variable `.Event.Output` qui vaut `c0ffee - beef- facade      -a5a5a5`.

```json
{
    "payload" : "{\"message\": {{ .Event.Output | split \"-\" 2 | trim | json }} }"
}
```

On découpe l'output avec `split` qui nous retourne ` facade      `, puis trim enlève les blancs en début et fin de chaîne pour obtenir `facade` et enfin, `json` va rendre `facade` compatible pour un document JSON.

```json
{
    "message": "facade"
}
```

Voici un deuxième exemple qui combine `formattedDate` puis `json_unquote` pour générer un message concernant la date de l'événement.

```json
{
    "payload" : "{\"note\": \"The event happened on a {{ .Event.Timestamp | formattedDate \"Monday\" | json_unquote }}.\" }"
}
```

```json
{
    "note": "The event happened on a Tuesday."
}
```
