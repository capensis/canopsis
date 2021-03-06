# Événement

L'API event permet d'envoyer un [événement](../struct-event.md#structure-basique-dun-evenement) sur l'exchange de Canopsis dans le [bus AMQP](../../struct-event/#focus-amqp).

!!! note
    La concaténation des champs `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]` générés par vos évènements ne doit pas dépasser 256 caractères, sans quoi l'évènement ne pourra pas être traité par Canopsis.
    
    Consultez la documentation sur les [limitations des évènements](../../guide-utilisation/limitations/index.md#limitations-des-evenements] pour en savoir plus.

## Envoi d'un événement

**URL** : `/api/v2/event`

**Méthode** : `POST`

**Headers** : `Content-Type: application/json`

**Authentification requise** : Oui

**Permissions requises** : Aucune

**Exemple de corps de requête** :
```json
{
"event_type": "check",
"connector": "superviseur",
"connector_name": "superviseur_de_Paris",
"component": "serveur_de_salle_machine_DHCP",
"resource": "disk2",
"source_type": "resource",
"author": "superviseur1",
"state": 2,
"debug": true,
"output": "Disque plein a 98%, 50GO occupe"
}
```

**Exemple de requête curl** pour utilisateur `root` avec mot de passe `root` qui veut ajouter le JSON ci-dessus :

```sh
curl -X POST -u root:root -H "Content-Type: application/json" -d '{
"event_type": "check",
"connector": "superviseur",
"connector_name": "superviseur_de_Paris",
"component": "serveur_de_salle_machine_DHCP",
"resource": "disk2",
"source_type": "resource",
"author": "superviseur1",
"state": 2,
"debug": true,
"output": "Disque plein a 98%, 50GO occupe"
}' 'http://localhost:8082/api/v2/event'
```

### Réponse en cas de réussite

**Condition** : L'événement a a été envoyé avec succès.

**Code** : `200 OK`

**Exemple du corps de la réponse** :

```javascript
{
  "sent_events": [
    {
      "resource": "disk2",
      "event_type": "check",
      "author": "superviseur1",
      "component": "serveur_de_salle_machine_DHCP",
      "connector": "superviseur",
      "source_type": "resource",
      "state": 2,
      "connector_name": "superviseur_de_Paris",
      "debug": true,
      "output": "Disque plein a 98%, 50GO occupe"
    }
  ],
  "retry_events": [],
  "failed_events": []
}
```
