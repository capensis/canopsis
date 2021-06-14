"""Dummy authentication backend"""

from bottle import request, redirect, HTTPError

from canopsis.auth.base import BaseBackend

class DummyAuthBackend(BaseBackend):
    """Handle dummy auth.XS"""

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
            sess = self.session.get()
            if sess.get('auth_on', False):
                if request.path in ['/logout', '/disconnect']:
                    self.undo_auth(sess)

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

    @classmethod
    def undo_auth(cls, session):
        """
        Handle SAML2 Dummy by destroying the session and redirecting the user to
        the login page.
        """
        if session.get('dummy', False) is True:
            session.delete()
            return redirect("/")

def get_backend(ws):
    """Return an instance of DummyAuthBackend."""
    return DummyAuthBackend(ws)
