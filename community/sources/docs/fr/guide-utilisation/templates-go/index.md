# Templates Go dans Canopsis

## Introduction

Les Templates Go sont un puissant système de génération de texte utilisé dans de nombreuses fonctionnalités de Canopsis. Ils vous permettent de créer du contenu dynamique en exploitant les données des événements, alarmes et entités de votre système de supervision.

**À quoi servent les Templates Go dans Canopsis ?**

*  Créer des messages personnalisés pour les notifications
*  Générer des payloads pour les webhooks
*  Construire des URLs dynamiques
*  Formatter les messages d'alarmes
*  Personnaliser l'affichage des informations
*  Automatiser la création de tickets

Cette documentation vous guidera pas à pas dans l'utilisation des Templates Go, des concepts fondamentaux aux cas d'usage avancés.

Des exemples avancés sont disponible sur [cette page dédiée](./templates-go-casdusage/)

## Concepts fondamentaux

Les Templates Go dans Canopsis s'appuient sur le [package officiel text/template de Go](https://pkg.go.dev/text/template). Cependant, Canopsis enrichit ce système avec ses propres fonctions adaptées au contexte de la supervision.

### Structure d'un template

Un template Go se compose de deux éléments principaux :

*  Du **texte statique** qui sera reproduit tel quel
*  Des **actions** entre doubles accolades `{{ }}` qui seront évaluées et remplacées par leur résultat

Par exemple, dans le template suivant :
```
L'alarme sur le composant {{ .Alarm.Value.Component }} est en criticité {{ .Alarm.Value.State.Value }}
```

Le texte statique est "L'alarme sur le composant " et " est en criticité ", tandis que `{{ .Alarm.Value.Component }}` et `{{ .Alarm.Value.State.Value }}` sont des actions qui seront remplacées par leurs valeurs respectives.

## Variables disponibles

Canopsis met à votre disposition plusieurs contextes de variables selon la source des données :

### Variables d'alarme

Les variables d'alarme sont accessibles via le préfixe `.Alarm`. Elles contiennent toutes les informations relatives à une alarme en cours.

```go
// Informations de base
{{ .Alarm.Value.Component }}     // Nom du composant
{{ .Alarm.Value.Resource }}      // Nom de la ressource
{{ .Alarm.Value.DisplayName }}   // Nom d'affichage
{{ .Alarm.Value.Connector }}     // Type de connecteur (nagios, snmp, etc.)
{{ .Alarm.Value.ConnectorName }} // Nom du connecteur

// Messages
{{ .Alarm.Value.Output }}        // Message actuel de l'alarme
{{ .Alarm.Value.InitialOutput }} // Message initial lors de la création de l'alarme
{{ .Alarm.Value.LongOutput }}    // Message long de l'alarme
{{ .Alarm.Value.InitialLongOutput }} // Message long initial

// État et Statut
{{ .Alarm.Value.State.Value }}   // Criticité (0=OK, 1=MINOR, 2=MAJOR, 3=CRITICAL)
{{ .Alarm.Value.State.Message }} // Message associé au dernier changement de criticité
{{ .Alarm.Value.Status.Value }}  // Statut de l'alarme

// Dates importantes
{{ .Alarm.Value.CreationDate }}   // Date de création
{{ .Alarm.Value.ActivationDate }} // Date d'activation
{{ .Alarm.Value.LastUpdateDate }} // Date du dernier changement de sévérité
{{ .Alarm.Value.LastEventDate }}  // Date du dernier événement reçu
{{ .Alarm.Value.Resolved }}       // Date de résolution

// Durées
{{ .Alarm.Value.InactiveDuration }}           // Durée d'inactivité
{{ .Alarm.Value.PbehaviorInactiveDuration }}  // Inactivité liée aux comportements périodiques
{{ .Alarm.Value.SnoozeDuration }}             // Durée de mise en veille

// Tickets et acquittements
{{ .Alarm.Value.Ticket.Author }}  // Auteur du ticket
{{ .Alarm.Value.Ticket.Ticket }}  // ID du ticket
{{ .Alarm.Value.Ticket.Message }} // Message du ticket
{{ .Alarm.Value.ACK.Author }}     // Auteur de l'acquittement
{{ .Alarm.Value.ACK.Message }}    // Message d'acquittement

// Commentaires
{{ .Alarm.Value.LastComment.Author }}  // Auteur du dernier commentaire
{{ .Alarm.Value.LastComment.Message }} // Message du dernier commentaire

// Compteurs
{{ .Alarm.Value.TotalStateChanges }}            // Nombre total de changements d'état
{{ .Alarm.Value.StateChangesSinceStatusUpdate }} // Changements depuis dernière mise à jour
{{ .Alarm.Value.EventsCount }}                   // Nombre d'événements reçus

// Infos dynamiques
{{ (index (index .Alarm.Value.Infos "rule_id") "info_name") }} // Accès aux infos dynamiques
```

