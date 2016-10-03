
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..base import CRUDElement
from ...core import Request


class CRUDTest(UTCase):

    def test_init_defaul(self):

        crud = CRUDElement()

        self.assertIsNone(crud.request)

    def test_init(self):

        request = Request()

        crud = CRUDElement(request=request)

        self.assertIs(request, crud.request)

    def test__call__(self):

        requests = []

        class Driver(object):

            def process(self, request, **kwargs):

                requests.append(request)

        request = Request(driver=Driver())

        crud = CRUDElement(request=request)
        crud()

        self.assertIn(request, requests)

        self.assertIn(crud, request.cruds)

if __name__ == '__main__':
    main()
