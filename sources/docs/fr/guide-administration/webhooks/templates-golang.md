# Templates (Go)

Dans les [webhooks](index.md), les champs `payload` et `url` sont personnalisables grâce aux templates Go. Les templates permettent de générer du texte en fonction de l'état de l'alarme, de l'évènement ou de l'entité.

Les templates des champs `payload` et `url` peuvent se décomposer en deux parties : la déclaration de variables et le corps du texte lui-même.

La déclaration de variables doit être positionnée avant le corps du message. Les variables se distinguent du corps du message par le fait qu'elles sont entourés de doubles accolades.

Pour plus d'informations, vous pouvez consulter la [documentaion officielle de Go sur les templates](https://golang.org/pkg/text/template).

## Déclaration de variables

!!! note
    La déclaration n'est pas obligatoire mais elle est recommandée si beaucoup de variables seront utilisées pour générer du texte.

Les variables stockent des informations sur les alarmes, les événements et les entités. `{{ .Alarm }}` permet d'accéder aux propriétés d'une alarme, de même que `{{ .Event }}` pour un évènement et `{{ .Entity }}` pour une entité. Ces trois éléments contiennent plusieurs propriétés qu'on peut utiliser pour créer des chaînes dynamiques. Par exemple, `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}` va créer la chaîne de caractères `http://mon-api.xyz/nom-du-component/nom-de-la-ressource`.

Voici une liste des principales données et la manière de la récupérer.

| Nom du champ        | Valeur                              |
|:--------------------|:------------------------------------|
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

!!! attention
    Les champs enrichis depuis un événement ou via l'event filter se retrouvent au niveau de l'entité et sont sensibles à la casse. Par exemple un champ enrichi intitulé `switch` dans l'entité sera traduit en `{{ .Entity.Infos.switch.Value }}`.

## Génération de texte

Une fois qu'on possède les variables nécessaires, la seconde étape est la génération du texte.

Pour mieux comprendre comment fonctionne la génération et comment utiliser les variables, plusieurs exemples seront présentés.

### Génération simple

Dans un premier cas, on peut utiliser directement les variables pour générer le texte. En reprenant le template `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}`, on peut voir :
- la première variable `$comp` définie comme `.Alarm.Value.Component`;
- la seconde variable `$res` qui a pour valeur `.Alarm.Value.Resource`;
- enfin le texte lui-même `http://mon-api.xyz/{{$comp}}/{{$res}}` qui va donner `http://mon-api.xyz/nom-du-component/nom-de-la-ressource` après transformation.

!!! note
    Dans cet exemple simple, on aurait pu se passer de variables et utiliser directement `http://mon-api.xyz/{{ .Alarm.Value.Component }}/{{ .Alarm.Value.Resource }}`.

### Génération selon les variables

Dans une utilisation plus avancée, on peut  générer du texte en fonction de l'état de la variable. Dans le cas suivant `"{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}"`, on trouve :
- la déclaration de variable `$val`, `{{ $val := .Alarm.Value.Status.Value }}`;
- le texte `http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}` qui contient une condition sur la variable `$val`.

La condition `{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}` va regarder si `$val` vaut soit 0, soit 2, soit 4. Si c'est le cas, on obtiendra `http://127.0.0.1:5000/even`, `http://127.0.0.1:5000/odd` sinon.

Ici, nous avons utilisé le `or` et le `eq`, mais il est possible de tester les conditions avec `and`, `not`, `ne` (not equal), `lt` (less than), `le` (less than or equal), `gt` (greater than) ou `ge` (greater than or equal).

### Transformation des variables

En plus des fontions de base pour tester la valeur des variables, il existe plusieurs fonctions pour transformer le contenu de la variable.

#### `urlquery`

`urlquery` va tarnsformer le contenu de la variable en une chaîne de caractères compatible avec le format des URL. Cette fonction a son intérêt si l'adresse du service externe dépend de l'état de l'alarme ou du ticket et que le contenu contient des caractères spéciaux. Un exemple d'adresse serait `http://une-api.org/edit/{{ .Alarm.Value.Ticket.Value | urlquery }}` pour modifier un ticket déjà existant.

#### `json` et `json_unquote`

Une fonctionnalité pas présente de base dans les templates Go mais implémentée dans Canopsis est la fonction `json`, ainsi que sa variante `json_unquote`. Ces fonctions vont transformer un champ en un élément directement insérable dans un fichier JSON et en conservant le contenu exact. Les caractères spéciaux dans les chaînes de caractères seront envoyés en tant que tel, sans traitement pas les engines Go.

`json_unquote` va retirer les guillemets si le résultat est une chaîne de caractères, sinon il retourne la même chose que `json`.

Ainsi, `{{ .Alarm.Value.ACK.Message | json }}` va renvoyer la chaîne `"ACK by someone"`, directement avec les guillemets et donc insérable dans un JSON.

`json_unquote` va lui générer la chaîne `ACK by someone` qui n'est pas utilisable directement dans un JSON mais qui peut servir pour créer des textes plus complexes. `\"Voici un message : '{{ .Alarm.Value.ACK.Message | json_unquote }}'\"` va donner par exemple `"Voici un message : 'ACK by someone'"`.

!!! note
    `json` et `json_unquote` sont disponibles uniquement dans les templates des [webhooks](index.md), pas ceux de l'event-filter.

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

Ici l'adresse sera généré en fonction de la variable, ici le numéro de ticket.

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
    "payload" : "{{ $c := .Alarm.Value.Component }} {{ $r := .Alarm.Value.Resource }} {\"component\":\"{{$c}}\",\"ressource\":\"{{$r}}\"}"
}
```

Ici, le payload sera différent en fonction du composant et de la ressource. Le payload pourra ressembler à ça

```json
{
    "component":"127.0.0.1",
    "ressource":"HDD"
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
    "ressource":"HDD",
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