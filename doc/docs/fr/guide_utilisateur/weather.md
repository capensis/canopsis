# Météo de service - weather

## Filtrage

Vos météo de service sont configurables : vous pouvez choisir quels informations afficher afin de donner un sens à une météo : n’afficher que les alertes en cours, que les éléments en pause etc…

Vous avez à disposition trois filtres spéciaux :

 * `active_pb_some` : un `booléen` valant `true`/`false`, `True`/`False`, ou `1`/`0`, indiquant si une tuile (un watcher) doit contenir au moins un `pbehavior` actif pour être affiché.
 * `active_pb_all` : même chose que pour `active_pb_some` sauf qu’il faut que **toutes** les entités liées aient un `pbehavior` actif pour être affiché.
 * `active_pb_type` : ce filtre peut être utilisé plusieurs fois et va vous permettre de demander d’afficher les tuiles contenant des `pbehavior` d’un certain `type`. Par exemple si vous voulez filtrer sur les tuiles ayant des `pbehavior` en type `pause` ou `maintenance` ou les deux. Dans le cas où une tuile dispose de `pbehavior` d’autre type mais qu’au moins un d’entre eux est du type indiqué par ce filtre, la tuile sera affichée.

Ces options ne sont disponibles qu’au travers du filtre créé avec le `query builder` :

Voici un exemple de filtre créé sur une vue MdS :

![weather_filter](weather_filter.png)

Et le filtre généré correspondant (pour vérification) :

![weather_generated_filter](weather_generated_filter.png)

### Comportement des filtres

Les trois options mises ensemble se comportent comme un `ET` logique. Exemple :

 * `active_pb_some: True`
 * `active_pb_type: pause`

Cette règle peut se lire de la façon suivante : *afficher les tuiles ayant au moins un pbehavior actif et dont l’un au moins d’entre eux est de type pause*.

Autre exemple :

 * `active_pb_all: False`
 * `active_pb_some: True`
 * `active_pb_type: maintenance`

*afficher les tuiles ayant au moins une entité sans pbehavior, avec avec au moins un pbehavior, et de type maintenance au moins.*

Autre exemple :

 * `active_pb_all: True`

*afficher les tuiles dont toutes les entités ont un pbehavior actif, peu importe leur type.*

Autre exemple :

 * `active_pb_type: maintenance`
 * `active_pb_type: pause`

*s’il y a des pbehavior actifs, ils doivent être de type pause ou maintenance.*