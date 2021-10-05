# Guide de développement Canopsis

Vous trouverez ici toute la documentation nécessaire au développement autour de Canopsis.

## Documentation des API

Canopsis repose sur un ensemble d'API REST ([voir une définition](https://www.redhat.com/fr/topics/api/what-is-a-rest-api)), pour son fonctionnement interne et pour son interfaçage avec d'autres programmes.

Ces API ont connu 4 versions différentes : l'APIv4 est la version actuelle, mais d'anciennes APIv1 et APIv2 sont encore utilisées à ce jour. L'APIv3 a été totalement abandonnée.

### Nouvelles APIv4

La documentation des nouvelles APIv4 est disponible [avec Swagger](./swagger).

Ces nouvelles API suivent l'[OpenAPI Specification 2.0](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/2.0.md).

### Anciennes APIv1 et APIv2

* [`/api/contextgraph/import`](api/api-v2-import.md)
* [`/api/internal`](api/api-internal.md)
* [`/api/v2/broadcast-message`](api/api-v2-broadcast-message.md)
* [`/api/v2/dynamic-infos`](api/api-v2-dynamic-infos.md)
* [`/api/v2/event`](api/api-v2-event.md) et [structure d'un évènement](struct-event.md)
* [`/api/v2/eventfilter`](api/api-v2-event-filter.md)
* [`/api/v2/healthcheck`](api/api-v2-healthcheck.md)
* [`/api/v2/metaalarmrule`](api/api-v2-meta-alarm-rule.md)
* [`/api/v2/watcherng`](api/api-v2-watcherng.md)
* [`/api/v2/weather`](api/api-v2-weather.md)

## Collections de base de données

* [Collection `default_entities` pour les entités](base-de-donnees/default-entities.md)
* [Collection `periodical_alarm` pour les alarmes](base-de-donnees/periodical-alarm.md)

## Guides de développement

* [Développement d'un linkbuilder](dev-linkbuilder.md)
