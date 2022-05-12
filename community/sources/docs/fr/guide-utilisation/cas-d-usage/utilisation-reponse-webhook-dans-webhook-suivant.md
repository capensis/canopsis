# Utiliser la réponse d'un webhook dans le webhook suivant

## Exemple

### Contexte

Dans notre exemple nous allons utiliser une API factice : Image docker : `docker.canopsis.net/docker/community/fakeapi:latest` .

Notre scénario comportera deux webhooks :

#### 1 - Un premier webhook qui va appeler l'api avec l'action `getObjectId` pour récupérer l'id de l'objet.

fakeapi retourne le JSON ci dessous pour ce premier appel:
```json
{"id": 42}
```
#### 2 - Un second webhook qui va appeler l'api avec l'action `getObjectById` avec en paramètre l'id récupéré a l'étape précédente.

dans le corps du deuxième webhook, nous utilisons le format `{{index .Response "JSON_KEY"}}` pour récupérer l'information dans le corps du webhook précédent (ici `{{index .Response "id"}}` pour récupérer la clé `id`).
!!! note Note
	Nous aurions aussi pu utiliser le format `{{index .Header "Header-Name"}}` pour récupérer le contenu d'un header de la réponse du webhook précédent.

### Démarrage de l'api

!!! info Attention
	Attention vous devez avoir démarré une instance Canopsis sur la même machine pour réaliser cet exemple.

```
docker run --network=canopsis-pro_default --rm --name fakeapi docker.canopsis.net/docker/community/fakeapi:latest
```

### Configuration

Envoyez la configuration à Canopsis avec cette requête :
```
curl -X POST -u root:root -H "Content-type: application/json" -d '{
	"name": "MultiStepScenario",
	"author": "root",
	"enabled": true,
	"disable_during_periods": [],
	"triggers": [
		"create",
		"changestate"
	],
	"actions": [
		{
			"type": "webhook",
			"comment": "Step one of the api calls : Get the object ID",
			"parameters": {
				"declare_ticket": null,
				"request": {
					"auth": null,
					"headers": {},
					"method": "GET",
					"payload": "{\"action\": \"getObjectId\"}",
					"skip_verify": true,
					"url": "http://fakeapi/main.php"
				},
				"retry_count": 3,
				"retry_delay": {
					"unit": "s",
					"value": 10
				}
			},
			"alarm_patterns": [
				{
					"v": {
						"state": {
							"val": {
								">": 0
							}
						}
					}
				}
			],
			"entity_patterns": null,
			"drop_scenario_if_not_matched": true,
			"emit_trigger": false
		},
		{
			"type": "webhook",
			"comment": "Step 2, get the object by id ",
			"parameters": {
				"declare_ticket": {
					"empty_response": false,
					"is_regexp": false,
					"ticket_id": "ticket_id"
				},
				"request": {
					"auth": null,
					"headers": {},
					"method": "GET",
					"payload": "{\"action\": \"getObjectById\", \"id\": {{ index .Response \"id\"}} }",
					"skip_verify": false,
					"url": "http://fakeapi/main.php"
				},
				"retry_count": 3,
				"retry_delay": {
					"unit": "s",
					"value": 10
				}
			},
			"alarm_patterns": [
				{
					"v": {
						"state": {
							"val": {
								">": 0
							}
						}
					}
				}
			],
			"entity_patterns": null,
			"drop_scenario_if_not_matched": false,
			"emit_trigger": false
		}
	],
	"priority": 1,
	"delay": null,
	"created": 1651156396,
	"updated": 1651156396
}' 'http://localhost:8082/api/v4/scenarios'
```

Lors de l'envoi d'un évènement un ticket avec l'id `R000042` devrait lui être assigné. 
