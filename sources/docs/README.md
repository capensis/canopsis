# Sources de la documentation Canopsis

Ce répertoire contient les sources servant à construire le documentation finale de Canopsis. Elle n'a pas vocation à être lue directement depuis le dépôt Git.

Rendez vous sur [**doc.canopsis.net**](https://doc.canopsis.net) pour parcourir la documentation officielle.

## Maintenance de la documentation

La documentation est maintenue selon ces principes :

*  Pour l'instant, la rédaction de la documentation se fait uniquement en français (quelques exceptions sont possibles uniquement pour le Guide de développement).
*  Ne décrire que l'UIv3 et les moteurs Go par défaut.
*  La documentation est unifiée, que le composant soit open-source ou non. Ajouter un bandeau lorsque la fonctionnalité décrite n'est pas disponible dans l'édition open-source.
*  La documentation est versionnée pour chaque nouvelle branche : 3.2.0, 3.3.0… Pas de mise à jour de la documentation pour les versions mineures, car le comportement de Canopsis ne doit pas changer lors d'une mise à jour mineure ([semver.org](https://semver.org)).
*  Chaque nouvelle branche (3.2.0, 3.3.0…) occasionne obligatoirement l'écriture d'un Guide de migration, validé par l'équipe d'intégration. (Si une nouvelle branche ne nécessite pas d'étape de migration, le préciser.)
*  Une fonctionnalité n'existe pas tant qu'elle n'est pas documentée.

## Licence

La documentation est sous licence CC BY-SA 3.0 FR (https://creativecommons.org/licenses/by-sa/3.0/fr/).
