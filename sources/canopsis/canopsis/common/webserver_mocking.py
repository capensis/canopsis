#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import requests
import SocketServer
import socket
from threading import Thread

DEFAULT_ADDRESS = 'localhost'


class MockedWebServer(object):
    """
    Generic mocking class for a HTTP server.

    It run an HTTP server in a thread ; you must personnalize GET, POST...
    responses through an hadler class.
    """

    def __init__(self, handler, port):
        """
        :param SimpleHTTPRequestHandler handler: class to handle http requests
        :param int port: a usable port number
        """
        self.port = port

        self._server = SocketServer.TCPServer(('', port), handler)
        self._server.timeout = 1
        self._thread = Thread(target=self.run)
        self._thread.deamon = True

    @staticmethod
    def get_free_port():
        """
        Automatically find a free port to listen to.
        """
        s = socket.socket(socket.AF_INET, type=socket.SOCK_STREAM)
        s.bind((DEFAULT_ADDRESS, 0))
        address, port = s.getsockname()
        s.close()

        return port

    def run(self):
        """
        Run the HTTP server.
        """
        self._server.running = True
        while self._server.running:
            self._server.handle_request()

    def start(self):
        """
        Run the thread.
        """
        self._thread.start()

    def shutdown(self):
        """
        Shutdown the HTTP server and then the thread.
        """
        self._server.running = False
        requests.get('http://{}:{}'.format(DEFAULT_ADDRESS, self.port))
