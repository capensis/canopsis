
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..base import Driver
from ..py import PyDriver, processcrud, create, read, update, delete
from ..generator import func2crudprocessing, obj2driver, DriverAnnotation
from ...request.core import Request, Context
from ...request.crud.create import Create
from ...request.crud.read import Read
from ...request.crud.update import Update
from ...request.crud.delete import Delete


class Func2CrudProcessingTest(UTCase):

    def test_function(self):

        def func(a, b):

            return [a + b]

        genfunc = func2crudprocessing(func)

        crud = Create(None, {'a': 1})

        request = Request(ctx=Context({'b': 2}))

        _request = genfunc(crud=crud, request=request)

        self.assertIs(_request, request)
        self.assertEqual(_request.ctx[crud], [3])

    def test_object(self):

        class Test(object):

            def test(self, a, b):

                return a + b

        result = func2crudprocessing(Test)

if __name__ == '__main__':
    main()
