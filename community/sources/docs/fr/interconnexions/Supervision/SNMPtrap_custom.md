# Traps SNMP Custom

Réceptionne des traps SNMP, les traduit grâce à des traitements spécifiques
souhaités et les convertit en évènements.

!!! info
    Ce connecteur n'est disponible que dans l'édition Pro de Canopsis.

## Fonctionnement

Le moteur SNMP permet le traitement des traps SNMP récupérés par le connecteur
`snmp2canopsis`, grâce à des règles de correspondance utilisant les MIB.
En l'abscence de MIB, il est possible d'effectuer un traitement spécifique
en passant par une classe Python construite pour l'occasion.

Lorsqu'aucune règle de correspondance du moteur SNMP ne correspond au trap reçu,
le moteur de traitement custom se déclenche.

Le moteur va alors utiliser les différentes classes « custom » à sa disposition
pour identifier le trap reçu (fonction `match` qui retourne `True` ou `False`).

Toute classe qui « reconnaît » ainsi le trap est sélectionnée pour exécuter la
seconde opération : la fonction `build_event`, qui construit l'évènement envoyé
à Canopsis.

## Conception d'une classe custom

Toute logique de reconnaissance et traduction de trap en un évènement Canopsis
doit faire l'objet d'une classe Python.

La classe doit hériter de `SnmpTrap` et implémenter au moins les méthodes :

- `match(trap)`

    Cette méthode doit retourner `True` ou `False`.

    Le paramètre `trap` passé correspond au JSON produit en sortie du
    connecteur `snmp2canopsis`.

- `build_event(trap)`

    Cette méthode doit en principe retourner un évènement au format attendu par
    Canopsis.

    Il est aussi possible de retourner une valeur fausse (le booléen `False` ou
    toute autre valeur évaluée comme fausse en Python) si au dernier moment,
    lors de l'étape `build_event`, on décide de ne pas produire d'évènement.

    Le paramètre `trap` est toujours le dictionnaire sorti par `snmp2canopsis`.
    Il est en fait déjà construit pour respecter la forme d'un évènement
    Canopsis, dont il ne reste qu'à corriger, transformer ou enrichir les
    valeurs dans `build_event`.

    Conseil : utiliser la méthode `format_event(trap)` pour s'assurer de la
    présence des attributs obligatoires avant de retourner l'évènement.

Squelette minimal de classe :

```python
from canopsis_cat.snmp.custom_trap import SnmpTrap

class MyCustomTrap(SnmpTrap):

    DELIMITER = '#'
    RULE_TAG = 'MYRULENAME'
    STOP_AT_FIRST_MATCH = False

    def match(self, trap):
        """
        :rtype: bool
        """
        return False

    def build_event(self, trap):
        """
        :rtype: dict
        """
        return self.format_event(trap)
```

### Utilitaires inclus

Des opérations courantes pour l'extraction d'informations dans les traps sont
fournies via les méthodes suivantes :

* `self.date_format(timestamp)`

    Convertit le timestamp donné en chaîne de caractères formattée.
    Exemple : `self.date_format(62) == '1970-01-01 00:01:2'`

* `self.trap_oid(trap)`

    Donne l'OID du trap.
    Exemple : `self.trap_oid({ 'snmp_trap_oid': '6.2' }) == '6.2'`

* `self.var_from_oid(trap, oid)`

    Retrouve un OID dans les variables du trap snmp.
    Exemple : `self.var_from_oid({ 'snmp_vars': {'5.9': 0} }, '5.9') == 0`

