# Connecteur PRTG

## Description

Convertit les alertes de PRTG en évènements Canopsis

## Principe de fonctionnement

Les notifications de PRTG peuvent être configurées afin d'effectuer une action HTTP.

En utilisant les [placeholders](https://kb.paessler.com/en/topic/373-what-placeholders-can-i-use-with-prtg) fournis par PRTG nous pouvons contruire un paquet qui sera ensuite envoyé à l'API de Canopsis

Les statuts renvoyés par PRTG ne sont pas dans le format attendu par Canopsis et seront transformés en utilisant des règles d'enrichissement.

### Traduction des états

La traduction des états entre PRTG et Canopsis est la suivante :

| PRTG (FR)          | PRTG (EN)      | Canopsis     |
|--------------------|----------------|--------------|
| Erreur             | Down           | MAJOR (2)    |
| Avertissement      | Warning        | MINOR (1)    |
| Inhabituel         | Unusual        | MINOR (1)    |
| Erreur (partielle) | Down (Partial) | MAJOR (2)    |
| OK                 | Up             | INFO (0)     |
| Inconnu            | Unknown        | CRITICAL (3) |

Ces correspondances sont données à titre indicatif, lors de la création des règles d'enrichissement vous pourrez choisir les correspondances que vous souhaitez utiliser.

## Configuration de PRTG

Cet exemple de configuration est fait avec une interface utilisateur de PRTG en français. Si votre PRTG est dans une autre langue vous devez adapter les exemples donnés en conséquence.

Il faut tout d'abord `Ajouter un modèle de notification` depuis le menu `Modèles de notifications`

Pour accéder à ce menu aller dans `Configuration` puis `Paramètres de compte` et enfin `Modèles de notifications`

![Menu](img/PRTG_notifications.png)

Dans la partie `Paramétrages de base` donnez un nom à votre modèle.

Dans la partie `Résumé des notifications` cochez la puce `Toujours aviser le plus tôt possible, ne jamais résumer`

Il faudra ensuite cocher le bouton `Exécuter une action HTTP` et remplir les champs comme ci-dessous : 

![Configuration action HTTP](img/PRTG_Action_HTTP.png)

Le champ `cargaison` est le suivant :

```
event_type=check&connector=PRTG&connector_name=PRTG&component=%host&resource=%shortname&source_type=resource&state=3&prtg_state=%laststatus&output=%message
```

Il ne vous reste plus qu'à utiliser ce modèle dans vos `Déclencheurs de notifications`

![Déclencheurs de notifications](img/PRTG_declencheur.png)

## Configuration de Canopsis

Maintenant que PRTG envoie ses alertes vers Canopsis il faut les enrichir afin qu'elles puissent être correctement traitées par Canopsis.

Il faut pour cela créer une première règle d'enrichissment qui permet d'activer l'enrichissement

![Options pour la copie de l'entité](img/PRTG_canopsis_entity_copy_options.png)

![Données externe pour la copie de l'entité](img/PRTG_canopsis_entity_copy_external_data.png)

![Action pour la copie de l'entité](img/PRTG_canopsis_entity_copy_action.png)

Puis ajouter une règle afin d'enrichir le contexte de l'entité avec le `prtg_state`

![Options pour l'enrichessement de l'entité](img/PRTG_canopsis_context_enrich_options.png)

![Action pour l'enrichissement de l'entité](img/PRTG_canopsis_context_enrich_action.png)

Une fois ces deux règles créées il reste à créer une règle par statut et par langue de PRTG

![Options pour l'enrichissement de l'alarme](img/PRTG_canopsis_entity_enrich_options.png)

![Action pour l'enrichissement de l'alarme](img/PRTG_canopsis_entity_enrich_action.png)
