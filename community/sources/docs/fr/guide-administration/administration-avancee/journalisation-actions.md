# Journalisation des actions utilisateurs

Canopsis permet de journaliser certaines actions réalisées par un utilisateur dans un journal.

## Actions et types pris en charge

!!! attention "Avertissement"
    En version Canopsis 4.1, seuls les types `pbehaviortype`, `pbehaviorreason`, `pbehaviorexception`, `pbehavior`, `heartbeat`, `jobconfig`, `job`, `instruction` sont traités. 
    Les autres types seront pris en charge dans les versions à venir de Canopsis

Actions  | Description  
--|---
`create`  | Création d'un objet (ex : Création d'un groupe de vues)
`update`  | Mise à jour d'un objet (ex : Mise à jour d'une règle Heartbeat)
`delete`  | Suppression d'un objet (ex : Suppression d'une règle de corrélation)
`export`  | Export d'un objet (ex : Export d'une vue)
`import`  | Import d'un objet (ex : Import d'une entité)

Types  | Description  
--|---
`view`                  | [Les vues](../../guide-utilisation/interface/vues/index.md)
`viewgroup`             | [Les groupes de vues](../../guide-utilisation/interface/vues/index.md)
`eventfilter`           | [Les règles de filtrage/enrichissement](../../guide-utilisation/menu-exploitation/filtres-evenements.md)
`metaalarmrule`         | [Les règles de méta alarmes/corrélation](../../guide-utilisation/menu-exploitation/regles-metaalarme.md)
`dynamicinfo`           | [Les règles d'enrichissement d'alarmes](../../guide-utilisation/cas-d-usage/enrichissement.md)
`service`               | [Les entités de type `service`](../../guide-utilisation/services/index.md)
`pbehavior`<br/>`pbehaviortype`<br/>`pbehaviorreason`<br/>`pbehaviorexception`  | [Les comportements périodiques](../../guide-utilisation/cas-d-usage/comportements_periodiques.md)
`heartbeat`             | [Les règles de lignes de vie](../../guide-utilisation/menu-exploitation/regles-inactivite.md)
`instruction`<br/>`jobconfig`             | [Les objets de remédiation](../../guide-utilisation/remediation/index.md)


## Exploitation des journaux

### Consultation

Les actions sont consignées dans le journal de l'API Canopsis. Voici la méthode qui vous permet de le consulter en fonction de votre type d'installation

=== "Paquets"

    ```sh
    journalctl -u canopsis-service@canopsis-api.service
    ```

=== "Docker Compose Community"

    ```sh
    CPS_EDITION=community docker compose logs api
    ```

=== "Docker Compose Pro"

    ```sh
    CPS_EDITION=pro docker compose logs api
    ```


### Anatomie du journal

Une ligne de journal présente ainsi : une date, une action, un type, et un identifiant d'objet.

```
févr. 10 09:18:09 localhost.localdomain canopsis-api[2699]: 2021-02-10T09:18:09+01:00 INF root/canopsis/go-engines/lib/api/logger/action_logger.go:86 > ActionLog:  action=create author=root value_id=2c2a146f-0861-411b-ac5d-02e153514c0c value_type=pbehaviorreason
```

## Collection MongoDB des dernières actions

Par ailleurs, une collection mongodb, `action_log`, propose de conserver la dernière action effectuée sur un objet.

Cette collection sera exploitée par l'interface web de Canopsis dans une future version.

```
MongoDB shell version v3.6.8
connecting to: mongodb://localhost:27017/canopsis
Implicit session: session { "id" : UUID("cff59761-6f70-4d8a-b592-d33d13f43ec5") }
MongoDB server version: 3.6.21
> db.action_log.find().pretty()

{
	"_id" : ObjectId("602396c113fa8b223784d93c"),
	"value_id" : "2c2a146f-0861-411b-ac5d-02e153514c0c",
	"value_type" : "pbehaviorreason",
	"action" : "create",
	"author" : "root",
	"time" : ISODate("2021-02-10T08:20:35.806Z")
}
```
