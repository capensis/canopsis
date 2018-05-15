## Authentification

### Utilisation de l'authentification

Afin d'utiliser les API de canopsis, il est nécessaire de s'authentifier.
Il existe trois solutions :

  - ouvrir une session canopsis dans son navigateur web ;
  - utiliser un token d'authentification ;
  - utiliser le mécanisme d'authentification HTTP Basic.

Nous ne détaillerons que les deux dernières méthodes.

#### Token d'authentification

Pour récupérer le token d'authentification associé à une session, il faut se
connecter à canopsis et afficher la modale "profile" en cliquant sur l'icône
avec le nom d'utilisateur (en haut à droite).

Avec ce dit token, vous pouvez ensuite vous créer une session canopsis en
envoyant une requête HTTP GET à l'URL suivante :
```
http://<host>:<port>/autologin?authkey=<authkey>
```

La requête retourne un cookie d'authentification qui doit être stocké et
utilisé pour les requêtes suivantes (l'option `-b` avec curl).

Si votre client HTTP ne supporte pas les sessions HTTP, vous devrez
utiliser l'authentification HTTP Basic.


#### HTTP Basic authentication
Pour utiliser l'authentification HTTP Basic, il suffit juste d'ajouter
dans l'URL interrogée, l'identifiant et le mot de passe associé en suivant
l'exemple suivant :

```
http://user:password@<host>:<port>/api/v2/context
```

En fonction du client, d'autres actions peuvent être nécessaire.

Une fois la session crée, elle reste valide pendant un certain temps. Vous
n'avez pas besoin de mettre identifiant et mot de passe à chaque nouvelle
requête. Néanmoins, certains clients HTTP ne permettent pas de gérer les
sessions HTTP, vous devrez donc vous référer à la documentation pour savoir
si vous pouvez utiliser cette fonctionnalité. Dans le cas contraire, vous
devrez saisir l'identifiant et le mot de passe associé à chaque requête.

#### Cloture d'une session

Pour fermer une session canopsis en cours, il faut juste faire une
requête HTTP GET à l'URL suivante :
```
http://<host>:<port>/logout
```

Vous n'avez pas besoin de fermer uns session canopsis si votre client HTTP
ne gère pas les sessions HTTP.

### Création d’un nouveau backend

Une nouvelle méthode d’authentification de Canopsis est relativement simple à
mettre en place avec les outils actuels :

Pour un backend nommé avec beaucoup d’imagination « NewBackend », voici le
fichier que vous devrez déposer dans `sources/python/webcore/canopsis/webcore/auth/newbackend.py` :

```python
from copy import deepcopy

from bottle import request, HTTPError

from canopsis.auth.base import BaseBackend

class NewBackend(BaseBackend):

    name = 'NewBackend'
    handle_logout = True # seulement si le logout doit effectivement etre gere par vous.

    AUTH_TRACER = 'auth_newtype_of_auth_replace_me'

    def __init__(self, ws, *args, **kwargs):
        super(NewBackend, self).__init__(*args, **kwargs)

        self.ws = ws

        self.session_module = self.ws.require('session')
        self.rights_module = self.ws.require('rights')

    # Fonction requise pour decorer toutes les routes nécessitant une authentification.
    def apply(self, callback, context):
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = self.session.get()

            if s.get('auth_on', False) is True:
                if request.path in ['/logout', '/disconnect']:
                    self.undo_auth(s)

            else:
                username = request.params.get('username', default=None)
                password = request.params.get('password', default=None)
                auth_res = do_auth(s, username, password)
                if not auth_res:
                    self.logger.error('Impossible to authenticate user')
                    return HTTPError(403, 'Forbidden')

            # On appelle effectivement la route demandee
            return callback(*args, **kwargs)

    def do_auth(self, session, username, password):
        # authenticate user...

        if not do_extra_authentication_steps(username, password):
            return False

        user = rights_mgr.get_user(username)
        if not user:
            user = {
                '_id': username,
                'external': True,
                'enable': True,
                'contact': {},
                'role': 'stringclassifiedcrecordselector' # role non privilégié
            }
            self.rights_module.save_user(self.ws, deepcopy(user))
        else:
            user['_id'] = user.get('_id', username)

        session = self.session_module.create(user)
        session[self.AUTH_TRACER] = True
        session.save()

        return True

    def undo_auth(self, session):
        # On vérifie si l'utilisateur a été authentifié par notre méthode
        if session.get(self.AUTH_TRACER, False) is True:
            do_extra_deauthentication_steps()
            session['auth_on'] = False

        return True

def get_backend(ws):
    return NewBackend(ws)
```
