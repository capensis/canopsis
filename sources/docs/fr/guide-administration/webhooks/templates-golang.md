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
    Le champ dans les informations d'entité sont sensibles à la casse. Par exemple un champ `switch` dans l'entité sera traduit en `{{ .Entity.Infos.switch.Value }}`.

## Génération de texte

Une fois qu'on possède les variables nécessaires, la seconde étape est la génération du texte.

Pour mieux comprendre comment fonctionne la génération et comment utiliser les variables, plusieurs exemples seront présentés.

### Génération simple, sans transformation

Dans un cas simple, on peut utilise directement les variables pour générer le texte.

### Exemples

Dans l'exemple suivant, `{{ $comp := .Alarm.Value.Component }}{{ $res := .Alarm.Value.Resource }}http://mon-api.xyz/{{$comp}}/{{$res}}`, on déclare d'abord deux variables, `$comp` et `$res`. Ensuite, on utilise ces deux variables pour générer l'adress URL que va appeler le moteur axe, par exemple `http://mon-api.xyz/nom-du-component/nom-de-la-ressource`.

On peut également générer du texte en fonction de l'état de la variable. Dans le cas suivant `"{{ $val := .Alarm.Value.Status.Value }}http://127.0.0.1:5000/{{if ((eq $val 0) or (eq $val 2) or (eq $val 4))}}even{{else}}odd{{end}}"`, on obtiendra `"http://127.0.0.1:5000/even"` si le statut de l'alarm vaut 0, 2 ou 4, `"http://127.0.0.1:5000/odd"` sinon.

De même façon que le `or` et le `eq`, il est possible de tester les conditions avec `and`, `not`, `ne` (not equal), `lt` (less than), `le` (less than or equal), `gt` (greater than) ou `ge` (greater than or equal).

La fonction `js`, qui renvoie une chaîne de caractères échappée, peut être également mentionnée. Si, par exemple, la valeur dans `{{ .Event.Output }}` contient des caractères spéciaux comme des guillemets ou des backslashs, `{{ .Event.Output | js }}` permet d'échapper ces caractères.
