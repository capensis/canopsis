# Développement d'un linkbuilder

## Contexte

Le principe du *linkbuilder* est de « calculer » un ou plusieurs liens en
rapport avec une alarme ou l'entité concernée par l'alarme, pour présenter ces
liens sur l'interface graphique de Canopsis (bac à alarmes ou météo de
services).

Le résultat obtenu et l'utilisation de base sont illustrés dans la
documentation [linkbuilder][0] du guide d'administration.

Un *linkbuilder* nommé `BasicAlarmLinkBuilder` est fourni par défaut dans
Canopsis : celui-ci permet de produire un simple lien basé sur un modèle
d'URL configurable (une base et un ou des « trous » qui seront remplacés par
des éléments de l'alarme ou de l'entité concernée).

Pour des cas où plusieurs types de liens doivent être produits, selon une
logique métier plus complexe, il est possible de développer une classe
personnalisée de *linkbuilder*.

Dans le cadre d'une souscription Canopsis auprès de Capensis, le développement
du linkbuiler adapté au processus du demandeur peut être réalisé par Capensis.
Cependant, toute personne peut faire ce développement, avec ou sans
souscription Canopsis.

Ce guide de développement apporte des éléments utiles à la conception d'un
*linkbuilder* personnalisé.

[0]: ../../guide-administration/linkbuilder/index.md

## Considérations générales

Le *linkbuilder* entre en jeu au sein du service http `oldapi` de Canopsis,
c'est-à-dire un composant de l'environnement Python de Canopsis.

Dès que l'instance Canopsis est configurée pour utiliser un *linkbuilder*,
celui-ci sera invoqué à chaque chargement d'alarme (par exemple via un widget
bac à alarmes). On veillera donc à ce que le processus de construction des
liens soit le moins complexe possible, et ne repose pas sur l'accès à des
services externes à Canopsis.

## Développer un linkbuilder

Le développement du *linkbuilder* se fait en Python 2 (2.7).

Il s'agit de créer un module Python qui devra être placé dans le dossier
`common/link_builder/` des sources de Canopsis.

