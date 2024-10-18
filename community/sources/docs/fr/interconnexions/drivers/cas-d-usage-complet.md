# Cas d'usage du driver API (`import-context-graph`)

La documentation de référence du programme `import-context-graph` est [accessible par ce lien](driver-api.md).

Ce cas d'usage décrit une organisation de services et de composants répartis sur 3 niveaux.

Le schéma ci-après présente ces niveaux :

```mermaid
classDiagram
    Applications <|-- Composants
    Composants <|-- Serveurs
        
    class Serveurs{
      type: component
      name : codeServeur
      category: Serveur
      infos.codeComposant = codeComposant
      infos.type = type
    }

    class Composants{
      type: service
      name: codeComposant
      category: Composant
      infos.applications: applications ['app1', 'app2']
      pattern(): infos.codeComposant eq code
    }

    class Applications{
      type: service
      name: codeApplication
      category: Application
      infos.libelle: libelle
      pattern(): infos.applications has_one_of codeApplication

    }
```

!!! note "Information"

    Les `composants` décrits dans cette organisation n'ont pas de lien avec les entités de type `composant` de Canopsis.
    Ils représentent ici un ensemble de `Serveurs applicatifs`.


Les imports sont lancés via une ligne de commande qu'il vous faut adapter à votre configuration.  
Voici un exemple pour les serveurs d'applications.

```sh
docker run \
-e EXTERNAL_API_USERNAME='[testuser]' \
-e EXTERNAL_API_PASSWORD='[testpassword]' \
-e CPS_MONGO_URL='mongodb://cpsmongo:canopsis@mongodb:27017/canopsis' \
-e CPS_AMQP_URL='amqp://cpsrabbit:canopsis@rabbitmq:5672/canopsis' \
-e CPS_POSTGRES_URL='postgresql://cpspostgres:canopsis@timescaledb:5432/canopsis' \
--network=canopsis-pro_default -it --rm -v \
"${PWD}/import-serveursapplications.yml:/opt/canopsis/share/config/import-context-graph/api.yml" \
docker.canopsis.net/docker/pro/import-context-graph:'23.04.1'
```

## Import des Serveurs d'applications

### API

??? note "Réponse de l'API"

	```json
	{
	  "took": 1724,
	  "timed_out": false,
	  "_shards": {
	    "total": 1,
	    "successful": 1,
	    "skipped": 0,
	    "failed": 0
	  },
	  "hits": {
	    "total": {
	      "value": 3,
	      "relation": "eq"
	    },
	    "hits": [
	      {
		"_id": "Doc-ServeurAPP1",
		"_source": {
		  "id": "Doc-ServeurAPP1",
		  "desc": "Serveur Application1",
		  "type": "WEBSERVICE",
		  "env" : "Production",
		  "codeComposant" : "Doc-Composant1"
		}
	      },
	      {
		"_id": "Doc-ServeurAPP2",
		"_source": {
		  "id": "Doc-ServeurAPP2",
		  "desc": "Serveur Application2",
		  "type": "WEBSERVICE",
		  "env" : "Production",
		  "codeComposant" : "Doc-Composant1"
		}
	      },
	      {
		"_id": "Doc-ServeurAPP3",
		"_source": {
		  "id": "Doc-ServeurAPP3",
		  "desc": "Serveur Application3",
		  "type": "WEBSERVICE",
		  "env" : "Production",
		  "codeComposant" : "Doc-Composant2"
		}
	      }
	    ]
	  }
	}
	```

### Mapping

La configuration suivante de `mapping` est utilisée :

```yaml
mapping:
  path: hits.hits
  is_map: false
  component:
    name: 
      type: copy
      value: _id
    description: 
      type: template
      field: _source.desc
      value: '{{uppercase .Field}}'
    category:
      type: set
      value: ServeurApplication
    infos:
      codeComposant:
        value:
          type: template
          field: _source.codeComposant
          value: '{{uppercase .Field}}'
        description: Code composant
      env:
        value:
          type: copy
          value: _source.env
        description: Environnement
```

### Résultats

Une fois l'import exécuté, vous retrouvez les Serveurs d'applications dans l'explorateur de contexte de Canopsis.  

![SA1](./img/casdusage-sa1.png)

## Import des Composants

### API

??? note "Réponse de l'API"

	```json
	{
	  "took": 43,
	  "timed_out": false,
	  "_shards": {
	    "total": 1,
	    "successful": 1,
	    "skipped": 0,
	    "failed": 0
	  },
	  "hits": {
	    "total": {
	      "value": 3,
	      "relation": "eq"
	    },
	    "hits": [
	      {
		"_id": "Doc-Composant1",
		"_source": {
		  "codeComposant": "Doc-Composant1",
		  "env": "Production",
		  "desc" : "Composant 1",
		  "applications": [
		    "Doc-APP1"
		  ]
		}
	      },
	      {
		"_id": "Doc-Composant2",
		"_source": {
		  "codeComposant": "Doc-Composant2",
		  "env": "Production",
		  "desc" : "Composant 2",
		  "applications": [
		    "Doc-APP1"
		  ]
		}
	      },
	      {
		"_id": "Doc-Composant3",
		"_source": {
		  "codeComposant": "Doc-Composant3",
		  "env": "Production",
		  "desc" : "Composant 3",
		  "applications": [
		    "Doc-APP2"
		  ]
		}
	      }
	    ]
	  }
	}
	```

