# Guide de développement Canopsis

Vous trouverez ici toute la documentation nécessaire au développement autour de Canopsis.

## Documentation des API

Canopsis repose sur un ensemble d'API REST ([voir une définition](https://www.redhat.com/fr/topics/api/what-is-a-rest-api)), pour son fonctionnement interne et pour son interfaçage avec d'autres programmes.

Ces API ont connu 4 versions différentes. L'APIv4 est la version actuelle.

### APIv4

La documentation des nouvelles APIv4 est disponible [par le biais de Swagger](./swagger).

Ces nouvelles API suivent l'[OpenAPI Specification 2.0](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md).

### Anciennes API

Les anciennes API v1 ou v2 ne sont plus utilisées par les dernières versions de Canopsis, et n'ont donc plus lieu d'être utilisées ou documentées.

## Collections de base de données

* [Collection `default_entities` pour les entités](base-de-donnees/default-entities.md)
* [Collection `periodical_alarm` pour les alarmes](base-de-donnees/periodical-alarm.md)

## Aides au développement

* [Développement d'un linkbuilder](dev-linkbuilder.md)
