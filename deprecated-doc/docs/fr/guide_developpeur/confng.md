## ConfNG

 * Projet : `canopsis/canopsis` -> `sources/python/confng`

### Description

Module facilitant la récupération de la configuration dans Canopsis.

### Utilisation

```python
from canopsis.confng import Configuration, Ini

class Worker(object):

    def __init__(self, *args, **kwargs):
        self.config = Configuration.load('etc/conf.ini', Ini)

        # shortcuts for sections
        conf_section = self.config.get('section', {})

        # optional conf
        self.value = conf_section.get('key', 'def_value')
        self.other_value = conf_section.get('okey', 'def_value')

        # mandatory conf
        try:
            self.required_value = conf_section.get('required_key')
            self.other_req_value = conf_section.get('other_req_key')
        except KeyError, ex:
            work_with_exception(ex)

    def work(self):
        return self.required_value, self.other_req_value, self.value
```

Règles d’utilisation proposées :

 * Préférez la récupération des sections dans des variables distinctes : on évite les `config.get('section', {})` répétitifs. Biensûr, adapter le nom : l’utilisation de "section" ici est générique.
 * L’objet `config` est en fait un simple dictionnaire. Néanmoins, passez toujours par `.get('key', DEF_VAL)` : cela évite les problèmes liés à de la configuration absente.
 * Le plus possible, permettre à votre configuration d’avoir des valeurs par défaut.

