# Sources de la documentation Canopsis

Ce répertoire contient les sources servant à construire le documentation de Canopsis. Elle n'a pas vocation à être lue directement depuis le dépôt Git.

Rendez vous sur [doc.canopsis.net](https://doc.canopsis.net) pour parcourir la documentation officielle.

## Maintenance de la documentation

La documentation est maintenue selon ces principes :

* Utilisation de Markdown et de `mkdocs`, avec les extensions `mkdocs-material` et `pymdown-extensions`. Du HTML peut être utilisé avec parcimonie, par exemple pour intégrer des vidéos YouTube.
* Pour l'instant, la rédaction de la documentation se fait uniquement en français (quelques exceptions sont possibles uniquement pour le Guide de développement).
* On documente à la fois Canopsis Community et Canopsis Pro. Lorsqu'une fonctionnalité n'est disponible qu'avec l'édition Pro, le mentionner dans un bandeau explicite en haut du document.
* Lors de la publication d'une nouvelle version de Canopsis, des notes de versions doivent obligatoirement être rédigées. En dehors des mise sà jour mineures (ex : 4.5.2), un Guide de migration doit aussi être rédigé, en règle générale.
* Une fonctionnalité n'existe pas tant qu'elle n'est pas documentée.

## Particularités de la syntaxe `mkdocs`

> **Note :** les limitations suivantes ne sont peut-être plus d'actualité dans une version suffisamment à jour de mkdocs. Ajuster ce document au besoin.

### Format des niveaux de titres

Lors de l'utilisation de sections et de sous-sections, il est obligatoire de laisser une ligne vide après chaque titre, et de laisser un espace entre le caractère `#` et le titre.

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

Pour `mkdocs`, le premier élément d'une liste doit toujours être précédé d'une ligne vide.

Ainsi, il ne faut **pas** utiliser la syntaxe suivante :

```md
Ceci est une liste INVALIDE :
* premier ;
* deuxième ;
* troisième.
```

**mais** il faut utiliser cette syntaxe :

```md
Ceci est une liste valide :

* premier ;
* deuxième ;
* troisième.
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

De la même façon, les ancres doivent toujours être saisies en minuscules et en ASCII :

```md
## Base de données

Blabla.

## Autre section

Rendez-vous dans la [section Base de données](ce-fichier.md#base-de-donnees).
```

### Mise en forme avancée

Pour réaliser de la mise en forme avancée, regarder les documents existants, ainsi que la documentation officielle des extensions Markdown que nous utilisons :

* <https://squidfunk.github.io/mkdocs-material/>
* <https://facelessuser.github.io/pymdown-extensions/>

**Attention :** certaines fonctionnalités peuvent nécessiter [une mise à jour](requirements.txt), ou l'édition Insiders (que nous n'utilisons pas).

## Visualisation de la documentation

### Installation des pré-requis

On a besoin d'installer une version récente de [mkdocs](https://www.mkdocs.org/) avec le thème [`material`](https://squidfunk.github.io/mkdocs-material/) associé :

```sh
python3 -m pip install -r community/sources/docs/requirements.txt
```

Les réglages de mkdocs sont dans le fichier [`community/mkdocs.yml`](../../mkdocs.yml).

### Visualisation

En étant dans le répertoire `/community` du monorepo :

```sh
mkdocs serve
```

Ensuite, vous pouvez faire pointer votre navigateur sur l'URL donnée.

## Licence

La documentation est la propriété de Capensis et est sous licence CC BY-SA 3.0 FR (<https://creativecommons.org/licenses/by-sa/3.0/fr/>).
