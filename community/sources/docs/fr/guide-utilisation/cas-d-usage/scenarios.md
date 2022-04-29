# Scenarios

## Période de confirmation pour les nouvelles alarmes

Certaines sources de données peuvent remonter des alarmes qui ont tendance à se résoudre sans intervention au bout d'un certain temps.
Canopsis permet de trier ces faux positifs.

### Configuration

Rendez vous dans l'interface de gestion des *scenarios* et créez en un nouveau.

![Nouveau scenario](./img/scn_snooze_new_scn.png)

Configurez le comme ci-dessous avec comme alarm pattern :
```json
{
    "v": {
        "state": {
            "val": {
                ">": 0
            }
        }
    }
}
```

![Configuration scenario](./img/scn_snooze_configuration.png)

Lors de la réception d'une alarme, elle sera automatiquement ignorée pendant 5 minutes. Ce délai peut permettre à un évènement résolvant cette alarme d'arriver.

Si la configuration fonctionne vous devriez voir cet icône ![Icone snooze](img/scn_snooze_result.png) sur les alarme arrivant dans le bac.

## Création de tickets dans Itop à la récéption d'une alarme

### Configuration

Créez le scenario suivant dans Canopsis:

!!! info Note
	Pensez à mettre à jour l'URL ainsi que les paramètres d'authentification pour qu'ils correspondent à votre instance Itop.

<details>

<summary>Requête CURL pour envoyer la configuration à Canopsis.</summary>
```bash
curl -X POST -u root:root -H "Content-type: application/json" -d '{
	"name" : "create_itop_ticket",
	"author" : "root",
	"enabled" : true,
	"disable_during_periods" : [ ],
	"triggers" : [
		"create"
	],
	"actions" : [
		{
			"type" : "webhook",
			"comment" : "",
			"parameters" : {
				"declare_ticket" : {
					"empty_response" : false,
					"is_regexp" : true,
					"ticket_id" : "objects\\.UserRequest::.*\\.fields\\.friendlyname"
				},
				"request" : {
					"auth" : {
						"username" : "admin",
						"password" : "ChAtX713IHw8"
					},
					"headers" : {
						"Content-type" : "application/x-www-form-urlencoded"
					},
					"method" : "POST",
					"payload" : "json_data={\n  \"operation\":\"core/create\",\n  \"comment\":\"Alarm created by Canopsis\",\n  \"class\":\"UserRequest\",\n  \"output_fields\":\"id, friendlyname\",\n  \"fields\":\n  {\n    \"org_id\":\"SELECT Organization WHERE name = \\\"Demo\\\"\",\n    \"title\":\"Alarm on : {{ .Alarm.Value.Component }} {{ .Alarm.Value.Resource }}\",\n    \"description\":\"Message : {{ .Alarm.Value.State.Message }}\",\n    \"functionalcis_list\" : [{\"functionalci_id\":\"SELECT Server WHERE name=\\\"{{ .Alarm.Value.Component}}\\\"\"}]\n  }\n}",
					"skip_verify" : true,
					"url" : "http://itop/webservices/rest.php?version=1.3&login_mode=basic"
				},
				"retry_count" : 3,
				"retry_delay" : {
					"unit" : "m",
					"value" : 1
				}
			},
			"alarm_patterns" : [
				{
					"v" : {
						"state" : {
							"val" : {
								">" : 0
							}
						}
					}
				}
			],
			"entity_patterns" : null,
			"drop_scenario_if_not_matched" : false,
			"emit_trigger" : false
		}
	],
	"priority" : 3,
	"delay" : null
}' 'http://localhost:8082/api/v4/scenarios'
```

</details>

![Configuration Webhook ITOP](./img/scn_itop_config.png)



Lors de la réception d'une alarme, un ticket sera automatiquement créé sur Itop.

Vous devriez voir apparaitre un ticket dans Itop :

![Ticket dans l'interface Itop](./img/scn_itop_ticket.png)


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
