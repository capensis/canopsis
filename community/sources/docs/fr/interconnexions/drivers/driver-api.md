# Driver API (`import-context-graph`)

!!! Info
    Disponible uniquement en édition Pro.

Canopsis embarque un programme permettant d'interroger une API externe afin de compléter son référentiel interne. Il peut être exécuté grâce à une ligne de commande ou via un conteneur Docker.

Ce programme se situe :

* à l'emplacement `/opt/canopsis/bin/import-context-graph` lors d'une installation par paquets.
* dans le conteneur `pro/import-context-graph` en Docker.

#### Table des matières
1. [Options](#options)<br>
2. [Variables d'environnement](#variables-denvironnement)<br>
3. [Configuration](#configuration)<br>
4. [Exemple](#exemple)<br>
5. [Exécution](#execution)<br>
6. [Résultats dans Canopsis](#resultats-dans-canopsis)

## Utilisation

### Options

| Option  | Argument                 | Description                                           |
|---------|--------------------------|-------------------------------------------------------|
| `-help` |                          | Lister toutes les options acceptées                   |
| `-d`    |                          | Activer le mode debug                                 |
| `-c`    | `/chemin/du/fichier.yml` |Indiquer le chemin complet du fichier de configuration |

### Variables d'environnement

L'identifiant et le mot de passe de connexion à l'API sont définis via des variables d'environnement :

 * `EXTERNAL_API_USERNAME` pour l'identifiant
 * `EXTERNAL_API_PASSWORD` pour le mot de passe

Seule l'authentification [basique](https://fr.wikipedia.org/wiki/Authentification_HTTP#M%C3%A9thode_%C2%AB_Basic_%C2%BB){target=_blank} est supportée.

### Configuration

Le format de fichier de configuration est le YAML. Il doit être entièrement rédigé à partir du JSON retourné par l'API externe.

Un exemple de fichier de configuration est disponible sur le dépôt [Canopsis Pro](https://git.canopsis.net/canopsis/canopsis-pro/-/tree/develop/pro/go-engines-pro/config/import-context-graph/api.yml.example){target=_blank}.
Un exemple donné à titre indicatif est également disponible [ci-après](#exemple).

Ce fichier de configuration doit comprendre trois sections :

 1. `api` qui permet d'indiquer les propriétés de la requête HTTP à envoyer à l'API.
 2. `import` qui permet de spécifier les paramètres de l'import.
 3. `mapping` qui permet de spécifier les associations entre les champs de l'entité dans Canopsis et la réponse de l'API.

Concernant les types d'actions d'import :

  *  `action` : Action effectuée sur l'entité reçue dans la réponse de l'API
     *  `create` : Création de l'entité en base
     *  `set`: La même chose que `create` mais peut être partielle si l'API ne retourne pas tous les champs
     *  `update` : Mise à jour de l'entité en base. Si l'entité n'existe pas, aucune modification ne sera appliquée
     *  `enable` : Active l'entité (passage à `true` du champ `enable` de l'entité en base)
  *  `missing_action` : Action effectuée sur les entités manquantes de la réponse de l'API
     *  `delete` : Suppression de l'entité en base
     *  `disable` : Désactive l'entité (action inverse de l'action `enable`)
     *  (chaîne vide) : le paramètre `missing_action` peut ne pas avoir d'action et donc ne rien faire

Pour plus de détails, référez-vous à la [documentation d'import de contexte](https://doc.canopsis.net/guide-developpement/swagger/#/contextgraph-import){target=_blank} disponible dans Swagger.

### Exemple

Prenons l'exemple d'une API externe retournant le JSON ci-dessous :

``` json
[
    {
        "ci": "composant1",
        "nom": "composant 1",
        "societe": "CAPENSIS",
        "statut": "actif",
        "responsable": "Service technique",
        "localisation": "DC1",
        "commentaire": "Il s agit du composant 1",
        "impact": 5
    }
]
```

Le fichier de configuration avec les commentaires explicatifs :

``` yaml
# Ce bloc permet de spécifier les propriétés de la requête HTTP à envoyer à l'API externe.
api:
  # Spécifier l'URL de l'API
  url: http://mon.api/?filtre=filtre1
  # Spécifier le type de requête HTTP (GET, POST, PUT etc)
  method: GET
  # Spécifier les en-têtes de la requête
  headers:
  # Le payload contient le body
  body:
    key: value
  # Activer (false) ou désactiver (true) la vérification de la chaîne de certification TLS du serveur
  insecure_skip_verify: false

# Ce bloc spécifie les paramètres de l'import
import:
  # Définit la source de l'import à l'entité
  source: import-context-graph
  # Définit le type d'action à effectuer sur les entités reçues de la réponse de l'API
  # Valeurs possibles : create, set, update, enable
  action: create
  # Définit le type d'action à effectuer sur les entités manquantes de la réponse de l'API
  # Les entités manquantes sont trouvées par source d'import
  # Valeurs possibles : delete, disable. Si rien n’est défini, les entités manquantes ne seront pas mises à jour.
  missing_action: disable

# Ce bloc spécifie l'association entre les champs de l'entité et la réponse de l'API
# Le champ _id est requis, les autres champs sont optionnels
# Seul le JSON qui contient un tableau d'objets est pris en charge
mapping:
  # path spécifie le chemin d’où importer les données
  # Par exemple, "data.nested" si la réponse est {"data": {"nested": [ {...}, {...} ]}}.
  # Peut être vide.
  path:
  # is_map spécifie le type des données importées.
  # Devrait être false si la réponse est array [ {...}, {...} ].
  # Devrait être true si la réponse est map {"key1": {...}, "key2": {...}}.
  is_map: false
  # Association de l'ID de l'entité
  _id: ci
  # Association de la description
  description: commentaire
  # Association du niveau d'impact
  impact_level: impact
  # Association d'informations complémentaires
  infos:
    nom:
      value: nom
      description: nom
    societe:
      value: societe
      description: societe
    statut:
      value: statut
      description: statut
    responsable:
      value: responsable
      description: responsable
```

Exemple pour ITop :
``` yaml
api:
  url: http://mon.itop/webservices/rest.php?version=1.3
  method: POST
```

### Exécution

!!! attention
    Ces exemples d'exécution ne répondent pas aux bonnes pratiques de sécurité. Veillez donc à bien adapter cette exécution selon votre politique de sécurité interne.
    
    Dans le cadre d'un usage via Docker, il est conseillé d'utiliser "[Docker Secrets](https://docs.docker.com/engine/swarm/secrets/){target=_blank}" ou une autre solution de coffre-fort.

#### Programme installé par packets :

``` shell
export EXTERNAL_API_USERNAME=test
export EXTERNAL_API_PASSWORD=test
/opt/canopsis/bin/import-context-graph -c import-context-api.yml 
```
Retour :
``` shell
2021-09-29T10:39:46+02:00 INF git.canopsis.net/canopsis/canopsis-pro/pro/go-engines-pro/cmd/import-context-graph/main.go:65 > import finished deleted=0 exec_time=3.784369ms updated=1
```

#### Programme dans Docker :

``` shell
docker run \
-e EXTERNAL_API_USERNAME='[testuser]' \
-e EXTERNAL_API_PASSWORD='[testpassword]' \
-e CPS_MONGO_URL='mongodb://cpsmongo:canopsis@mongodb/canopsis' \
-e 'CPS_AMQP_URL=amqp://cpsrabbit:canopsis@rabbitmq/canopsis' \
-e CPS_POSTGRES_URL='postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis' \
--network=canopsis-pro_default -it --rm -v \
"[/chemin/vers]/pro/deployment/canopsis/docker/files/api.yml:/opt/canopsis/share/config/import-context-graph/api.yml" \
docker.canopsis.net/docker/pro/import-context-graph:'[4.5]'
```
Retour :
``` shell
2022-02-18T14:10:15Z INF git.canopsis.net/canopsis/canopsis-pro/pro/go-engines-pro/cmd/import-context-graph/main.go:78 > import finished deleted=0 exec_time=16.775252ms updated=3
```

### Résultats dans Canopsis :

![](./img/imported_entity.png)

![](./img/linked_alarm.png)

![](./img/var_alarm.png)
