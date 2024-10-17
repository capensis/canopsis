# Guide de développement Canopsis

Vous trouverez ici toute la documentation nécessaire au développement autour de Canopsis.

## Documentation des API

Canopsis repose sur un ensemble d'API REST ([voir une définition](https://www.redhat.com/fr/topics/api/what-is-a-rest-api)), pour son fonctionnement interne et pour son interfaçage avec d'autres programmes.

Ces API ont connu 4 versions différentes. L'APIv4 est la version actuelle.

### APIv4

La documentation des APIv4 `edition community` est disponible [par le biais de Swagger-community](./swagger-community/index.md).  
La documentation des APIv4 `edition pro` est disponible [par le biais de Swagger-pro](./swagger-pro/index.md).

Ces nouvelles API suivent l'[OpenAPI Specification 2.0](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md).

!!! important
    Les API REST imposent aux clients de respecter la spécification HTTP, et notamment de prendre en charge toute [redirection HTTP](https://www.rfc-editor.org/rfc/rfc7231#section-6.4) qui serait envoyée par l'API.

    Une erreur courante, par exemple avec l'outil `curl` (est sa bibliothèque), est que certains clients HTTP n'activent pas la prise en charge des redirections HTTP par défaut, ce qui est incorrect. Veillez donc à toujours utiliser les options de type `curl -L` en ligne de commande (ou l'option `CURLOPT_FOLLOWLOCATION` dans sa bibliothèque) afin de vous interfacer correctement avec l'ensemble des API REST.

!!! note "Note"
    Il est possible d'activer la documentation swagger de la version installée en ajoutant l'option `-docs` aux options de démarrage du binaire `canopsis-api`. 
    Une url `/swagger/index.html` sera ainsi accessible et disponible derrière le port d'écoute de l'api ( par défaut TCP/8082 )


### URL de l'API

L'API Canopsis peut être interrogée sur plusieurs URL différentes :

 - `http(s)://<canopsis>/api/v4/`: via le reverse-proxy Nginx (avec les [en-têtes CORS](https://developer.mozilla.org/fr/docs/Web/HTTP/CORS))
 - `http://<canopsis>:8082/api/v4/`: moteur `canopsis-api` directement (sans les en-têtes CORS)


!!! Warning
    Depuis Canopsis 4.4, l'URL `http(s)://<canopsis>/backend/api/v4/` est dépréciée. Pensez à mettre à jour vos scripts et applications clients de celle-ci.

> Remplacer `<canopsis>` par l'adresse IP ou le nom de domaine du Canopsis déployé.

Si l'API est interrogée via un navigateur (Firefox, Chrome, Safari, etc) ou un framework emulant un navigateur (Angular, Electron, etc) et pour lesquels en-têtes CORS sont nécessaires, alors il faut utiliser l'URL du reverse-proxy Nginx.

L'usage de l'URL du moteur `canopsis-api` est possible pour des requêtes dites "classiques", par exemple via des scripts, via l'outil `curl` ou encore via des webhooks de solutions externes.


### Authentification

#### Authentification basique HTTP

L'API Canopsis supporte l'authentification [basique HTTP](https://fr.wikipedia.org/wiki/Authentification_HTTP).  
Exemple : 

```
curl -X GET -u "user:mdp" -H "Content-Type: application/json" "http://localhost:8082/api/v4/cat/healthcheck/status"
```

#### Authentification par token

Il est possible d'utiliser un jeton (token) d'authentification pour manipuler l'API.  
Ce jeton s'obtient grâce à une authentification sur la route `/api/v4/login`.  
Le jeton est ensuite passé dans un entête `Authorization:Bearer {{ token }}`.  

Exemple :

Obtention du jeton

```
curl -X POST -H "Content-Type: application/json" -d '{"username" : "mon_user", "password" : "mon_pass"}' "http://localhost:8082/api/v4/login"
{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiJjYW5vcHNpcyIsImlzcyI6ImNhbm9wc2lzIiwiZXhwIjoxNjc2MTk2Nzk5LCJpYXQiOjE2NzM1MTgzOTksImp0aSI6ImFiNWM3MTllLWY2YzktNGI1ZS05MTQ4LWZkZWRjN2NkMjM3YiJ9.5tV7L9K-sBswjyXjO8aOQFMwor_LqoGacvg5RuBF51A"}
```

Utilisation du jeton

```
curl -X GET -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiJjYW5vcHNpcyIsImlzcyI6ImNhbm9wc2lzIiwiZXhwIjoxNjc2MTk2Nzk5LCJpYXQiOjE2NzM1MTgzOTksImp0aSI6ImFiNWM3MTllLWY2YzktNGI1ZS05MTQ4LWZkZWRjN2NkMjM3YiJ9.5tV7L9K-sBswjyXjO8aOQFMwor_LqoGacvg5RuBF51A" "http://localhost:8082/api/v4/cat/healthcheck/status"
```

!!! note "Note"

    Le jeton peut être soumis à un délai d'expiration.  
    RDV dans le profil de l'utilisateur dont vous souhaitez utiliser le jeton.  
    ![token1](./img/token1.png)

#### Clé d'authentification / Authkey

Canopsis permet une authentification à l'API par `authkey`.
Cette `authkey`, disponible pour tout utilisateur enregistré dans Canopsis, est utilisable directement dans l'URL d'API appelée ou
en tant qu'entête d'une requête vers l'API.

!!! note "Note"

    Vous pouvez récupérer la clé d'authentification d'un utilisateur depuis l'interface graphique.
    ![authkey1](./img/authkey1.png)
    ![authkey1](./img/authkey2.png)

**Dans l'URL / querystring**

Dans ce cas, il vous suffit d'ajouter le paramètre `authkey` dans l'URL.
Exemple avec la commande cURL :  

```
curl -X GET -H "Content-Type: application/json" "http://localhost:8082/api/v4/cat/healthcheck/status?authkey=5ba7537d-132d-49e1-8581-8c057de2b558"
```

**Entête (headers)**

La clé d'authentification peut également être passée en tant qu'entête de la requête HTTP à l'API.
Le nom de l'entête est dans ce cas `x-canopsis-authkey`.
Exemple avec la commande cURL :

```
curl -X GET -H "Content-Type: application/json" -H "x-canopsis-authkey: 5ba7537d-132d-49e1-8581-8c057de2b558" "http://localhost:8082/api/v4/cat/healthcheck/status"
```

## Filtres

* [Langage utilisé par les filtres](filtres/index.md)

## Structure des évènements

* [Structure des évènements](structures/index.md)

## Aides au développement

* [:warning: Obsolète :warning: Développement d'un linkbuilder](linkbuilder/index.md)

## Schémas 

* [Interactions entre les différents composants de Canopsis](./schemas/all-engines.md)
