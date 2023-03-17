# Driver API (`import-context-graph`)

!!! Info
    Disponible uniquement en édition Pro.

Canopsis embarque un programme permettant d'interroger une API externe afin de compléter son référentiel interne. Il peut être exécuté grâce à une ligne de commande ou via un conteneur Docker.

Ce programme se situe :

* à l'emplacement `/opt/canopsis/bin/import-context-graph` lors d'une installation par paquets.
* dans le conteneur `pro/import-context-graph` en Docker.

Le `import-context-graph` peut importer 

* Des composants
* Des ressources
* Des services

#### Table des matières
1. [Options](#options)<br>
2. [Variables d'environnement](#variables-denvironnement)<br>
3. [Configuration](#configuration)<br>
4. [Importer des entités](#importer-des-entités)<br>
5. [Importer des informations](#importer-des-informations-complementaires)<br>
6. [Exécution](#execution)<br>
7. [Résultats dans Canopsis](#resultats-dans-canopsis)

## Options

| Option  | Argument                 | Description                                           |
|---------|--------------------------|-------------------------------------------------------|
| `-help` |                          | Lister toutes les options acceptées                   |
| `-d`    |                          | Activer le mode debug                                 |
| `-c`    | `/chemin/du/fichier.yml` |Indiquer le chemin complet du fichier de configuration |

## Variables d'environnement

L'identifiant et le mot de passe de connexion à l'API sont définis via des variables d'environnement :

 * `EXTERNAL_API_USERNAME` pour l'identifiant
 * `EXTERNAL_API_PASSWORD` pour le mot de passe

Seule l'authentification [basique](https://fr.wikipedia.org/wiki/Authentification_HTTP#M%C3%A9thode_%C2%AB_Basic_%C2%BB){target=_blank} est supportée.

Si votre API ne nécessite pas d'authentification, il vous suffit de ne pas spécifier ces varaibles d'environnement.

## Configuration

Le format de fichier de configuration est le YAML. Il doit être entièrement rédigé à partir du JSON retourné par l'API externe.

Un exemple de fichier de configuration est disponible sur le dépôt [Canopsis Pro](https://git.canopsis.net/canopsis/canopsis-pro/-/tree/develop/pro/go-engines-pro/config/import-context-graph/api.yml.example){target=_blank}.
Un exemple donné à titre indicatif est également disponible [ci-après](#exemple).

Ce fichier de configuration doit comprendre trois sections :

 1. `api` qui permet d'indiquer les propriétés de la requête HTTP à envoyer à l'API.
 2. `import` qui permet de spécifier les paramètres de l'import.
 3. `mapping` qui permet de spécifier les associations entre les champs de l'entité dans Canopsis et la réponse de l'API.


### Section `api`

Vous définissez ici l'URL, la méthode, les entêtes, et le corps de la requête vers l'API tierce.

| Paramètre            | Valeur                                                                    |
| -------------------- | ------------------------------------------------------------------------- |
| url                  | URL de l'API à interroger                                                 | 
| method               | Méthode HTTP (GET, POST, PUT, etc)                                        |
| headers              | Entêtes HTTP (ex : Content-type: application/json)                        |
| body                 | Corps de la requête sous forme de clé/valeur (ex : key: value)            |
| insecure_skip_verify | Si `true` alors vérification de la chaine de certification et du hostname |

### Section `import`

| Paramètre            | Valeur                                                                                                                 |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| source         | Permet de spécifier le nom de cet import. Ce nom sera présent dans les entités importées                                     |
| action         | Vaut `set` ou `enable`                                                                                                       |
| missing_action | Vaut `delete` ou `disable`. Action exécutée lorsque les entités spécifiées par "source" sont absentes de la réponse de l'API |


### Section `mapping`

Dans cette section, vous pouvez définir les règles de translation entre la réponse fournie par l'API et les entités Canopsis qui vont être importées.

| Paramètre            | Valeur                                                                                                                 |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| path      | Chemin à partir duquel les données seront lues (peut être vide). <br>Par exemple, "data.nested" si la réponse est de la forme  `{"data": {"nested": [ {...}, {...} ]}}` | 
| is_map    | Type de structure renvoyée par l'API. <br>Si la réponse de l'API est un tableau `[ {...}, {...} ]`, positionner `is_map: false`.<br> Si la réponse de l'API est une map `{"key1": {...}, "key2": {...}}`, positionner `is_map: true`. | 
| component | Spécifie les paramètres propres aux imports de composants  |
| resource  | Spécifie les paramètres propres aux imports de ressources  |
| service   | Spécifie les paramètres propres aux imports de services  |


### Importer des entités

#### Composants

A compléter

#### Ressources

A compléter

#### Serices

A compléter

### Importer des informations complémentaires

Il existe 4 moyens d'importer des informations complémentaires aux API

1. `set` : depuis une constante
2. `copy` : en recopiant une réponse de l'API
3. `template` : en appliquant une template (pouvant utiliser une réponse de l'API)
4. `regexp` : en exécutant un groupe de capture via une expression régulière


#### Action **set**

```yaml
infos:
  constant:
    value:
      type: set
      value: my-constant
    description: set constant
```

Resultat:

![Constant](./img/driver-api-constant.png)

#### Action **copy**

```yaml
infos:
  copy:
    value:
      type: copy
      value: env.codeEnv
    description: copy field
```

Résultat:

![Copy](./img/driver-api-copy.png)

#### Action **template**

Pour exécuter un template, vous devez spécifier le nom du champ de la réponse d'API dans `field` et le template lui-même dans `value`.  
La valeur du champ `field` est accessible grâce à la var"iable `.Field`.

```yaml
infos:
  template:
    value:
      type: template
      field: env.codeEnv
      value: 'my template value {{ .Field }}'
    description: execute template
```

Résultat:

![Template](./img/driver-api-template.png)

Il est également possible d'utiliser des fonctions dans les templates 

```yaml
infos:
  template:
    value:
      type: template
      field: env.codeEnv
      value: 'my template value {{ lowercase .Field }}'
    description: execute template with function
```

Résultat:

![Template-function](./img/driver-api-template-function.png)

#### Action **regexp**

Vous pouvez appliquer une expression régulière et utiliser des groupes de capture dans les valeurs finales.  
Pour cela, vous devez définir l'expression regulière dans `regexp` et utiliser le résultat grâce à la variable `.RegexMatch`.

```yaml
infos:
  regexp:
    value:
      type: template
      field: test_regexp
      regexp: CMDB:(?P<SI_CMDB>.*?)($|,)
      value: "regexp value {{ .RegexMatch.SI_CMDB }}"
    description: regexp example
```

Résultat:

![Regexp](./img/driver-api-regexp.png)


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
docker.canopsis.net/docker/pro/import-context-graph:'[23.04.1]'
```
Retour :
``` shell
2022-02-18T14:10:15Z INF git.canopsis.net/canopsis/canopsis-pro/pro/go-engines-pro/cmd/import-context-graph/main.go:78 > import finished deleted=0 exec_time=16.775252ms updated=3
```

### Résultats dans Canopsis :

![](./img/imported_entity.png)

![](./img/linked_alarm.png)

![](./img/var_alarm.png)

