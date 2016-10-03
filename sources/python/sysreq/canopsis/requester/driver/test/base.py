
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..base import Driver


class DriverTest(UTCase):

    def setUp(self):

        self.requests = []
        self.lkwargs = []

        class TestDriver(Driver):

            name = 'test'

            def process(self_, request, **kwargs):

                self.requests.append(request)
                self.lkwargs.append(kwargs)

                return request

        self.drivercls = TestDriver
        self.driver = TestDriver()

    def test_class_name(self):

        self.assertEqual(self.drivercls.name, 'test')

    def test_custom_name(self):

        self.assertEqual(self.drivercls(name='example').name, 'example')

    def test_kwargs(self):

        kwargs = {'bar': 'foo'}

        result = self.driver.process(request=True, **kwargs)

        self.assertEqual(self.requests, [True])
        self.assertEqual(self.lkwargs, [kwargs])

if __name__ == '__main__':
    main()
