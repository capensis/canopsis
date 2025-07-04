# Exemples et cas d'usage concrets pour les Templates Go dans Canopsis

Cette section présente des exemples pratiques de templates Go pour répondre à des besoins courants dans Canopsis. Chaque exemple est accompagné d'explications et de conseils pour l'adapter à vos besoins.

## 1. Notification Microsoft Teams

Microsoft Teams utilise le format de carte adaptative pour les notifications. Le template suivant crée une notification Teams avec différentes couleurs selon la criticité de l'alarme.

```go
{
  "type": "message",
  "attachments": [
    {
      "contentType": "application/vnd.microsoft.card.adaptive",
      "content": {
        "type": "AdaptiveCard",
        "body": [
          {
            "type": "TextBlock",
            "size": "Large",
            "weight": "Bolder",
            "color": "{{ if eq .Alarm.Value.State.Value 3 }}Attention{{ else if eq .Alarm.Value.State.Value 2 }}Warning{{ else if eq .Alarm.Value.State.Value 1 }}Accent{{ else }}Good{{ end }}",
            "text": "{{ if eq .Alarm.Value.State.Value 3 }}⚠️ CRITIQUE{{ else if eq .Alarm.Value.State.Value 2 }}⚠️ MAJEUR{{ else if eq .Alarm.Value.State.Value 1 }}ℹ️ MINEUR{{ else }}✅ OK{{ end }} - {{ .Alarm.Value.Component }}"
          },
          {
            "type": "FactSet",
            "facts": [
              {
                "title": "Composant",
                "value": "{{ .Alarm.Value.Component }}"
              },
              {
                "title": "Ressource",
                "value": "{{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}N/A{{ end }}"
              },
              {
                "title": "Criticité",
                "value": "{{ if eq .Alarm.Value.State.Value 3 }}CRITIQUE (3){{ else if eq .Alarm.Value.State.Value 2 }}MAJEUR (2){{ else if eq .Alarm.Value.State.Value 1 }}MINEUR (1){{ else }}OK (0){{ end }}"
              },
              {
                "title": "Date",
                "value": "{{ .Alarm.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }}"
              },
              {
                "title": "Connecteur",
                "value": "{{ .Alarm.Value.ConnectorName }}"
              }
            ]
          },
          {
            "type": "TextBlock",
            "text": "**Message:**",
            "wrap": true
          },
          {
            "type": "TextBlock",
            "text": "{{ .Alarm.Value.Output | replace "\n" "\n\n" }}",
            "wrap": true
          },
          {{ if .Alarm.Value.LongOutput }}
          {
            "type": "TextBlock",
            "text": "**Détails:**",
            "wrap": true
          },
          {
            "type": "TextBlock",
            "text": "{{ .Alarm.Value.LongOutput | replace "\n" "\n\n" }}",
            "wrap": true
          },
          {{ end }}
          {
            "type": "TextBlock",
            "text": "**Événements:** {{ .Alarm.Value.EventsCount }} | **Changements d'état:** {{ .Alarm.Value.TotalStateChanges }}",
            "wrap": true
          }
        ],
        "actions": [
          {
            "type": "Action.OpenUrl",
            "title": "Voir dans Canopsis",
            "url": "{{ .Env.CanopsisURL }}/alarms/{{ .Alarm.ID }}"
          }
        ],
        "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
        "version": "1.2"
      }
    }
  ]
}
```

**Cas d'usage:** Ce template est idéal pour notifier les équipes d'exploitation lors de problèmes critiques nécessitant une attention rapide. L'utilisation des couleurs et des icônes permet une identification visuelle immédiate de la gravité.

**À noter:**  

