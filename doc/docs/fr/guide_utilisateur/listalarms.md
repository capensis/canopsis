# Bac à alarmes

## Recherche

Il est possible de filtrer les résultats du bac à alarmes suivant deux méthodes (dans le même champs de recherche):

 - Soit en saisissant simplement un mot-clef, qui sera recherché dans l'ensemble des champs affichés dans le bac à alarmes ;
 - Soit en précisant explicitement les champs dans lequel faire le filtrage, au moyen d'une grammaire (voir ci-dessous).

### Recherche avec grammaire

Pour activer ce mode, il faut commencer la recherche par un tiret, puis l'on précise le nom d'un champs, suivi d'un égal et de la valeur désirée (ne pas oublier les " pour les champs textuels). Il est possible de composer plusieurs filtres ensemble avec des opérateurs (AND, OR). Voici quelques exemples de filtres valides :

```
- Connector = "connector_1"
- Connector="connector_1"
      - Connector = "connector_1"
- Connector="connector_1" AND Resource="resource_3"
- Connector="connector_1" OR Resource="resource_3"
```
