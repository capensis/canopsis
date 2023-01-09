# Guide de développement Canopsis

Vous trouverez ici toute la documentation nécessaire au développement autour de Canopsis.

## Documentation des API

Canopsis repose sur un ensemble d'API REST ([voir une définition](https://www.redhat.com/fr/topics/api/what-is-a-rest-api)), pour son fonctionnement interne et pour son interfaçage avec d'autres programmes.

Ces API ont connu 4 versions différentes. L'APIv4 est la version actuelle.

### APIv4

La documentation des nouvelles APIv4 est disponible [par le biais de Swagger](./swagger).

Ces nouvelles API suivent l'[OpenAPI Specification 2.0](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md).

!!! important
    Les API REST imposent aux clients de respecter la spécification HTTP, et notamment de prendre en charge toute [redirection HTTP](https://www.rfc-editor.org/rfc/rfc7231#section-6.4) qui serait envoyée par l'API.

    Une erreur courante, par exemple avec l'outil `curl` (est sa bibliothèque), est que certains clients HTTP n'activent pas la prise en charge des redirections HTTP par défaut, ce qui est incorrect. Veillez donc à toujours utiliser les options de type `curl -L` en ligne de commande (ou l'option `CURLOPT_FOLLOWLOCATION` dans sa bibliothèque) afin de vous interfacer correctement avec l'ensemble des API REST.

### Anciennes API

Les anciennes API v1 ou v2 ne sont plus utilisées par les dernières versions de Canopsis, et n'ont donc plus lieu d'être utilisées ou documentées.

### URL de l'API

L'API Canopsis peut être interrogée sur plusieurs URL différentes :

 - `http(s)://<canopsis>/api/v4/`: via le reverse-proxy Nginx (avec les [en-têtes CORS](https://developer.mozilla.org/fr/docs/Web/HTTP/CORS))
 - `http://<canopsis>:8082/api/v4/`: moteur `canopsis-api` directement (sans les en-têtes CORS)


!!! Warning
    Depuis Canopsis 4.4, l'URL `http(s)://<canopsis>/backend/api/v4/` est dépréciée. Pensez à mettre à jour vos scripts et applications clients de celle-ci.

> Remplacer `<canopsis>` par l'adresse IP ou le nom de domaine du Canopsis déployé.

Si l'API est interrogée via un navigateur (Firefox, Chrome, Safari, etc) ou un framework emulant un navigateur (Angular, Electron, etc) et pour lesquels en-têtes CORS sont nécessaires, alors il faut utiliser l'URL du reverse-proxy Nginx.

L'usage de l'URL du moteur `canopsis-api` est possible pour des requêtes dites "classiques", par exemple via des scripts, via l'outil `curl` ou encore via des webhooks de solutions externes.

## Filtres

* [Langage utilisé par les filtres](filtres/index.md)

## Structure des évènements

* [Structure des évènements](structures/index.md)

## Collections de base de données

* [Collection `default_entities` pour les entités](base-de-donnees/default-entities.md)
* [Collection `periodical_alarm` pour les alarmes](base-de-donnees/periodical-alarm.md)

## Aides au développement

* [Développement d'un linkbuilder](linkbuilder/index.md)
