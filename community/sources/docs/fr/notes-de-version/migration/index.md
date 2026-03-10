# Guide de migration vers Canopsis 25.10.0

Ce guide donne les instructions vous permettant de mettre à jour Canopsis 25.04 (dernière version disponible) vers [la version 25.10.0](../25.10.0.md).

## Prérequis

L'ensemble de cette procédure doit être lue avant son exécution.

Ce document ne prend en compte que Canopsis Community et Canopsis Pro : tout développement personnalisé dont vous pourriez bénéficier ne fait pas partie du cadre de ce Guide de migration.

Les fichiers de référence qui sont mentionnés dans ce guide sont disponibles à ces adresses

| Édition           | Sources                                                                                                                              |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Édition Community | [https://git.canopsis.net/canopsis/canopsis-community/-/releases](https://git.canopsis.net/canopsis/canopsis-community/-/releases)   |
| Édition pro       | [https://git.canopsis.net/sources/canopsis-pro-sources/-/releases](https://git.canopsis.net/sources/canopsis-pro-sources/-/releases) |

## Réalisation d'une sauvegarde

Des sauvegardes sont toujours recommandées, qu'elles soient régulières ou lors de modifications importantes.

## Procédure de mise à jour

<div class="grid cards" markdown>

-   :material-docker:{ .lg .middle } __Docker-Compose__

    ---

    Réaliser la migration de votre Canopsis dans un environnement `Docker-Compose`

    [:octicons-arrow-right-24: Guide de migration](docker.md)

-   :material-package-variant:{ .lg .middle } __Paquets RPM__

    ---

    Réaliser la migration de votre Canopsis dans un environnement `RPM`

    [:octicons-arrow-right-24: Guide de migration](./rpm.md)

-   :material-ship-wheel:{ .lg .middle } __Helm__

    ---

    Réaliser la migration de votre Canopsis dans un environnement `Kubernetes/Helm`

    [:octicons-arrow-right-24: Guide de migration](./helm.md)

</div>

