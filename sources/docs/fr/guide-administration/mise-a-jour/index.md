# Procédure de mise à jour de Canopsis

!!! note
    Cette procédure ne décrit que la mise à jour d'une instance mono-nœud de Canopsis open-core.

## Principes des numéros de version de Canopsis

Une nouvelle version de Canopsis sort généralement tous les 15 jours.

Afin de décrire ses changements, Canopsis suit le principe du *versionnage sémantique*.

### Mises à jour majeures

On parle d'une mise à jour **majeure** lorque le premier chiffre change ; par exemple, lors du passage de Canopsis 2 à Canopsis 3.

Ceci signifie que des changements fondamentaux ont été apportés à Canopsis (changement d'interface, de paradigmes, de technologies…), ou qu'une incompatibilité majeure est survenue. Ces changements sont exceptionnels.

La mise à jour d'une installation existante n'est **pas garantie** lors d'une mise à jour majeure. Il faudra généralement prévoir une réinstallation et une reconfiguration complètes.

### Mises à jour intermédiaires

Les mises à jour **intermédiaires** sont les plus courantes. Il s'agit, par exemple, d'une mise à jour de Canopsis 3.25.0 vers Canopsis 3.26.0.

Ce type de mise à jour apporte de nouvelles fonctionnalités, souvent des corrections de bugs et parfois quelques changements de comportement.

Elle nécessite de suivre la procédure de mise à jour decrite ci-dessous.

### Mises à jour mineures

La mise à jour d'un Canopsis 3.26.0 vers 3.26.1 constitue une mise à jour **mineure**.

Ces mises à jour ne comportent uniquement que des correctifs. Il n'y a aucun ajout de fonctionnalités ou de changement fonctionnel.

Elle nécessite de suivre la procédure de mise à jour décrite ci-dessous.

## Procédure de mise à jour

Cette procédure s'applique aux mises à jour intermédiaires (3.25.0 vers 3.26.0) et aux mises à jour mineures (3.26.0 vers 3.26.1).

Vous devez *impérativement* lire **chacune** des [notes de version](../../index.md) vous séparant de votre version précédente à votre version cible, avant de procéder à une mise à jour.

### Mise à jour en installation par paquets

Veuillez tout d'abord obligatoirement lire les [notes de version](../../index.md) qui vous concernent **avant** de démarrer toute manipulation. Des prérequis cruciaux peuvent y être mentionnés.

Puisque la [procédure d'installation par paquets](../installation/installation-paquets.md) ajoute un dépôt Canopsis dans votre gestionnaire de paquets, celui-ci vous permettra de mettre à jour Canopsis avec le reste de votre système.

!!! attention
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

Il faut ensuite lancer le script `canopsinit` pour appliquer les éventuelles procédures automatisées de mise à jour.

!!! attention
    Avant Canopsis 3.27.0, la procédure suivante réinitialise les identifiants `root` de la base utilisateur, interne à Canopsis, ainsi que sa authkey associée.

Si vous utilisez la configuration « moteurs Python uniquement » et « édition Canopsis open-core » (qui sont les réglages par défaut), lancer :
```sh
su - canopsis -c "canopsinit"
```

En revanche, si vous utilisez une configuration « moteurs Go » et « édition Canopsis CAT », lancer :
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
canoctl restart
```

Ne pas oublier d'appliquer toute éventuelle procédure supplémentaire décrite dans chacune des [notes de version](../../index.md) qui vous concernent.

**Si vous bénéficiez d'un développement spécifique** (modules ou add-ons ayant été spécifiquement développés pour votre installation), assurez-vous de suivre toute procédure complémentaire vous ayant été communiquée.

Vous pouvez alors vous connecter à nouveau sur l'interface Canopsis pour valider que tout fonctionne correctement.

### Mise à jour en environnement Docker

!!! todo
    Cette procédure est en cours de rédaction.
