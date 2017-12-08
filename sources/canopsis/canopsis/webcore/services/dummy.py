"""Dummy auth routes."""

import copy

from canopsis.auth.dummy import DummyAuthBackend
from canopsis.webcore.services import session, rights

from bottle import request, redirect

import copy

USERNAME = "username"
PASSWORD = "password"


def exports(ws):

    @ws.application.get(DummyAuthBackend.DUMMY_LOGIN_URL, skip=ws.skip_login)
    def dummy_login():
        page = """<html>
                   <body>
                       <form action=\"{0}\" method="post" >
                           <label>Username</label>
                           <input type="textname=" name="username" required>
                           <label>Password</label>
                           <input type="password" name="password" required>
                           <button type="submit">Login</button>
                       <form>
                   <body>
               <html>"""
        page = page.format(DummyAuthBackend.DUMMY_VALIDATE_URL)
        return page

    @ws.application.post(DummyAuthBackend.DUMMY_VALIDATE_URL, skip=ws.skip_login)
    def dummy_validate_login():
        try:
            username = request.forms.get(USERNAME)
            password = request.forms.get(PASSWORD)
        except KeyError:
            ws.logger.debug("Dummy auth : Missing username or password")
            redirect(DummyAuthBackend.DUMMY_LOGIN_URL)

        if username in DummyAuthBackend.BLACKLIST_USERS:
            redirect(DummyAuthBackend.DUMMY_LOGIN_URL)

        user = rights.get_manager().get_user(copy.deepcopy(username))

        ws.logger.error("User : {}\n".format(user))

        if user is None:
            if not user:
                user = {
                    '_id': username,
                    'external': True,
                    'enable': True,
                    'contact': {},
                    'role': 'stringclassifiedcrecordselector'
                }
            rights.save_user(ws, copy.deepcopy(user))
        else:
            # fix for canopsis_cat.webcore.services.session.get_user() that need
            # the _id key to work...
            user['_id'] = user.get('_id', username)

        sess = session.create(user)
        sess['auth_dummy'] = True
        sess['dummy_username'] = username
        sess.save()

        redirect("/")
