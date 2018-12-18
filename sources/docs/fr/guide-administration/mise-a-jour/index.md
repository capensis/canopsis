# Procédure de mise à jour de Canopsis

## Rappels sur les numéros de version

Depuis la version 3.5.0, Canopsis veille à suivre le principe du *versionnage sémantique*.

Pour une mise à jour vers une version `X.Y.Z` de Canopsis :

*  une incrémentation du `X` désigne une mise à jour **majeure** de Canopsis.
    *  Exemple : Canopsis 2.6.8 vers 3.1.0.
    *  Cela signifie que des changements fondamentaux ont été apportés à Canopsis, ou qu'une incompatibilité majeure est survenue. La mise à jour n'est pas garantie, et il faudra généralement prévoir une réinstallation et une reconfiguration complète.
*  une incrémentation du `Y` désigne une mise à jour **intermédiaire** de Canopsis (on parle aussi d'une nouvelle *branche*).
    *  Exemple : Canopsis 3.4.0 vers 3.5.0.
    *  Elle peut apporter de nouvelles fonctionnalités, ou des changements de comportement. Elle nécessite de suivre la [mise à jour standard](#procedure-standard-de-mise-a-jour) et de vérifier la présence d'un [guide de migration](#guides-de-migration).
*  une incrémentation du `Z` apporte uniquement des **correctifs**. Il n'y a aucun changement fonctionnel.
    *  Exemple : Canopsis 3.5.0 vers 3.5.1.
    *  Pour cette raison, une simple [mise à jour standard](#procedure-standard-de-mise-a-jour) est suffisante.

## Procédure standard de mise à jour

Cette procédure s'applique aux mises à jour intermédiaires (3.4.0 vers 3.5.0) et aux mises à jour de correctifs (3.5.0 vers 3.5.1).

> **TODO :** cette procédure est en cours de rédaction.

## Guides de migration

Cette procédure s'applique au mises à jour intermédiaires (3.4.0 vers 3.5.0). Ce sont des manipulations supplémentaires et **obligatoires**.

Celles-ci sont décrites dans les documents suivants, branche par branche :

*  [3.5.0](../../notes-de-version/3.5.0.md)
*  [3.6.0](../../notes-de-version/3.6.0.md)