### Mapping

La configuration suivante de `mapping` est utilisée :

```yaml
mapping:
  path: hits.hits
  is_map: false
  service:
    name: 
      type: template
      field: _id
      value: '{{uppercase .Field}}'
    description: 
      type: template
      field: _source.desc
      value: '{{uppercase .Field}}'
    category:
      type: set
      value: Composant
    pattern:
      cond: eq
      field: infos.codeComposant
      field_type: string
      value:
        type: template
        field: _source.codeComposant
        value: '{{uppercase .Field}}'
    infos:
      env:
        value:
          type: copy
          value: _source.env
        description: Environnement
      applications:
        value:
          type: copy
          value: _source.applications
        description: Liste des applications
```

### Résultats

Une fois l'import exécuté, vous retrouvez les Composants dans l'explorateur de contexte de Canopsis.  

![CO1](./img/casdusage-co1.png)

## Import des Applications

### API

??? note "Réponse de l'API"

	```json
	{
	  "took": 23,
	  "timed_out": false,
	  "_shards": {
	    "total": 1,
	    "successful": 1,
	    "skipped": 0,
	    "failed": 0
	  },
	  "hits": {
	    "total": {
	      "value": 3,
	      "relation": "eq"
	    },
	    "hits": [
	      {
		"_id": "Doc-APP1",
		"_source": {
		  "codeApplication": "Doc-APP1",
		  "env": "Production",
		  "desc" : "Application 1"
		}
	      },
	      {
		"_id": "Doc-APP2",
		"_source": {
		  "codeApplication": "Doc-APP2",
		  "env": "Production",
		  "desc" : "Application 2"
		}
	      },
	      {
		"_id": "Doc-APP3",
		"_source": {
		  "codeApplication": "Doc-APP3",
		  "env": "Production",
		  "desc" : "Application 3"
		}
	      }
	    ]
	  }
	}
	```

### Mapping

La configuration suivante de `mapping` est utilisée :

```yaml
mapping:
  path: hits.hits
  is_map: false
  service:
    name: 
      type: template
      field: _id
      value: '{{.Field}}'
    description: 
      type: template
      field: _source.desc
      value: '{{uppercase .Field}}'
    category:
      type: set
      value: Application
    pattern:
      cond: has_one_of
      field: infos.applications
      field_type: string_array
      value:
        type: copy
        value: _source.codeApplication
    infos:
      env:
        value:
          type: copy
          value: _source.env
        description: Environnement
```

### Résultats

Une fois l'import exécuté, vous retrouvez les Applications dans l'explorateur de contexte de Canopsis.  

![APP1](./img/casdusage-app1.png)

## Exploitation de cette organisation

### Arbre de dépendances

Vous pouvez visualiser l'arbre de dépendances de l'application `Doc-APP1` en vous rendant dans l'explorateur de contexte.

![Exploit1](./img/casdusage-exploit1.png)

Vous constatez alors que `Doc-APP1` dépend de `DOC-COMPOSANT1` et `DOC-COMPOSANT2` eux mêmes dépendants de `Doc-ServeurAPP1`, `Doc-ServeurAPP2`, et `Doc-ServeurAPP3`.

Vous pour également utiliser le widget `Cartographie` pour visualiser cet arbre : 

![Exploit2](./img/casdusage-exploit2.png)

### Bac à alarmes

Si une alarme est émise sur une dépendance alors vous pourrez observer dans un bac à alarmes les conséquences.

```sh
curl -H "Content-Type: application/json" -X POST -u root:root -d '{
   "connector" : "doc-connector",
   "connector_name" : "doc",
   "event_type" : "check",
   "source_type" : "component",
   "component" : "Doc-ServeurAPP1",
   "state" : 2,
   "output" : "Instance KO|cpu=80%"
 }' 'http://localhost:8082/api/v4/event'
```

![Exploit3](./img/casdusage-exploit3.png)

L'alarme émise à eu un impact sur le composant `Doc-COMPOSANT1` qui lui-même a eu un impact sur l'application `Doc-APP1`.  
L'arbre de dépendance est visible dans la colonne `TRACK SOURCE` des alarmes de services.

Par ailleurs, la cartographie est elle aussi mise à jour en conséquence

![Exploit4](./img/casdusage-exploit4.png)
