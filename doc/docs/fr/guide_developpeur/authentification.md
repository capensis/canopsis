## Authentification

### Création d’un nouveau backend

Une nouvelle méthode d’authentification de Canopsis est relativement simple à mettre en place avec les outils actuels :

Pour un backend nommé avec beaucoup d’imagination « NewBackend », voici le fichier que vous devrez déposer dans `sources/python/webcore/canopsis/webcore/auth/newbackend.py` :

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