Le module contiendra la classe personnalisée, qui doit hériter de la classe
abstraite `HypertextLinkBuilder`. La classe personnalisée doit au minimum
implémenter une méthode `build()` et cette méthode doit retourner la
*liste de liens* selon une structure de données définie.

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""My specific linkbuilder class"""

from __future__ import unicode_literals

from canopsis.common.link_builder.link_builder import HypertextLinkBuilder


class MySpecificLinkBuilder(HypertextLinkBuilder):
    """Link builder for my very own project.

    Features:
    - not much! (yet)
    """

    def build(self, entity, options={}):
        # This is how a linkbuilder returns no links: empty dict
        return {}
```

### Structure de données attendue

Pour une alarme donnée, l'interface graphique de Canopsis présentera les liens
(URL) regroupés par **catégorie** et un **label** peut être ajouté pour chaque
URL.

La structure de données que le *linkbuilder* doit renvoyer est un dictionnaire.

- Les clefs sont les catégories
- Sous chaque catégorie, la valeur est une liste de liens
- Un lien est un dictionnaire avec

    * une clef `label` pour le titre du lien
    * une clef `link` pour l'URL cible

```python
{
    "cat1": [
        {
            "label": "link1",
            "link": "http://url1.example.com/page/42"
        },
        {
            "label": "link2",
            "link": "http://url2.example.com/page/43"
        }
    ],
    "cat2": [
        {
            "label": "link3",
            "link": "http://url3.example.com/"
        }
    ]
}
```

### Recettes de code

#### Récupérer l'entité en cours

L'entité est passée à la méthode `build()` (argument `entity`).

Ci-dessous un exemple où l'on va construire un lien de manière différente selon
le type d'élément en alerte (component ou resource).

```python
from canopsis.common.link_builder.link_builder import HypertextLinkBuilder


class MySpecificLinkBuilder(HypertextLinkBuilder):
    # ...

    def build(self, entity, options={}):
        if entity['type'] == 'component':
            label = u'Fiche serveur'
            url = 'https://wiki.example.com/?server={}'.format(entity['name'])
        elif entity['type'] == 'resource':
            label = u'Consigne service'
            url = 'https://wiki.example.com/?server={}&service={}'.format(
                entity['component'], entity['name'])
        else:
            # Ni component ni resource ? alors ne rien renvoyer
            return {}
        return {
            u'Wiki': [
                {
                    'label': label,
                    'link': url
                }
            ]
        }
```

#### Récupérer l'alarme en cours

L'alarme est passée parmi le dictionnaire `options` de la méthode `build()`.

```python
# ...

class MySpecificLinkBuilder(HypertextLinkBuilder):
    # ...

    def build(self, entity, options={}):
        alarm = options.pop('alarm', None)
        if not alarm:
            return {}
        last_update = alarm['v']['last_update_date']
        if 946681200 <= last_update < 978303600:
            # Servons en lien en rapport avec la date d'apparition de l'alarme
            return {
                u'Histoire': [
                    {
                        'label': u"Avez-vous pensé au bug de l'an 2000 ?",
                        'link': 'https://en.wikipedia.org/wiki/Y2K'
                    }
                ]
            }
        return {}
```

#### Récupérer une propriété d'un objet Canopsis

Une fois qu'une entité ou une alarme est récupérée, il peut être souhaité
d'accéder à des propriétés de l'objet (telles que les valeurs dans la structure
de données `infos` ou de nombreux champs de l'alarme).

On trouvera généralement plus commode de le faire en utilisant la notation
pointée, par exemple :

- dans une entité, une info comme `infos.crit.value`
- dans une alarme, un champ comme `v.state.val`

Il est possible de travailler avec cette notation pointée dans le code d'un
*linkbuilder* en utilisant la méthode `get_sub_key` qui se trouve dans les
utilitaires communs de Canopsis.

```python
from canopsis.common.utils import get_sub_key

# ...

alarm_state = get_sub_key(alarm, 'v.state.val')

# Aussi possible d'avoir une valeur par défaut si le champ n'existe pas
entity_crit = get_sub_key(entity, 'infos.crit.value', 'NA')
```

#### Encodage des caractères

Pour toute chaîne en dur dans le code avec des caractères non-ASCII, on
veillera à privilégier les objets Unicode (Python 2) :

```python
label = 'Accéder au ticket'   # soucis en vue
label = u'Accéder au ticket'  # ok
```

### Exploitation de la configuration

La configuration du *linkbuilder*, réalisée par l'API *associativetable*,
peut définir un ensemble libre de clefs-valeurs.

Grâce à ces possibilités de configuration, on peut par exemple avoir un seul
code de *linkbuilder* sur plusieurs instances Canopsis (dev, prod) et insérer
des éléments différents en fonction de l'environnement (bases d'URL, valeurs
diverses à placer dans les liens).

Pour utiliser cette configuration, on utilise l'attribut `options` de
l'instance du *linkbuilder*.

Par exemple, si l'on rend une base d'URL configurable et que l'on s'attend à
avoir cette configuration dans *associativetable* :

```json
{
    "my_specific_link_builder" : {
        "base_url" : "http://wiki-dev.lan/"
    }
}
```

… alors on pourra récupérer cette `base_url` lors de la construction de liens
dans le *linkbuilder* en procédant ainsi :

```python
DEFAULT_BASE_URL = 'http://wiki-prod.example.com'

    def build(self, entity, options={}):
        base_url = self.options.get('base_url', DEFAULT_BASE_URL)
        # ...
```

Note : différencier les deux dictionnaires d'options

- `self.options`, attribut d'instance, donne accès aux paramètres du
linkbuilder ;
- `options`, argument passé à la méthode `build()`, contient des éléments de
contexte propres à chaque sollicitation du linkbuilder (exemple : l'alarme en
cours).

Il est possible de fusionner les deux dictionnaires dès le début de chaque
traitement avec la fonction `merge_two_dicts` :

```python
from canopsis.common.utils import merge_two_dicts

# ...
    def build(self, entity, options={}):
        opt = merge_two_dicts(self.options, options)

        base_url = opt.get('base_url', DEFAULT_BASE_URL)
        alarm = opt.get('alarm')
        # ...
```

### Gestion des logs

Il peut être souhaité de journaliser des erreurs lors de la construction des
liens.

Attention cependant au nombre d'instructions de log et au niveau de criticité
des messages : lorsqu'un cas de log est systématique, il faut avoir à l'esprit
que le fichier de log cible va se retrouver alimenté avec cette ligne pour
chaque appel au linkbuilder (c'est-à-dire à chaque récupération d'une alarme).
Il convient donc de ne faire de messages de log que dans des proportions utiles.

Exemple de connexion à un logger :

```python
# ...
from canopsis.logger import Logger
# ...


LOG_PATH = 'var/log/linkbuilder.log'


class MySpecificLinkBuilder(HypertextLinkBuilder):
    def __init__(self, options={}):
        super(MySpecificLinkBuilder, self).__init__(options=options)
        self.logger = Logger.get('linkbuilder', LOG_PATH)
        # ...

    def build(self, entity, options={}):
        if False:
            self.logger.error('Where is truth now?')
        return {}
```

Les messages de log du *linkbuilder* pourront alors être consultés dans le
fichier de log au chemin renseigné, sur la machine où le service `oldapi` de
Canopsis est exécuté.

## Mettre en service le linkbuilder

Placer le fichier Python dans le dossier `common/link_builder/` des sources de
Canopsis. Redémarrer le service `oldapi`.

Préparer et envoyer la configuration via l'API *associativetable*.

La clef doit porter le nom du module Python contenant le *linkbuilder* à
activer (nom du fichier).

Exemple :

- Nom du fichier Python : `my_specific_link_builder.py`
- Clef à écrire pour la config : `my_specific_link_builder`

```json
{
    "my_specific_link_builder" : {
        "base_url" : "http://wiki-dev.lan/",
        "tout_autre_parametre": "..."
    }
}
```

Se référer au [guide d'administration Linkbuilder][0] pour la requête API
attendue.