* La variable `.Env.CanopsisURL` doit être définie dans le fichier [Canopsis.toml](../../../../guide-administration/administration-avancee/modification-canopsis-toml/#section-canopsistemplatevars)
* La fonction `replace` avec `\n` permet d'améliorer le formatage du texte dans Teams. Notez également que l'interface de Canopsis se chargera de protéger les backslashs.

## 2. Notification Mattermost

Mattermost utilise le markdown pour formatter les messages. Voici un template pour une notification Mattermost riche et informative:

```go
{
  "text": "{{ if eq .Alarm.Value.State.Value 3 }}:red_circle: **CRITIQUE**{{ else if eq .Alarm.Value.State.Value 2 }}:large_orange_diamond: **MAJEUR**{{ else if eq .Alarm.Value.State.Value 1 }}:warning: **MINEUR**{{ else }}:white_check_mark: **OK**{{ end }}",
  "username": "Canopsis",
  "icon_url": "{{ .Env.CanopsisLogoURL }}",
  "attachments": [
    {
      "color": "{{ if eq .Alarm.Value.State.Value 3 }}#FF0000{{ else if eq .Alarm.Value.State.Value 2 }}#FFA500{{ else if eq .Alarm.Value.State.Value 1 }}#FFFF00{{ else }}#00FF00{{ end }}",
      "title": "{{ .Alarm.Value.Component }}{{ if .Alarm.Value.Resource }}/{{ .Alarm.Value.Resource }}{{ end }}",
      "title_link": "{{ .Env.CanopsisURL }}/alarms/{{ .Alarm.ID }}",
      "fields": [
        {
          "short": true,
          "title": "Composant",
          "value": "{{ .Alarm.Value.Component }}"
        },
        {
          "short": true,
          "title": "Ressource",
          "value": "{{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}N/A{{ end }}"
        },
        {
          "short": true,
          "title": "Criticité",
          "value": "{{ if eq .Alarm.Value.State.Value 3 }}CRITIQUE (3){{ else if eq .Alarm.Value.State.Value 2 }}MAJEUR (2){{ else if eq .Alarm.Value.State.Value 1 }}MINEUR (1){{ else }}OK (0){{ end }}"
        },
        {
          "short": true,
          "title": "Statut",
          "value": "{{ if eq .Alarm.Value.Status.Value 1 }}En cours{{ else if eq .Alarm.Value.Status.Value 0 }}Résolue{{ end }}"
        },
        {
          "short": true,
          "title": "Date",
          "value": "{{ .Alarm.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" }}"
        },
        {
          "short": true,
          "title": "Source",
          "value": "{{ .Alarm.Value.Connector }}/{{ .Alarm.Value.ConnectorName }}"
        }
      ],
      "text": "**Message:** {{ .Alarm.Value.Output }}\n{{ if .Alarm.Value.LongOutput }}**Détails:** {{ .Alarm.Value.LongOutput }}{{ end }}"
    }
  ]
}
```

**Cas d'usage:** Parfait pour les équipes qui utilisent Mattermost comme principal outil de communication d'équipe. L'information est structurée et complète, tout en restant visuelle.

**À noter:**

* Les variables `.Env.CanopsisURL` et `.Env.CanopsisLogoURL` doivent être définies
* Les émoticônes peuvent être personnalisées selon votre instance Mattermost

## 3. Rapport complet au format Markdown

Ce template génère un rapport détaillé au format Markdown pour une alarme spécifique:

```go
# Rapport d'alarme - {{ .Alarm.Value.Component }}{{ if .Alarm.Value.Resource }}/{{ .Alarm.Value.Resource }}{{ end }}
{{ $state := .Alarm.Value.State.Value }}

## Informations générales

| Propriété | Valeur |
|-----------|--------|
| **ID** | `{{ .Alarm.ID }}` |
| **Composant** | {{ .Alarm.Value.Component }} |
| **Ressource** | {{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}N/A{{ end }} |
| **Nom d'affichage** | {{ if .Alarm.Value.DisplayName }}{{ .Alarm.Value.DisplayName }}{{ else }}Non défini{{ end }} |
| **Criticité** | {{ if eq $state 3 }}:red_circle: CRITIQUE (3){{ else if eq $state 2 }}:orange_circle: MAJEUR (2){{ else if eq $state 1 }}:yellow_circle: MINEUR (1){{ else }}:green_circle: OK (0){{ end }} |
| **Statut** | {{ if eq .Alarm.Value.Status.Value 0 }}:person_raising_hand: Résolue{{ else if eq .Alarm.Value.Status.Value 1 }}:wrench: En cours{{ end }} |
| **Connecteur** | {{ .Alarm.Value.Connector }}/{{ .Alarm.Value.ConnectorName }} |

## Chronologie

| Événement | Date |
|-----------|------|
| **Création** | {{ .Alarm.Value.CreationDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }} |
| **Activation** | {{ if .Alarm.Value.ActivationDate }}{{ .Alarm.Value.ActivationDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }}{{ else }}N/A{{ end }} |
| **Dernier changement d'état** | {{ .Alarm.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }} |
| **Dernier événement** | {{ .Alarm.Value.LastEventDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }} |
| **Résolution** | {{ if .Alarm.Value.Resolved }}{{ .Alarm.Value.Resolved | localtime "02/01/2006 15:04:05" "Europe/Paris" }}{{ else }}Non résolue{{ end }} |

## Messages

### Message actuel

{{ .Alarm.Value.Output }}


{{ if .Alarm.Value.LongOutput }}
### Détails supplémentaires


{{ .Alarm.Value.LongOutput }}

{{ end }}

### Message initial

{{ .Alarm.Value.InitialOutput }}


## Statistiques

| Métrique | Valeur |
|----------|--------|
| **Nombre d'événements** | {{ .Alarm.Value.EventsCount }} |
| **Changements d'état totaux** | {{ .Alarm.Value.TotalStateChanges }} |
| **Changements depuis dernière mise à jour** | {{ .Alarm.Value.StateChangesSinceStatusUpdate }} |
| **Durée d'inactivité** | {{ .Alarm.Value.InactiveDuration }} |

{{ if .Alarm.Value.ACK }}
## Informations d'acquittement

| Propriété | Valeur |
|-----------|--------|
| **Auteur** | {{ .Alarm.Value.ACK.Author }} |
| **Message** | {{ if .Alarm.Value.ACK.Message }}{{ .Alarm.Value.ACK.Message }}{{ else }}Aucun message{{ end }} |
{{ end }}

{{ if .Alarm.Value.Ticket }}
## Ticket associé

| Propriété | Valeur |
|-----------|--------|
| **Numéro** | {{ .Alarm.Value.Ticket.Ticket }} |
| **Auteur** | {{ .Alarm.Value.Ticket.Author }} |
| **Message** | {{ if .Alarm.Value.Ticket.Message }}{{ .Alarm.Value.Ticket.Message }}{{ else }}Aucun message{{ end }} |
{{ end }}

{{ if .Alarm.Value.LastComment }}
## Dernier commentaire

| Propriété | Valeur |
|-----------|--------|
| **Auteur** | {{ .Alarm.Value.LastComment.Author }} |
| **Message** | {{ .Alarm.Value.LastComment.Message }} |
{{ end }}

## Informations sur l'entité

| Propriété | Valeur |
|-----------|--------|
| **ID** | {{ .Entity.ID }} |
| **Nom** | {{ .Entity.Name }} |

{{ if .Entity.Infos }}
### Informations supplémentaires
{{ range $key, $info := .Entity.Infos }}
- **{{ $key }}**: {{ $info.Value }}
{{ end }}
{{ end }}

---
*Rapport généré automatiquement le {{ .Alarm.Value.CreationDate | localtime "02/01/2006 à 15:04:05" "Europe/Paris" }} - Canopsis {{ .Env.CanopsisVersion }}*
```

**Cas d'usage:** Ce template est utile pour générer des rapports complets sur une alarme spécifique, à partager par email ou à intégrer dans une documentation d'incident. Il contient toutes les informations pertinentes sur l'alarme et son contexte.

**À noter:**

* La variable `.Env.CanopsisVersion` doit être définie dans le fichier Canopsis.toml
* Le rapport utilise des tableaux Markdown pour une meilleure lisibilité

## 4. Rapport Markdown pour une méta-alarme

Template pour générer un rapport détaillé sur une méta-alarme et ses alarmes conséquences:

```go
# Rapport de méta-alarme - {{ .Alarm.Value.Component }}
{{ $state := .Alarm.Value.State.Value }}

## Résumé de la méta-alarme

| Propriété | Valeur |
|-----------|--------|
| **ID** | `{{ .Alarm.ID }}` |
| **Composant** | {{ .Alarm.Value.Component }} |
| **Ressource** | {{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}N/A{{ end }} |
| **Criticité globale** | {{ if eq $state 3 }}:red_circle: CRITIQUE (3){{ else if eq $state 2 }}:orange_circle: MAJEUR (2){{ else if eq $state 1 }}:yellow_circle: MINEUR (1){{ else }}:green_circle: OK (0){{ end }} |
| **Statut** | {{ if eq .Alarm.Value.Status.Value 1 }}:person_raising_hand: En cours{{ else if eq .Alarm.Value.Status.Value 0 }}:x: Résolue{{ end }} |
| **Date de dernier changement** | {{ .Alarm.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }} |
| **Nombre d'alarmes associées** | {{ len .Children }} |

## Message principal


{{ .Alarm.Value.Output }}


{{ if .Alarm.Value.LongOutput }}
### Détails supplémentaires

{{ .Alarm.Value.LongOutput }}

{{ end }}

## Alarmes conséquences

{{ range $index, $child := .Children }}
{{ $childState := $child.Value.State.Value }}
### {{ add $index 1 }}. {{ $child.Value.Component }}{{ if $child.Value.Resource }}/{{ $child.Value.Resource }}{{ end }} - {{ if eq $childState 3 }}:red_circle: CRITIQUE{{ else if eq $childState 2 }}:orange_circle: MAJEUR{{ else if eq $childState 1 }}:yellow_circle: MINEUR{{ else }}:green_circle: OK{{ end }}

| Propriété | Valeur |
|-----------|--------|
| **ID** | `{{ $child.ID }}` |
| **Composant** | {{ $child.Value.Component }} |
| **Ressource** | {{ if $child.Value.Resource }}{{ $child.Value.Resource }}{{ else }}N/A{{ end }} |
| **Criticité** | {{ if eq $childState 3 }}CRITIQUE (3){{ else if eq $childState 2 }}MAJEUR (2){{ else if eq $childState 1 }}MINEUR (1){{ else }}OK (0){{ end }} |
| **Statut** | {{ if eq $child.Value.Status.Value 1 }}en cours{{ else if eq $child.Value.Status.Value 0 }}Résolue{{ end }} |
| **Dernière mise à jour** | {{ $child.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }} |
| **Nombre d'événements** | {{ $child.Value.EventsCount }} |

#### Message

{{ $child.Value.Output }}


{{ if $child.Value.LongOutput }}
#### Détails

{{ $child.Value.LongOutput }}

{{ end }}

{{ end }}

## Analyse des impacts

| Criticité | Nombre d'alarmes |
|-----------|------------------|
{{ $crit3 := 0 }}
{{ $crit2 := 0 }}
{{ $crit1 := 0 }}
{{ $crit0 := 0 }}
{{ range $child := .Children }}
  {{ if eq $child.Value.State.Value 3 }}{{ $crit3 = add $crit3 1 }}{{ end }}
  {{ if eq $child.Value.State.Value 2 }}{{ $crit2 = add $crit2 1 }}{{ end }}
  {{ if eq $child.Value.State.Value 1 }}{{ $crit1 = add $crit1 1 }}{{ end }}
  {{ if eq $child.Value.State.Value 0 }}{{ $crit0 = add $crit0 1 }}{{ end }}
{{ end }}
| **CRITIQUE (3)** | {{ $crit3 }} |
| **MAJEUR (2)** | {{ $crit2 }} |
| **MINEUR (1)** | {{ $crit1 }} |
| **OK (0)** | {{ $crit0 }} |

## Recommandations

{{ if gt $crit3 0 }}
- Traiter en priorité les {{ $crit3 }} alarmes critiques
{{ end }}
{{ if gt $crit2 0 }}
- Traiter ensuite les {{ $crit2 }} alarmes majeures
{{ end }}
{{ if gt $crit1 0 }}
- Planifier le traitement des {{ $crit1 }} alarmes mineures
{{ end }}
{{ if gt $crit0 0 }}
- Vérifier l'état de fonctionnement normal des {{ $crit0 }} composants OK
{{ end }}

{{ if .Alarm.Value.ACK }}
## Informations d'acquittement

| Propriété | Valeur |
|-----------|--------|
| **Auteur** | {{ .Alarm.Value.ACK.Author }} |
| **Message** | {{ if .Alarm.Value.ACK.Message }}{{ .Alarm.Value.ACK.Message }}{{ else }}Aucun message{{ end }} |
{{ end }}

{{ if .Alarm.Value.Ticket }}
## Ticket associé

| Propriété | Valeur |
|-----------|--------|
| **Numéro** | {{ .Alarm.Value.Ticket.Ticket }} |
| **Auteur** | {{ .Alarm.Value.Ticket.Author }} |
| **Message** | {{ if .Alarm.Value.Ticket.Message }}{{ .Alarm.Value.Ticket.Message }}{{ else }}Aucun message{{ end }} |
{{ end }}

---
*Rapport de méta-alarme généré automatiquement le {{ .Alarm.Value.CreationDate | localtime "02/01/2006 à 15:04:05" "Europe/Paris" }} - Canopsis {{ .Env.CanopsisVersion }}*
```

**Cas d'usage:** Idéal pour documenter et analyser les incidents majeurs qui impactent plusieurs composants. Ce template permet d'avoir une vue d'ensemble de la situation tout en conservant les détails de chaque alarme conséquence.

**À noter:**

* L'utilisation de la fonction `len` pour compter le nombre d'alarmes conséquences
* L'utilisation des variables temporaires (`$crit3`, `$crit2`, etc.) pour compter les alarmes par criticité
* La fonction `add` pour l'incrémentation des compteurs et la numérotation

## 5. Notification au format HTML/Email

Ce template génère un email HTML complet prêt à être envoyé:

```go
From: Canopsis <{{ .Env.CanopsisEmailSender }}>
To: {{ .Env.AlertRecipients }}
Subject: {{ if eq .Alarm.Value.State.Value 3 }}[CRITIQUE]{{ else if eq .Alarm.Value.State.Value 2 }}[MAJEUR]{{ else if eq .Alarm.Value.State.Value 1 }}[MINEUR]{{ else }}[OK]{{ end }} Alerte sur {{ .Alarm.Value.Component }}{{ if .Alarm.Value.Resource }}/{{ .Alarm.Value.Resource }}{{ end }}
Content-Type: text/html; charset=UTF-8

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; color: #333; }
        .header { padding: 10px; border-radius: 5px 5px 0 0; }
        .header-critical { background-color: #FFE0E0; color: #D32F2F; border-bottom: 3px solid #D32F2F; }
        .header-major { background-color: #FFF3E0; color: #F57C00; border-bottom: 3px solid #F57C00; }
        .header-minor { background-color: #FFFDE7; color: #FBC02D; border-bottom: 3px solid #FBC02D; }
        .header-ok { background-color: #E8F5E9; color: #388E3C; border-bottom: 3px solid #388E3C; }
        .title { margin: 0; padding: 10px 0; font-size: 20px; }
        .content { padding: 20px; background-color: #f9f9f9; border: 1px solid #ddd; border-top: none; border-radius: 0 0 5px 5px; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th { text-align: left; background-color: #f2f2f2; }
        th, td { padding: 8px; border: 1px solid #ddd; }
        .message { background-color: white; padding: 10px; border: 1px solid #ddd; border-radius: 3px; margin-top: 20px; white-space: pre-wrap; }
        .footer { margin-top: 20px; font-size: 12px; color: #777; border-top: 1px solid #ddd; padding-top: 10px; }
        .button { display: inline-block; padding: 10px 15px; background-color: #1976D2; color: white; text-decoration: none; border-radius: 3px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="header {{ if eq .Alarm.Value.State.Value 3 }}header-critical{{ else if eq .Alarm.Value.State.Value 2 }}header-major{{ else if eq .Alarm.Value.State.Value 1 }}header-minor{{ else }}header-ok{{ end }}">
        <h1 class="title">{{ if eq .Alarm.Value.State.Value 3 }}⚠️ ALERTE CRITIQUE{{ else if eq .Alarm.Value.State.Value 2 }}⚠️ ALERTE MAJEURE{{ else if eq .Alarm.Value.State.Value 1 }}ℹ️ ALERTE MINEURE{{ else }}✅ RÉSOLUTION{{ end }}</h1>
    </div>
    <div class="content">
        <p>Une {{ if eq .Alarm.Value.State.Value 3 }}alerte critique{{ else if eq .Alarm.Value.State.Value 2 }}alerte majeure{{ else if eq .Alarm.Value.State.Value 1 }}alerte mineure{{ else }}résolution d'alerte{{ end }} a été détectée sur le composant <strong>{{ .Alarm.Value.Component }}</strong>{{ if .Alarm.Value.Resource }} / <strong>{{ .Alarm.Value.Resource }}</strong>{{ end }}.</p>
        
        <p>Cet email a été généré suite à un <strong>{{ .AdditionalData.Trigger }}</strong>.</p>
        
        <h2>Informations sur l'alarme</h2>
        <table>
            <tr>
                <th>ID</th>
                <td>{{ .Alarm.ID }}</td>
            </tr>
            <tr>
                <th>Composant</th>
                <td>{{ .Alarm.Value.Component }}</td>
            </tr>
            <tr>
                <th>Ressource</th>
                <td>{{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}N/A{{ end }}</td>
            </tr>
            <tr>
                <th>Criticité</th>
                <td>{{ if eq .Alarm.Value.State.Value 3 }}<span style="color:#D32F2F">CRITIQUE (3)</span>{{ else if eq .Alarm.Value.State.Value 2 }}<span style="color:#F57C00">MAJEUR (2)</span>{{ else if eq .Alarm.Value.State.Value 1 }}<span style="color:#FBC02D">MINEUR (1)</span>{{ else }}<span style="color:#388E3C">OK (0)</span>{{ end }}</td>
            </tr>
            <tr>
                <th>Statut</th>
                <td>{{ if eq .Alarm.Value.Status.Value 1 }}En cours{{ else if eq .Alarm.Value.Status.Value 0 }}Résolue{{ end }}</td>
            </tr>
            <tr>
                <th>Date de dernière mise à jour</th>
                <td>{{ .Alarm.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" "Europe/Paris" }}</td>
            </tr>
            <tr>
                <th>Source</th>
                <td>{{ .Alarm.Value.Connector }}/{{ .Alarm.Value.ConnectorName }}</td>
            </tr>
            <tr>
                <th>Événements</th>
                <td>{{ .Alarm.Value.EventsCount }}</td>
            </tr>
            <tr>
                <th>Changements d'état</th>
                <td>{{ .Alarm.Value.TotalStateChanges }}</td>
            </tr>
        </table>
        
        <h2>Message d'alarme</h2>
        <div class="message">{{ .Alarm.Value.Output }}</div>
        
        {{ if .Alarm.Value.LongOutput }}
        <h2>Détails supplémentaires</h2>
        <div class="message">{{ .Alarm.Value.LongOutput }}</div>
        {{ end }}
        
        {{ if .Alarm.Value.ACK }}
        <h2>Informations d'acquittement</h2>
        <table>
            <tr>
                <th>Auteur</th>
                <td>{{ .Alarm.Value.ACK.Author }}</td>
            </tr>
            <tr>
                <th>Message</th>
                <td>{{ if .Alarm.Value.ACK.Message }}{{ .Alarm.Value.ACK.Message }}{{ else }}Aucun message{{ end }}</td>
            </tr>
        </table>
        {{ end }}
        
        {{ if .Alarm.Value.Ticket }}
        <h2>Ticket associé</h2>
        <table>
            <tr>
                <th>Numéro</th>
                <td>{{ .Alarm.Value.Ticket.Ticket }}</td>
            </tr>
            <tr>
                <th>Auteur</th>
                <td>{{ .Alarm.Value.Ticket.Author }}</td>
            </tr>
            <tr>
                <th>Message</th>
                <td>{{ if .Alarm.Value.Ticket.Message }}{{ .Alarm.Value.Ticket.Message }}{{ else }}Aucun message{{ end }}</td>
            </tr>
        </table>
        {{ end }}
        
        <a href="{{ .Env.CanopsisURL }}/alarms/{{ .Alarm.ID }}" class="button">Consulter dans Canopsis</a>
        
        <div class="footer">
            <p>Email généré automatiquement par Canopsis le {{ .Alarm.Value.CreationDate | localtime "02/01/2006 à 15:04:05" "Europe/Paris" }}.</p>
            <p>Pour plus d'informations, connectez-vous à <a href="{{ .Env.CanopsisURL }}">Canopsis</a>.</p>
        </div>
    </div>
</body>
</html>
```

**Cas d'usage:** Parfait pour l'envoi de notifications par email aux équipes d'astreinte ou aux responsables. L'email est formaté pour être lisible et mettre en évidence les informations importantes, avec un code couleur selon la criticité.

**À noter:**

* Les variables `.Env.CanopsisEmailSender`, `.Env.AlertRecipients` et `.Env.CanopsisURL` doivent être définies
* L'utilisation de styles CSS inline pour garantir la compatibilité avec les clients email
* L'utilisation de `.AdditionalData.Trigger` pour indiquer le trigger qui a déclenché le scénario


## 6. Création de ticket ServiceNow

Ce template génère un payload JSON pour créer un ticket dans ServiceNow via leur API:

```go
{
  "sysparm_action": "insert",
  "table_name": "incident",
  "category": "monitoring",
  "subcategory": "canopsis",
  "short_description": "{{ if eq .Alarm.Value.State.Value 3 }}[CRITIQUE]{{ else if eq .Alarm.Value.State.Value 2 }}[MAJEUR]{{ else if eq .Alarm.Value.State.Value 1 }}[MINEUR]{{ end }} {{ .Alarm.Value.Component }}{{ if .Alarm.Value.Resource }}/{{ .Alarm.Value.Resource }}{{ end }}",
  "description": "{{ .Alarm.Value.Output | replace "\n" "\\n" }}",
  "work_notes": "{{ if .Alarm.Value.LongOutput }}{{ .Alarm.Value.LongOutput | replace "\n" "\\n" }}{{ else }}Pas de détails supplémentaires disponibles.{{ end }}",
  "cmdb_ci": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "cmdb_ci" }}{{ (index .Entity.Infos "cmdb_ci").Value }}{{ else }}{{ .Alarm.Value.Component }}{{ end }}{{ else }}{{ .Alarm.Value.Component }}{{ end }}",
  "impact": "{{ if eq .Alarm.Value.State.Value 3 }}1{{ else if eq .Alarm.Value.State.Value 2 }}2{{ else }}3{{ end }}",
  "urgency": "{{ if eq .Alarm.Value.State.Value 3 }}1{{ else if eq .Alarm.Value.State.Value 2 }}2{{ else }}3{{ end }}",
  "priority": "{{ if eq .Alarm.Value.State.Value 3 }}1{{ else if eq .Alarm.Value.State.Value 2 }}2{{ else if eq .Alarm.Value.State.Value 1 }}3{{ else }}4{{ end }}",
  "assignment_group": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "support_group" }}{{ (index .Entity.Infos "support_group").Value }}{{ else }}{{ .Env.DefaultSupportGroup }}{{ end }}{{ else }}{{ .Env.DefaultSupportGroup }}{{ end }}",
  "assigned_to": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "technical_contact" }}{{ (index .Entity.Infos "technical_contact").Value }}{{ end }}{{ end }}",
  "business_service": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "business_service" }}{{ (index .Entity.Infos "business_service").Value }}{{ end }}{{ end }}",
  "u_environment": "{{ .Env.Environment }}",
  "u_canopsis_alarm_id": "{{ .Alarm.ID }}",
  "u_canopsis_link": "{{ .Env.CanopsisURL }}/alarms/{{ .Alarm.ID }}",
  "caller_id": "{{ .Env.ServiceNowCaller }}",
  "comments": "Ticket créé automatiquement par Canopsis suite à une alerte de supervision"
}
```

**Cas d'usage:** Ce template permet de créer automatiquement des tickets d'incident dans ServiceNow lorsqu'une alarme est détectée. L'intégration est riche et utilise le maximum d'informations disponibles pour remplir correctement les champs du ticket.

**À noter:**

* Les champs impact, urgency et priority sont mappés automatiquement selon la criticité
* Le template exploite les informations de l'entité pour enrichir le ticket (assignation, CMDB, service)
* Les variables d'environnement `.Env.DefaultSupportGroup`, `.Env.Environment`, `.Env.CanopsisURL` et `.Env.ServiceNowCaller` doivent être définies
* La fonction `map_has_key` est utilisée pour vérifier que les informations sont disponibles avant de les utiliser

## 7. Création de ticket iTop

Ce template génère un payload pour la création d'un ticket dans iTop:

```go
{
  "operation": "core/create",
  "class": "Incident",
  "comment": "Ticket créé automatiquement depuis Canopsis",
  "fields": {
    "org_id": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "itop_organization" }}{{ (index .Entity.Infos "itop_organization").Value }}{{ else }}{{ .Env.DefaultITopOrg }}{{ end }}{{ else }}{{ .Env.DefaultITopOrg }}{{ end }}",
    "title": "{{ if eq .Alarm.Value.State.Value 3 }}[CRITIQUE]{{ else if eq .Alarm.Value.State.Value 2 }}[MAJEUR]{{ else if eq .Alarm.Value.State.Value 1 }}[MINEUR]{{ end }} {{ .Alarm.Value.Component }}{{ if .Alarm.Value.Resource }}/{{ .Alarm.Value.Resource }}{{ end }}",
    "description": "## Message d'alarme\n\n{{ .Alarm.Value.Output | replace "\n" "\\n" }}\n\n{{ if .Alarm.Value.LongOutput }}## Détails supplémentaires\n\n{{ .Alarm.Value.LongOutput | replace "\n" "\\n" }}{{ end }}\n\n## Informations techniques\n\n- **ID d'alarme**: {{ .Alarm.ID }}\n- **Composant**: {{ .Alarm.Value.Component }}\n- **Ressource**: {{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}N/A{{ end }}\n- **Connecteur**: {{ .Alarm.Value.Connector }}/{{ .Alarm.Value.ConnectorName }}\n- **Date de détection**: {{ .Alarm.Value.CreationDate | localtime "02/01/2006 15:04:05" }}\n- **Dernière mise à jour**: {{ .Alarm.Value.LastUpdateDate | localtime "02/01/2006 15:04:05" }}\n- **Nombre d'événements**: {{ .Alarm.Value.EventsCount }}\n- **Changements d'état**: {{ .Alarm.Value.TotalStateChanges }}",
    "impact": "{{ if eq .Alarm.Value.State.Value 3 }}3{{ else if eq .Alarm.Value.State.Value 2 }}2{{ else }}1{{ end }}",
    "urgency": "{{ if eq .Alarm.Value.State.Value 3 }}3{{ else if eq .Alarm.Value.State.Value 2 }}2{{ else }}1{{ end }}",
    "origin": "monitoring",
    "service_id": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "itop_service" }}{{ (index .Entity.Infos "itop_service").Value }}{{ else }}{{ .Env.DefaultITopService }}{{ end }}{{ else }}{{ .Env.DefaultITopService }}{{ end }}",
    "team_id": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "itop_team" }}{{ (index .Entity.Infos "itop_team").Value }}{{ else }}{{ .Env.DefaultITopTeam }}{{ end }}{{ else }}{{ .Env.DefaultITopTeam }}{{ end }}",
    "functionalcis_list": [
      {
        "functionalci_id": "{{ if .Entity.Infos }}{{ if map_has_key .Entity.Infos "itop_ci" }}{{ (index .Entity.Infos "itop_ci").Value }}{{ else }}{{ .Alarm.Value.Component }}{{ end }}{{ else }}{{ .Alarm.Value.Component }}{{ end }}",
        "impact_code": "not_impacted"
      }
    ],
    "caller_id": "{{ .Env.ITopCaller }}",
    "private_log": "Ticket créé automatiquement par Canopsis suite à une alarme de supervision.\nLien vers l'alarme: {{ .Env.CanopsisURL }}/alarms/{{ .Alarm.ID }}"
  }
}
```

**Cas d'usage:** Ce template est idéal pour créer des tickets d'incident dans iTop avec toutes les informations nécessaires pour traiter le problème. Il utilise le formatage markdown pour la description du ticket, ce qui améliore la lisibilité dans iTop.

**À noter:**

* Le template fait correspondre la criticité Canopsis aux niveaux d'impact et d'urgence d'iTop
* Les informations de l'entité sont utilisées pour enrichir le ticket avec les valeurs correctes de l'organisation, du service, de l'équipe et du CI
* Les variables d'environnement `.Env.DefaultITopOrg`, `.Env.DefaultITopService`, `.Env.DefaultITopTeam`, `.Env.ITopCaller` et `.Env.CanopsisURL` doivent être définies
* La fonction `map_has_key` est utilisée pour vérifier que les informations sont disponibles avant de les utiliser

## 8. Notification pour webhook générique

Template pour une notification à envoyer vers un webhook générique:

```go
{
  "alert": {
    "id": "{{ .Alarm.ID }}",
    "title": "{{ if eq .Alarm.Value.State.Value 3 }}[CRITIQUE]{{ else if eq .Alarm.Value.State.Value 2 }}[MAJEUR]{{ else if eq .Alarm.Value.State.Value 1 }}[MINEUR]{{ else }}[OK]{{ end }} {{ .Alarm.Value.Component }}{{ if .Alarm.Value.Resource }}/{{ .Alarm.Value.Resource }}{{ end }}",
    "message": "{{ .Alarm.Value.Output }}",
    "details": "{{ if .Alarm.Value.LongOutput }}{{ .Alarm.Value.LongOutput }}{{ else }}Pas de détails supplémentaires.{{ end }}",
    "source": {
      "component": "{{ .Alarm.Value.Component }}",
      "resource": "{{ if .Alarm.Value.Resource }}{{ .Alarm.Value.Resource }}{{ else }}null{{ end }}",
      "connector": "{{ .Alarm.Value.Connector }}",
      "connector_name": "{{ .Alarm.Value.ConnectorName }}"
    },
    "status": {
      "severity": {{ .Alarm.Value.State.Value }},
      "severity_label": "{{ if eq .Alarm.Value.State.Value 3 }}CRITICAL{{ else if eq .Alarm.Value.State.Value 2 }}MAJOR{{ else if eq .Alarm.Value.State.Value 1 }}MINOR{{ else }}OK{{ end }}",
      "status": {{ .Alarm.Value.Status.Value }},
      "status_label": "{{ if eq .Alarm.Value.Status.Value 1 }}ACKNOWLEDGED{{ else if eq .Alarm.Value.Status.Value 2 }}IN_PROGRESS{{ else }}NEW{{ end }}"
    },
    "timestamps": {
      "creation": "{{ .Alarm.Value.CreationDate }}",
      "last_update": "{{ .Alarm.Value.LastUpdateDate }}",
      "last_event": "{{ .Alarm.Value.LastEventDate }}",
      "activation": "{{ if .Alarm.Value.ActivationDate }}{{ .Alarm.Value.ActivationDate }}{{ else }}null{{ end }}",
      "resolved": "{{ if .Alarm.Value.Resolved }}{{ .Alarm.Value.Resolved }}{{ else }}null{{ end }}"
    },
    "counts": {
      "events": {{ .Alarm.Value.EventsCount }},
      "state_changes": {{ .Alarm.Value.TotalStateChanges }},
      "changes_since_update": {{ .Alarm.Value.StateChangesSinceStatusUpdate }}
    },
    "entity": {
      "id": "{{ .Entity.ID }}",
      "name": "{{ .Entity.Name }}"
    },
    "acknowledgement": {
      {{ if .Alarm.Value.ACK }}
      "author": "{{ .Alarm.Value.ACK.Author }}",
      "message": "{{ if .Alarm.Value.ACK.Message }}{{ .Alarm.Value.ACK.Message }}{{ else }}null{{ end }}"
      {{ else }}
      "author": null,
      "message": null
      {{ end }}
    },
    "ticket": {
      {{ if .Alarm.Value.Ticket }}
      "id": "{{ .Alarm.Value.Ticket.Ticket }}",
      "author": "{{ .Alarm.Value.Ticket.Author }}",
      "message": "{{ if .Alarm.Value.Ticket.Message }}{{ .Alarm.Value.Ticket.Message }}{{ else }}null{{ end }}"
      {{ else }}
      "id": null,
      "author": null,
      "message": null
      {{ end }}
    },
    "change_type": "{{ .AdditionalData.Trigger }}",
    "canopsis_url": "{{ .Env.CanopsisURL }}/alarms/{{ .Alarm.ID }}",
    "environment": "{{ .Env.Environment }}"
  }
}
```

**Cas d'usage:** Ce template crée un payload JSON complet et bien structuré pour l'intégration avec n'importe quel système externe pouvant recevoir des webhooks. Il contient toutes les informations importantes sur l'alarme dans un format facile à traiter.

**À noter:**

* Le format JSON est structuré de manière logique pour faciliter l'intégration
* Le template gère correctement les valeurs nulles pour éviter les erreurs de parsing
* Les variables d'environnement `.Env.CanopsisURL` et `.Env.Environment` doivent être définies
* Le type de trigger est inclus via `.AdditionalData.Trigger`


