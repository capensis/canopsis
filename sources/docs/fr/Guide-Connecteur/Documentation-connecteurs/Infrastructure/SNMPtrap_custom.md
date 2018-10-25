# /!\ Documentation CAT /!\

# Trap Snmp Custom

L'engine SNMP permet le traitement des traps SNMP (par le biais du connecteur
snmp2canopsis). En l'abscence de MIB, il est désormais possible d'effectuer un
traitement spécifique en passant par une classe python construite pour l'occasion.

## Fonctionnement

Le moteur de traitement custom pour les traps SNMP se déclenche lorsqu'aucune
rule ne correspond à la trap reçue.

Le moteur va alors exécuter une fonction d'identification (match) dans chacune
des classes héritant de SnmpTrap présentes dans un dossier spécifique.

Si l'identification réussie, une seconde fonction est exécutée ; celle-ci va
permettre de construire un event, qui sera envoyé à canopsis.

## Ajouter une classe custom

1. Placer le fichier class dans `/opt/canopsis/lib/python2.7/site-packages/canopsis_cat/snmp/custom_handler`
2. Redémarrer l'engine snmp, avec l'utilisateur canopsis:
```bash
supervisorctl restart amqp2engines:engine-snmp-0
```

### Exemple de class custom

```python
class MyRuleTrap(SnmpTrap):

    DELIMITER = '#'
    RULE_TAG = 'RULENAME'

    def match():
        """
        :rtype: bool
        """
        return False

    def build_event():
        """
        :rtype: dict
        """
        return self.format_event(trap)
```

### Outils inclus

* `self.date_format(timestamp)`: converti le timestamp donné en chaîne de
caractères formattée. Exemple: `self.date_format(62) == '1970-01-01 00:01:2'`
* `self.trap_oid(trap)`: donne l'oid de la trap.
Exemple : `self.trap_oid({ 'snmp_trap_oid': '6.2' }) == '6.2'`
* `self.var_from_oid(trap, oid)`: retrouve un oid dans les variables snmp.
Exemple: `self.var_from_oid({ 'snmp_vars': {'5.9': 0} }, '5.9') == 0`
* `self.word(variable, index)`: permet de récupérer le i-ème élément d'une
variable, en la découpant suivant un caractère de séparation (# par défaut).
Exemple: `self.word('a#b#c', 2) == 'b'`


## Tester une classe custom

Pour tester un trap contre un handler, créer un fichier `trap_event.json`
contenant un json telle qu'envoyé par le connecteur snmp, puis lancer le script
suivant, avec la ligne d'import de classe voulue :

```bash
custom-trap-tester --trap trap_event.json --classe canopsis_cat.snmp.custom_handler.<example>
```

## Exemple de trap envoyé par le connecteur snmp

```json
{
  "connector": "snmp",
  "connector_name": "snmp2canopsis",
  "component": "127.0.0.1",
  "event_type": "trap",
  "source_type": "component",
  "state_type": 1,
  "state": 3,
  "timestamp": 1517909965.44025,
  "snmp_version": "2c",
  "snmp_vars": {
    "6.0.1": 607,
    "9.9.9": "it#doesnt#matter#star#guitare#the#test"
  },
  "snmp_trap_oid": "1.2.3.4"
}
```
