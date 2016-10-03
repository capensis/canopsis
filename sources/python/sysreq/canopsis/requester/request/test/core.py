
"""request.base UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..base import BaseElement
from ..core import Request, Context
from ..expr import Expression as E, Function as F, FuncName as FN
from ..crud.create import Create
from ..crud.read import Read
from ..crud.update import Update
from ..crud.delete import Delete


class ContextTest(UTCase):

    def test_init(self):

        self.assertFalse(Context())

    def test_init_params(self):

        context = Context({'a': 1})

        self.assertTrue(context)

        self.assertIn('a', context)

    def test_crud(self):

        context = Context()
        crud = BaseElement()

        self.assertNotIn(crud, context)

        context[crud] = 1

        self.assertIn(crud, context)
        self.assertIn(crud.ctxname, context)

        self.assertEqual(context[crud], 1)
        self.assertEqual(context[crud.ctxname], 1)

        del context[crud]

        self.assertNotIn(crud, context)
        self.assertNotIn(crud.ctxname, context)

    def test_fill(self):

        context = Context({
            'a.b': [1, 2, 4],
            'c': [5]
        })

        _context = Context({
            'b': [3],
            'c': [4],
            'a': [5]
        })

        context.fill(_context)

        self.assertEqual(context['a'], [5])
        self.assertEqual(context['b'], [3])
        self.assertEqual(context['c'], [5, 4])
        self.assertEqual(context['a.b'], [1, 2, 4, 3])


class RequestTest(UTCase):

    def setUp(self):

        self.requests = []

        class TestDriver(object):

            def process(_self, request, **kwargs):

                self.requests.append(request)

        self.driver = TestDriver()

        self.request = Request(driver=self.driver)

    def test_init_default(self):

        request = Request()

        self.assertIsNone(request.driver)
        self.assertEqual(request.ctx, {})
        self.assertIsNone(request.query)
        self.assertEqual(request.cruds, [])

    def test_init_errorquery(self):

        self.assertRaises(TypeError, Request, query=1)

    def test_init_errorcruds(self):

        self.assertRaises(TypeError, Request, cruds=[1])

    def test_and_query(self):

        self.request.query = E.A

        self.assertEqual(self.request.query.name, 'A')

        self.request.query &= E.B

        self.assertEqual(self.request.query.name, FN.AND.value)
        self.assertEqual(
            self.request.query.params[0].name,
            'A'
        )
        self.assertEqual(
            self.request.query.params[1].name,
            'B'
        )

    def test_or_query(self):

        self.request.query = E.A

        self.assertEqual(self.request.query.name, 'A')

        self.request.query |= E.B

        self.assertEqual(self.request.query.name, FN.OR.value)
        self.assertEqual(
            self.request.query.params[0].name,
            'A'
        )
        self.assertEqual(
            self.request.query.params[1].name,
            'B'
        )

    def test_del_query(self):

        self.request.query = E.A

        self.assertEqual(self.request.query.name, 'A')

        del self.request.query

        self.assertIsNone(self.request.query)

    def test_and__query(self):

        self.request.query = E.A

        self.assertEqual(self.request.query.name, 'A')

        request = self.request.and_(E.B)

        self.assertIs(request, self.request)

        self.assertEqual(self.request.query.name, FN.AND.value)
        self.assertEqual(
            self.request.query.params[0].name,
            'A'
        )
        self.assertEqual(
            self.request.query.params[1].name,
            'B'
        )

    def test_or__query(self):

        self.request.query = E.A

        self.assertEqual(self.request.query.name, 'A')

        request = self.request.or_(E.B)

        self.assertIs(request, self.request)

        self.assertEqual(self.request.query.name, FN.OR.value)
        self.assertEqual(
            self.request.query.params[0].name,
            'A'
        )
        self.assertEqual(
            self.request.query.params[1].name,
            'B'
        )

    def test_commit(self):

        self.request.commit()

        self.assertEqual(self.requests, [self.request])

    def test_select(self):

        value = 'test'
        read = self.request.select(value)

        self.assertIsInstance(read, Read)
        self.assertEqual(read.select(), (value,))

    def test_offset(self):

        value = 1
        read = self.request.offset(value)

        self.assertIsInstance(read, Read)
        self.assertEqual(read.offset(), value)

    def test_limit(self):

        value = 1
        read = self.request.limit(value)

        self.assertIsInstance(read, Read)
        self.assertEqual(read.limit(), value)

    def test_groupby(self):

        value = 'test'
        read = self.request.groupby(value)

        self.assertIsInstance(read, Read)
        self.assertEqual(read.groupby(), (value,))

    def test_orderby(self):

        value = 'test'
        read = self.request.orderby(value)

        self.assertIsInstance(read, Read)
        self.assertEqual(read.orderby(), (value,))

    def test_join(self):

        value = 'full'
        read = self.request.join(value)

        self.assertIsInstance(read, Read)
        self.assertEqual(read.join(), value)

    def test_processcrud(self):

        cruds = [
            Create('create', {}),
            Read(),
            Update('update', {}),
            Delete()
        ]

        self.request.processcrud(*cruds)

        self.assertIn(self.request, self.requests)
        self.assertEqual(self.request.cruds, cruds)

    def test_create(self):

        name = 'test'
        values = {'a': 1, 'b': 2}

        self.request.create(name, **values)

        self.assertIn(self.request, self.requests)
        crud = self.request.cruds[0]

        self.assertIsInstance(crud, Create)
        self.assertIs(crud.request, self.request)
        self.assertEqual(crud.name, name)
        self.assertEqual(crud.values, values)

    def test_read(self):

        select = ('test',)
        limit = 1

        self.request.read(select=select, limit=limit)

        self.assertIn(self.request, self.requests)
        crud = self.request.cruds[0]

        self.assertIsInstance(crud, Read)
        self.assertIs(crud.request, self.request)
        self.assertEqual(crud.select(), select)
        self.assertEqual(crud.limit(), limit)

    def test_update(self):

        name = 'test'
        values = {'a': 1, 'b': 2}

        self.request.update(name, **values)

        self.assertIn(self.request, self.requests)
        crud = self.request.cruds[0]

        self.assertIsInstance(crud, Update)
        self.assertIs(crud.request, self.request)
        self.assertEqual(crud.name, name)
        self.assertEqual(crud.values, values)

    def test_delete(self):

        names = ('test',)

        self.request.delete(*names)

        self.assertIn(self.request, self.requests)
        crud = self.request.cruds[0]

        self.assertIsInstance(crud, Delete)
        self.assertIs(crud.request, self.request)
        self.assertEqual(crud.names, names)

    def test_exe(self):

        name = 'test'
        params = (1, 2)

        self.request.exe(name, *params)

        self.assertIn(self.request, self.requests)
        crud = self.request.cruds[0]

        self.assertIsInstance(crud, Read)
        self.assertIs(crud.request, self.request)
        self.assertEqual(crud.name, name)
        self.assertEqual(crud.params, params)

if __name__ == '__main__':
    main()
