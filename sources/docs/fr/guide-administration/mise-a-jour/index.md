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

!!! note
    Cette procédure ne décrit que la mise à jour d'une instance **mono-nœud** de Canopsis.

Cette procédure s'applique aux mises à jour intermédiaires (3.4.0 vers 3.5.0) et aux mises à jour de correctifs (3.5.0 vers 3.5.1).

Vous devez **impérativement** lire chacune des [notes de version](../../notes-de-version/index.md) vous séparant de votre version précédente à votre version cible, avant de procéder à une mise à jour.

### Mise à jour en installation par paquets

Puisque la [procédure d'installation](../installation/index.md) ajoute un dépôt Canopsis dans votre gestionnaire de paquets, il suffit de lancer une mise à jour pour bénéficier d'un paquet à jour.

**Cette mise à jour causera une interruption de service.**

Sur un environnement Debian :
```sh
apt update
apt upgrade
```

Sur un environnement CentOS :
```sh
yum update
```

!!! attention
    La procédure suivante réinitialise les identifiants `root` de la base utilisateur, interne à Canopsis, ainsi que sa authkey associée ([Bug #1431](https://git.canopsis.net/canopsis/canopsis/issues/1431)).

Il faut ensuite lancer le script `canopsinit` (en tant qu'utilisateur Unix `canopsis`) pour appliquer les éventuelles procédures automatisées de mise à jour.

Si vous utilisez la configuration « moteurs Python uniquement » et « édition Canopsis Core » (qui sont les réglages par défaut), lancer :
```sh
su - canopsis -c "canopsinit"
```

**En revanche**, si vous installez une version de Canopsis supérieure ou égale à 3.17.0, et que vous utilisez une configuration « moteurs Go » et « édition Canopsis CAT », lancer :
```sh
# Valeurs acceptées : --canopsis-edition core OU cat, --canopsis-stack python OU go.
# "core" et "python" sont les valeurs par défaut.
su - canopsis -c "canopsinit --canopsis-edition cat --canopsis-stack go"
```

S'assurer que toute modification des unités systemd soit bien prise en compte :
```sh
systemctl daemon-reload
```

Puis, redémarrer l'ensemble des moteurs Canopsis :
```sh
/opt/canopsis/bin/canopsis-systemd restart
```

Appliquer toute éventuelle procédure supplémentaire décrite dans les [notes de version](../../notes-de-version/index.md) qui vous concernent.

Vous pouvez alors vous connecter à nouveau sur l'interface Canopsis pour valider que tout fonctionne correctement.

### Mise à jour en environnement Docker

!!! todo
    Cette procédure est en cours de rédaction.

## Guides de migration

Cette procédure s'applique au mises à jour intermédiaires (3.4.0 vers 3.5.0). Ce sont des manipulations supplémentaires et **obligatoires**.

Celles-ci sont décrites dans les documents suivants, branche par branche :

*  [3.5.0](../../notes-de-version/3.5.0.md)
*  [3.6.0](../../notes-de-version/3.6.0.md)
*  [3.7.0](../../notes-de-version/3.7.0.md)
*  [3.8.0](../../notes-de-version/3.8.0.md)
*  [3.9.0](../../notes-de-version/3.9.0.md)
*  [3.10.0](../../notes-de-version/3.10.0.md)
*  [3.11.0](../../notes-de-version/3.11.0.md)
*  [3.12.0](../../notes-de-version/3.12.0.md)
*  [3.13.0](../../notes-de-version/3.13.0.md), [3.13.1](../../notes-de-version/3.13.1.md), [3.13.2](../../notes-de-version/3.13.2.md)
*  [3.14.0](../../notes-de-version/3.14.0.md)
*  [3.15.0](../../notes-de-version/3.15.0.md)
*  [3.16.0](../../notes-de-version/3.16.0.md)
*  [3.17.0](../../notes-de-version/3.17.0.md)
*  [3.18.0](../../notes-de-version/3.18.0.md), [3.18.1](../../notes-de-version/3.18.1.md)
*  [3.19.0](../../notes-de-version/3.19.0.md)
*  [3.20.0](../../notes-de-version/3.20.0.md)
*  [3.21.0](../../notes-de-version/3.21.0.md), [3.21.1](../../notes-de-version/3.21.1.md)
