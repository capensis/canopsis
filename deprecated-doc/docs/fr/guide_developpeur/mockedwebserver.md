## MockedWebServer

 * Project : `canopsis/canopsis` -> `sources/canopsis/common/webserver_mocking.py`

### Description

Classe générique pour permettre de mocker un webservice facilement. Très utile pour les tests unitaires.

Elle nécessite l'écriture d'un handler spécifique pour décrire le comportement des routes désirées.

### Utilisation

Ligne d'import :

```python
from canopsis.common.mocking_webserver import MockedWebService
```

Utilisation minimaliste :

```python
import json
import re
import requests
from SimpleHTTPServer import SimpleHTTPRequestHandler
import urllib2

from canopsis.common.webserver_mocking import MockedWebServer

class MonRequestHandler(SimpleHTTPRequestHandler):
    DUO_PATTERN = re.compile(r'/GetDuo')

    def do_GET(self):
        if re.search(self.DUO_PATTERN, self.path):
            self.send_response(requests.codes.ok)
            self.send_header('Content-Type', 'application/json; charset=utf-8')
            self.end_headers()

            dico = {'Titi': 'Gros Minet'}
            response_content = json.dumps(dico)
            self.wfile.write(response_content.encode('utf-8'))
            return

        self.send_response(requests.codes.not_found)
        self.end_headers()

port = MockedWebServer.get_free_port()
webserver = MockedWebServer(handler=MonRequestHandler, port=port)

webserver.start()

# On peut maintenant envoyer des requêtes sur notre nouvelle route
urllib2.urlopen('http://localhost:{}/GetDuo'.format(port)).read()
# '{"Titi": "Gros Minet"}'

webserver.shutdown()
```

### Fonctions intéréssantes

 - **get_free_port()**: (static) permet de deviner un numéro de port disponible sur localhost
 - **start()**: pour lancer le webserveur
 - **shutdown()**: pour éteindre le webserveur