* `self.word(variable, index)`

    Permet de récupérer le i-ème élément d'une variable, en la découpant suivant
    un caractère de séparation (# par défaut).
    Exemple : `self.word('a#b#c', 2) == 'b'`

## Ajout de la classe

Le fichier de classe doit être placé dans le dossier
`/opt/canopsis/lib/python2.7/site-packages/canopsis_cat/snmp/custom_handler`,
là où tourne le moteur `snmp`.

Dans le cas d'un environnement Docker, il est conseillé de créer un volume pour
ce répertoire, ce qui facilitera le placement des fichiers.

Une fois le fichier ajouté ou modifié, redémarrer le moteur `snmp` de Canopsis.

## Ordre et options de traitement

Les modules Python déposés dans le dossier `custom_handler` sont chargés
d'après l'ordre des noms de fichiers (ordre lexicographique). Lors du traitement
d'un trap, les classes sont « exécutées » dans ce même ordre.

Par défaut, si plusieurs classes reconnaissent un même trap, le mécanisme de
traitement de traps personnalisés peut produire plusieurs évènements pour ce
même trap. Ce comportement peut être modifié sur chaque classe en définissant
l'attribut de classe `STOP_AT_FIRST_MATCH` (booléen). Lorsque cet attribut est
vrai pour une classe et que cette classe a reconnu un trap, aucune autre classe
suivante ne sera essayée pour ce trap.

## Journalisation

Au sein de la classe custom, un [logger Python][logging] est utilisable,
les messages se trouveront alors dans le log du moteur snmp.

Le logger est accessible via `self.logger` :

```python
from canopsis_cat.snmp.custom_trap import SnmpTrap

class MyCustomTrap(SnmpTrap):

    DELIMITER = '#'
    RULE_TAG = 'MYRULENAME'
    STOP_AT_FIRST_MATCH = False

    def match(self, trap):
        """
        :rtype: bool
        """
        self.logger.info("It's a trap! but it will never match.")
        return False

    # ...
```

[logging]: https://docs.python.org/2.7/howto/logging.html

## Outil de test de classe

Il est possible de tester le comportement d'une règle custom pour vérifier la
bonne détection et la bonne génération d'un évènement Canopsis, d'après le
JSON qu'envoie le connecteur `snmp2canopsis`.

Cela évite d'attendre l'arrivée ou la reproduction de « vrais » traps.

Pour ce faire, écrire le JSON donné par `snmp2canopsis` dans un fichier et
invoquer le script `custom-trap-tester.py` avec :

- l'option `-t` et le chemin vers le fichier de trap (ci-dessous, `trap.json`)
- l'option `-c` et le nom complet du module qui fournit la classe à tester
(ci-dessous, pour un fichier déposé sous le nom `mycustomtrap.py`, le nom
complet sera `canopsis_cat.snmp.custom_handler.mycustomtrap`)
- éventuellement, l'option `--publish` pour déclencher la publication réelle de
l'evènement produit vers le bus de messages habituel de Canopsis

Le script `custom-trap-tester.py` est disponible dans le `PATH` de
l'utilisateur `canopsis`. Dans un environnement Docker, il faut se placer dans
le conteneur du moteur `snmp`.

```console
custom-trap-tester.py -t trap.json -c canopsis_cat.snmp.custom_handler.mycustomtrap [--publish]
```

## Test complet de mise en œuvre

Pour cet exemple, on considère le cas où l'on doit gérer un trap spécifique
« Guitare ».

Voici le contenu du JSON pour ce trap, tel que produit par le connecteur
`snmp2canopsis` :

```json
{
   "component":"172.18.0.1",
   "connector":"snmp",
   "connector_name":"snmp2canopsis",
   "event_type":"trap",
   "snmp_trap_oid":"1.3.6.1.4.1.20006.1.12",
   "snmp_vars":{
      "1.3.6.1.2.1.1.3.0":"527512",
      "1.3.6.1.4.1.20006.1.3.1.54":"607",
      "1.3.6.1.4.1.20006.1.3.1.55":"it#doesnt#matter#star#guitare#the#test",
      "1.3.6.1.6.3.1.1.4.1.0":"1.3.6.1.4.1.20006.1.12"
   },
   "snmp_version":"2c",
   "source_type":"component",
   "state":3,
   "state_type":1,
   "timestamp":1560936165.165884
}
```

Spécifications imaginaires :

> On considère que la reconnaissance de ce trap repose sur l'OID de trap
> `1.3.6.1.4.1.20006.1.12` et la présence du mot `guitare` en cinquième position
> de la variable à l'OID `1.3.6.1.4.1.20006.1.3.1.55`, cette variable étant une
> chaîne séparable par `#`.
>
> Pour créer l'évènement Canopsis, on prend le nombre dans la variable d'OID
> `1.3.6.1.4.1.20006.1.3.1.54`. Une valeur strictement supérieure à 100 est
> à interpréter comme CRITICAL, la valeur 0 est à interpréter comme OK, toute
> autre valeur comme MINOR.

Code de la classe associée :

```python
from __future__ import unicode_literals
from canopsis_cat.snmp.custom_trap import SnmpTrap


class GuitarTrap(SnmpTrap):
    """
    Snmp Trap handling class for guitar traps.
    """

    DELIMITER = '#'
    RULE_TAG = 'letagdeguitare'

    def match(self, trap):
        """
        Tell if a given trap will be treated by this class.

        :param dict trap: a snmp trap
        :rtype: bool
        """
        oid = self.trap_oid(trap)
        message = self.var_from_oid(trap, '1.3.6.1.4.1.20006.1.3.1.55')
        if oid == '1.3.6.1.4.1.20006.1.12' and self.word(message, 5) == 'guitare':
            return True

        return False

    def build_event(self, trap):
        """
        Take a trap, an process it to create an event.

        :param dict trap: a snmp trap
        :rtype: dict
        """
        message = self.var_from_oid(trap, '1.3.6.1.4.1.20006.1.3.1.55')
        value = int(self.var_from_oid(trap, '1.3.6.1.4.1.20006.1.3.1.54'))

        trap['output'] = self.word(message, 7)
        trap['component'] = 'Gibson'
        trap['resource'] = self.word(message, 3)

        if value > 100:
            trap['state'] = 3
        elif value == 0:
            trap['state'] = 0
        else:
            trap['state'] = 1

        return self.format_event(trap)
```

On place ce fichier sous le nom `guitare.py` dans
`/opt/canopsis/lib/python2.7/site-packages/canopsis_cat/snmp/custom_handler/`.

On redémarre le service ou conteneur `snmp` Canopsis. Cette étape dépend du
type d'installation. Par exemple, sur une installation Docker Compose :

```sh
docker-compose restart snmp
```

Pour la validation de fonctionnement de la classe avec le script, placer
le contenu du trap en JSON dans un fichier sur le serveur Canopsis ou sur
le conteneur `snmp`. L'exécution du script et son résultat se présentent comme
suit :

```console
(canopsis x.yy.z)[canopsis@ct ~]$ custom-trap-tester.py -t /tmp/event.json -c canopsis_cat.snmp.custom_handler.guitare
* Found class GuitarTrap
* Trap can be handled by this class. Building event...
* The following event has been generated: 
{
    "_id": "snmp.snmp2canopsis.check.resource.Gibson.matter", 
    "component": "Gibson", 
    "connector": "snmp", 
    "connector_name": "snmp2canopsis", 
    "event_type": "check", 
    "output": "test", 
    "resource": "matter", 
    "source_type": "resource", 
    "state": 3, 
    "state_type": 1, 
    "timestamp": 1560936165.165884
}
```

Ceci confirme la bonne détection du trap et la bonne transformation des infos
pour produire l'évènement Canopsis attendu.

Après validation du comportement de la classe custom, on peut constater le bon
fonctionnement de la chaîne complète en envoyant réellement le trap au
connecteur `snmp2canopsis`.

```console
snmptrap -v 2c -c public ${IP_RECEPTEUR} '' 1.3.6.1.4.1.20006.1.12 \
  1.3.6.1.4.1.20006.1.3.1.54 i 607 \
  1.3.6.1.4.1.20006.1.3.1.55 s 'it#doesnt#matter#star#guitare#the#test'
```

Dans les logs du service ou conteneur `snmp2canopsis`, on retrouve le JSON
montré en début de test.

Dans les logs du service ou conteneur du moteur `snmp`, on peut observer la
ligne suivante qui indique le déclenchement de notre *handler* custom :

```
2019-06-20 15:29:36,246 INFO snmp [snmp 109] Trap handled by custom handler : letagdeguitare
```
