# Sources de la documentation Canopsis

Ce répertoire contient les sources servant à construire le documentation finale de Canopsis. Elle n'a pas vocation à être lue directement depuis le dépôt Git.

Rendez vous sur [**doc.canopsis.net**](https://doc.canopsis.net) pour parcourir la documentation officielle.

## Maintenance de la documentation

La documentation est maintenue selon ces principes :

*  Utilisation de Markdown et de `mkdocs`.
*  Pour l'instant, la rédaction de la documentation se fait uniquement en français (quelques exceptions sont possibles uniquement pour le Guide de développement).
*  Ne décrire que l'UIv3 et les moteurs Go par défaut.
*  La documentation est unifiée, que le composant soit open-source ou non. Ajouter un bandeau lorsque la fonctionnalité décrite n'est pas disponible dans l'édition open-source.
*  La documentation est versionnée pour chaque nouvelle branche : 3.2.0, 3.3.0… Pas de mise à jour de la documentation pour les versions mineures, car le comportement de Canopsis ne doit pas changer lors d'une mise à jour mineure ([semver.org](https://semver.org)).
*  Chaque nouvelle branche (3.2.0, 3.3.0…) occasionne obligatoirement l'écriture d'un Guide de migration, rédigé par l'équipe d'intégration. (Si une nouvelle branche ne nécessite pas d'étape de migration, le préciser.)
*  Une fonctionnalité n'existe pas tant qu'elle n'est pas documentée.

## Particularités de la syntaxe `mkdocs`

### Format des niveaux de titres

Lors de la construction de plusieurs niveaux titres, il est obligatoire de laisser une ligne vide après chaque titre, et de laisser un espace entre le caractère `#` et le titre.

Ainsi, cette syntaxe n'est **pas** valide :
```md
# Titre1 invalide
##Titre2 invalide
Ceci est un texte invalide.
```

**mais** la syntaxe suivante fonctionnera :
```md
# Titre1 valide

## Titre2 valide

Ceci est un texte valide.
```

### Format des listes

La syntaxe Markdown de `mkdocs` est plus stricte que celle de Gitlab.

Pour `mkdocs`, le premier élément d'une liste doit toujours être précédé d'une ligne vide.

Ainsi, il ne faut **pas** utiliser la syntaxe suivante :
```md
Ceci est une liste INVALIDE :
*  premier ;
*  deuxième ;
*  troisième.
```

**mais** il faut utiliser cette syntaxe :
```md
Ceci est une liste valide :

*  premier ;
*  deuxième ;
*  troisième.
```

### Construction des hyperliens

La commande `mkdocs serve` est capable de signaler tout hyperlien interne qui ne serait plus valide.

Cependant, pour que cela fonctionne, les liens internes doivent toujours être construits en précisant le chemin **relatif** vers le **fichier Markdown** correspondant.

Ainsi, ce type de lien n'est pas sûr :
```md
Ceci est un [lien INVALIDE vers le moteur Che](/guide-administration/moteur/che/)
```

Il faut plutôt utiliser la syntaxe suivante :
```md
Ceci est un [lien valide vers le moteur Che](../guide-administration/moteur/che.md)
```

De la même façon, les liens internes comportant des ancres doivent suivre

## Visualisation de la documentation

### Installation des pré-requis

On a besoin d'installer une version récente de mkdocs (https://www.mkdocs.org/) avec le thème `material` (https://squidfunk.github.io/mkdocs-material/) associé :

```sh
$ pip install mkdocs
$ pip install mkdocs-material
```

Les réglages de mkdocs sont dans le fichier `mkdocs.yml` à la racine du dépôt.

### Visualisation

À la racine du dépôt :

```sh
$ mkdocs serve 
```

Ensuite, vous pouvez faire pointer votre navigateur sur l'URL donnée.

## Licence

La documentation est sous licence CC BY-SA 3.0 FR (https://creativecommons.org/licenses/by-sa/3.0/fr/).
