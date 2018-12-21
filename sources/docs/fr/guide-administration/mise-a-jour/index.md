# Procédure de mise à jour de Canopsis

## Rappels sur les numéros de version

Depuis la version 3.5.0, Canopsis veille à suivre le principe du *versionnage sémantique*.

Sur une version `X.Y.Z` de Canopsis :

*  le `X` décrit une mise à jour **majeure** de Canopsis.
    *  Cela signifie que des changements fondamentaux ont été apportés à Canopsis, ou qu'une incompatibilité majeure est survenue (exemple : mise à jour d'un Canopsis 2.6.8 vers un Canopsis 3.1.0). La mise à jour n'est pas garantie, et il faudra parfois prévoir une réinstallation suivie d'une migration.
*  le `Y` décrit une mise à jour **intermédiaire** de Canopsis (on parle aussi d'une nouvelle *branche*).
    *  Elle peut apporter de nouvelles fonctionnalités, ou des changements de comportement. Ces mises à jour peuvent nécessiter d'appliquer une [procédure de mise à jour](#guides-de-migration), comme décrit plus bas (exemple : mise à jour d'une version 3.4.0 vers 3.5.0).
*  le `Z` apporte uniquement des **correctifs**, et aucun changement fonctionnel.
    *  Pour cette raison, une simple [mise à jour standard](#procedure-standard-de-mise-a-jour) est suffisante (exemple : mise à jour d'une version 3.5.0 vers 3.5.1).

## Procédure standard de mise à jour

Cette procédure s'applique aux mises à jour intermédiaires (3.4.0 vers 3.5.0) et aux mises à jour de correctifs (3.5.0 vers 3.5.1).

> **TODO :** cette procédure est en cours de rédaction.

## Guides de migration

La mise à jour vers certaines branches de Canopsis nécessite parfois des manipulations supplémentaires et **OBLIGATOIRES**.

Celles-ci sont décrites dans les documents suivants, branche par branche :

*  [3.5.0](../../notes-de-version/3.5.0.md)
*  [3.6.0](../../notes-de-version/3.6.0.md)
*  [3.7.0](../../notes-de-version/3.7.0.md)
