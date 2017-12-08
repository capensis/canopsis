from bottle import request, redirect, HTTPError

from canopsis.auth.base import BaseBackend

class DummyAuthBackend(BaseBackend):

    BLACKLIST_USERS = ["davros", "mr.rabbit"]
    DUMMY_AUTH = "auth_dummy"
    DUMMY_LOGIN_URL = '/auth/dummy/login/'
    DUMMY_VALIDATE_URL = '/auth/dummy/login/validate'

    def __init__(self, *args, **kwargs):
        super(DummyAuthBackend, self).__init__(*args, **kwargs)

    def apply(self, callback, context):
        """
        Required decorator function for bottle.

        This will be applied every time a auth-ed route is accessed.
        """
        self.setup_config(context)

        def decorated(*args, **kwargs):
            s = self.session.get()
            if s.get('auth_on', False):
                if request.path in ['/logout', '/disconnect']:
                    self.undo_auth(s)

                return callback(*args, **kwargs)

            else:
                auth_res = self.do_auth()
                if not auth_res:
                    self.logger.error(u'Impossible to authenticate user')
                    return HTTPError(403, "Forbidden")

        return decorated

    def do_auth(self):
        """
        Handle SAML2 authentication by redirecting the user to the SSO endpoint.

        After this first redirect, the user will be redirected to the requested
        URL.
        """
        return redirect(self.DUMMY_LOGIN_URL)

    def undo_auth(self, session):
        """
        Handle SAML2 Dummy by destroying the session and redirecting the user to
        the login page.
        """
        if session.get('dummy', False) is True:
            session.delete()
            return redirect("/")

def get_backend(ws):
    return DummyAuthBackend(ws)