### Variables d'entité

Les entités représentent les objets surveillés dans Canopsis (serveurs, applications, services...).

```go
{{ .Entity.ID }}     // Identifiant unique de l'entité
{{ .Entity.Name }}   // Nom de l'entité

// Informations associées à l'entité
{{ (index .Entity.Infos "nom_info").Value }}          // Valeur d'une info spécifique
{{ (index .Entity.ComponentInfos "nom_info").Value }} // Info du composant
```

### Variables d'événement

Les événements sont les données brutes envoyées à Canopsis par les systèmes de monitoring.

```go
{{ .Event.Component }}     // Composant concerné
{{ .Event.Resource }}      // Ressource concernée
{{ .Event.Connector }}     // Type de connecteur
{{ .Event.ConnectorName }} // Nom du connecteur
{{ .Event.EventType }}     // Type d'événement
{{ .Event.Output }}        // Message de l'événement
{{ .Event.Author }}        // Auteur de l'événement

// Informations supplémentaires
{{ index .Event.ExtraInfos "nom_info" }} // Accès à une info supplémentaire
```

### Variables d'environnement

Vos propres variables définies dans le fichier [Canopsis.toml](../../../guide-administration/administration-avancee/modification-canopsis-toml/#section-canopsistemplatevars).

```go
{{ .Env.nom_variable }} // Accès à une variable d'environnement
```

### Variables additionnelles

Disponibles dans les scénarios et les règles de déclaration de tickets.

```go
{{ .AdditionalData.RuleName }}       // Nom de la règle
{{ .AdditionalData.Trigger }} // Type de déclencheur (ack, stateinc, etc.)
{{ .AdditionalData.Author }}         // Auteur de l'action
{{ .AdditionalData.Initiator }}      // Origine de l'action (user/system)
{{ .AdditionalData.Output }}         // Message de l'événement
```

### Variables pour les méta-alarmes

Pour parcourir les alarmes conséquences d'une méta alarme.

```go
{{ range $children := .Children }}
    ID: {{ $children.ID }} - Message: {{ $children.Value.State.Message }}
{{ end }}
```

## Utilisation des variables

### Accès simple

Pour accéder à une variable, utilisez simplement son chemin complet :

```go
{{ .Alarm.Value.Component }}
```

### Déclaration de variables temporaires

Vous pouvez déclarer des variables pour simplifier votre template :

```go
{{ $component := .Alarm.Value.Component }}
{{ $resource := .Alarm.Value.Resource }}

Problème détecté sur {{ $component }}/{{ $resource }}
```

### Accès aux informations structurées

Pour les données complexes comme les maps ou les listes :

```go
// Accès à un élément d'une map par son nom
{{ index .Event.ExtraInfos "location" }}

// Accès sécurisé avec vérification préalable
{{ if map_has_key .Event.ExtraInfos "location" }}
    {{ index .Event.ExtraInfos "location" }}
{{ else }}
    Emplacement inconnu
{{ end }}
```

## Opérations conditionnelles

Les templates Go permettent d'effectuer des opérations conditionnelles pour adapter le contenu généré.

### Conditions simples

```go
{{ if eq .Alarm.Value.State.Value 3 }}
    CRITIQUE: Intervention immédiate requise!
{{ else if eq .Alarm.Value.State.Value 2 }}
    MAJEUR: Intervention nécessaire
{{ else if eq .Alarm.Value.State.Value 1 }}
    MINEUR: À surveiller
{{ else }}
    OK: Tout fonctionne normalement
{{ end }}
```

### Opérateurs de comparaison

- `eq` : égalité (equal)
- `ne` : différence (not equal)
- `lt` : inférieur à (less than)
- `le` : inférieur ou égal à (less than or equal)
- `gt` : supérieur à (greater than)
- `ge` : supérieur ou égal à (greater than or equal)

### Opérateurs logiques

- `and` : ET logique
- `or` : OU logique
- `not` : négation

### Exemple combiné

```go
{{ if and (ge .Alarm.Value.State.Value 2) (lt .Alarm.Value.EventsCount 5) }}
    Alarme importante avec peu d'événements: possible faux positif
{{ end }}
```

## Boucles et itérations

### Parcourir une liste

```go
{{ range $index, $tag := .Alarm.Tags }}
    Tag #{{ $index }}: {{ $tag }}
{{ end }}
```

### Parcourir une map

```go
{{ range $key, $value := .Event.ExtraInfos }}
    {{ $key }}: {{ $value }}
{{ end }}
```

## Fonctions de transformation

Canopsis fournit de nombreuses fonctions pour transformer les données dans vos templates. Ces fonctions s'utilisent avec le pipe (`|`).

### Manipulation de texte

```go
// Conversion en majuscules/minuscules
{{ .Alarm.Value.Output | uppercase }}  // Convertit en majuscules
{{ .Alarm.Value.Output | lowercase }}  // Convertit en minuscules

// Suppression des espaces
{{ .Alarm.Value.Output | trim }}  // Supprime les espaces en début et fin

// Extraction de parties d'une chaîne
{{ .Alarm.Value.Output | split "/" 0 }}  // Split par "/" et prend le 1er élément
{{ substrLeft .Alarm.Value.Output 10 }}  // Prend les 10 premiers caractères
{{ substrRight .Alarm.Value.Output 10 }} // Prend les 10 derniers caractères
{{ substr .Alarm.Value.Output 5 10 }}    // Prend 10 caractères à partir du 5ème

// Remplacement de texte
{{ .Alarm.Value.Output | replace "erreur" "ERREUR" }} // Remplace un texte
```

### Traitement de dates

```go
// Formatage de date dans le fuseau horaire local
{{ .Alarm.Value.CreationDate | localtime "02/01/2006 15:04:05" }}

// Formatage dans un fuseau horaire spécifique
{{ .Alarm.Value.CreationDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }}
```

### Manipulation de données structurées

```go
// Encodage JSON
{{ .Entity.Infos | json }}

// Encodage JSON sans guillemets
{{ .Entity.Infos | json_unquote }}

// Vérification de présence
{{ tag_has_key .Alarm.Tags "Priorité" }}  // Vérifie si un tag existe
{{ get_tag .Alarm.Tags "Priorité" }}      // Récupère la valeur d'un tag
```

### Opérations mathématiques

```go
{{ add 2 3 }}  // Addition: 5
{{ sub 5 2 }}  // Soustraction: 3
{{ mult 2 3 }} // Multiplication: 6
{{ div 6 2 }}  // Division: 3 (division entière)
```

## Exemples de cas d'usage concrets

Voici quelques exemples d'usages simples et concrets utilisant des templates GO.  
Si vous souhaitez rédiger des templates plus complexes, n'hésitez pas à consulter [cette page dédiée](./templates-go-casdusage/).

### Cas #1: Notification personnalisée selon la criticité

```go
{{ $comp := .Alarm.Value.Component }}
{{ $res := .Alarm.Value.Resource }}
{{ $crit := .Alarm.Value.State.Value }}
{{ $msg := .Alarm.Value.Output }}

{{ if eq $crit 3 }}
[CRITIQUE] Incident sur {{ $comp }}/{{ $res }}: {{ $msg }}
Veuillez contacter l'équipe d'astreinte immédiatement au 01-23-45-67-89
{{ else if eq $crit 2 }}
[MAJEUR] Problème sur {{ $comp }}/{{ $res }}: {{ $msg }}
Merci d'intervenir dans les prochaines heures
{{ else if eq $crit 1 }}
[MINEUR] Anomalie sur {{ $comp }}/{{ $res }}: {{ $msg }}
À traiter dans les prochains jours
{{ else }}
[INFORMATION] {{ $comp }}/{{ $res }}: {{ $msg }}
{{ end }}
```

### Cas #2: Création d'un ticket avec informations enrichies

```go
{
  "title": "[{{ if eq .Alarm.Value.State.Value 3 }}CRITIQUE{{ else if eq .Alarm.Value.State.Value 2 }}MAJEUR{{ else }}MINEUR{{ end }}] Problème sur {{ .Alarm.Value.Component }}",
  "description": "## Détails de l'alarme\n\n- **Composant**: {{ .Alarm.Value.Component }}\n- **Ressource**: {{ .Alarm.Value.Resource }}\n- **Criticité**: {{ .Alarm.Value.State.Value }}\n- **Message**: {{ .Alarm.Value.Output }}\n\n## Informations complémentaires\n\n{{ if map_has_key .Event.ExtraInfos \"location\" }}**Localisation**: {{ index .Event.ExtraInfos \"location\" }}\n{{ end }}{{ if map_has_key .Event.ExtraInfos \"responsible\" }}**Responsable**: {{ index .Event.ExtraInfos \"responsible\" }}\n{{ end }}\n\n## Lien vers Canopsis\n\nhttp://canopsis-url/alarm/{{ .Alarm.ID }}",
  "priority": "{{ if ge .Alarm.Value.State.Value 3 }}P1{{ else if eq .Alarm.Value.State.Value 2 }}P2{{ else }}P3{{ end }}",
  "assignee": "{{ if map_has_key .Event.ExtraInfos \"responsible\" }}{{ index .Event.ExtraInfos \"responsible\" }}{{ else }}support_team{{ end }}"
}
```

### Cas #3: URL dynamique vers un système externe

```go
{{ $env := .Env.environment }}
{{ $component := .Alarm.Value.Component | lowercase | replace " " "-" }}
{{ $resource := .Alarm.Value.Resource | lowercase | replace " " "-" }}

https://wiki.example.com/{{ $env }}/components/{{ $component }}/{{ $resource }}/troubleshooting
```

### Cas #4: Résumé d'une méta-alarme

```go
{
  "summary": "Incident global: {{ .Alarm.Value.Component }}",
  "details": "## Vue d'ensemble\n\nUn incident global affecte {{ .Alarm.Value.Component }}.\n\n## Systèmes impactés\n\n{{ range $index, $children := .Children }}{{ if $index }}\n{{ end }}* **{{ $children.Value.Component }}/{{ $children.Value.Resource }}**: {{ $children.Value.State.Message }}{{ end }}\n\n## Recommandations\n\nVérifier d'abord les systèmes en criticité {{ .Alarm.Value.State.Value }}."
}
```

### Cas #5: Traitement de données JSON

```go
{{ $data := index .Event.ExtraInfos "metrics" | json_unquote }}
{{ $cpu := index $data "cpu_usage" }}
{{ $mem := index $data "memory_usage" }}

Utilisation des ressources:
- CPU: {{ $cpu }}% {{ if gt $cpu 90 }}(CRITIQUE){{ else if gt $cpu 75 }}(ÉLEVÉ){{ else }}(NORMAL){{ end }}
- Mémoire: {{ $mem }}% {{ if gt $mem 90 }}(CRITIQUE){{ else if gt $mem 75 }}(ÉLEVÉ){{ else }}(NORMAL){{ end }}
```

## Conseils pour des templates efficaces

1. **Testez vos templates** avant de les mettre en production pour éviter les surprises
2. **Vérifiez l'existence des variables** avant de les utiliser avec `map_has_key`
3. **Utilisez des variables temporaires** pour simplifier vos templates complexes
4. **Pensez à la lisibilité** de votre template et du résultat généré
5. **Commentez vos templates complexes** pour faciliter leur maintenance

## Exemples de formats de date

Les formats de date suivent la convention Go et utilisent comme référence `Mon Jan 2 15:04:05 MST 2006` (01/02 03:04:05PM '06 -0700).

| Format | Résultat |
|--------|----------|
| `{{ .Alarm.Value.CreationDate \| localtime "02/01/2006 15:04:05" }}` | 28/10/2021 09:05:00 |
| `{{ .Alarm.Value.CreationDate \| localtime "Monday, January 2, 2006" }}` | Thursday, October 28, 2021 |
| `{{ .Alarm.Value.CreationDate \| localtime "15:04:05" }}` | 09:05:00 |
| `{{ .Alarm.Value.CreationDate \| localtime "2006-01-02T15:04:05Z07:00" }}` | 2021-10-28T09:05:00+02:00 |

## Résolution des problèmes courants

### Le template ne s'évalue pas correctement

Vérifiez que :

* Les accolades sont bien fermées et correctement imbriquées
* Les variables utilisées existent bien dans le contexte actuel
* Les opérateurs de comparaison sont correctement utilisés

### Une variable est inaccessible

Si vous rencontrez des erreurs d'accès aux variables :

* Assurez-vous que la variable existe dans le contexte (alarme, événement, entité)
* Vérifiez l'orthographe et la casse (les variables sont sensibles à la casse)
* Pour les structures imbriquées, utilisez `index` pour un accès sécurisé

### Problèmes avec les fonctions personnalisées

Si une fonction ne fonctionne pas comme prévu :

* Vérifiez l'ordre des arguments
* Utilisez correctement l'opérateur pipe (`|`)
* Consultez les exemples de la documentation

## Conclusion

Les Templates Go dans Canopsis sont un outil puissant pour personnaliser vos messages et automatiser vos processus de supervision. En maîtrisant leur utilisation, vous pouvez créer des notifications précises et contextuelles, générer des payloads adaptés pour vos intégrations, et faciliter le travail de vos équipes de supervision.

N'hésitez pas à commencer avec des templates simples et à les enrichir progressivement en fonction de vos besoins spécifiques.
