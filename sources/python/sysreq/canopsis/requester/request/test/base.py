
"""request.base UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..base import BaseElement


class BaseElementTest(UTCase):

    def test_init_default(self):

        alias = BaseElement()

        self.assertIsNone(alias.alias)

    def test_init(self):

        alias = BaseElement(alias='alias')

        self.assertEqual(alias.alias, 'alias')

    def test_as_(self):

        alias = BaseElement()

        alias.as_(alias='alias')

        self.assertEqual(alias.alias, 'alias')

    def test_refers(self):

        alias = BaseElement.refers('test')

        self.assertIsInstance(alias, BaseElement)
        self.assertEqual(alias.ctxname, 'test')

    def test_refers_elt(self):

        elt = BaseElement()

        alias = BaseElement.refers(elt)

        self.assertIsInstance(alias, BaseElement)
        self.assertEqual(alias.ctxname, elt.ctxname)

    def test_uuid(self):

        base0, base1 = BaseElement(), BaseElement()

        self.assertNotEqual(base0.uuid, base1.uuid)

        base0, base1 = BaseElement(uuid=1), BaseElement(uuid=1)

        self.assertEqual(base0.uuid, base1.uuid)

    def test_ctxname(self):

        class TestElement(BaseElement):
            pass

        elt = TestElement()

        self.assertEqual(elt.ctxname, elt.uuid)

        elt.name = 1

        self.assertEqual(elt.ctxname, elt.name)

        elt.alias = True

        self.assertEqual(elt.ctxname, elt.alias)

    def test_eq(self):

        class TestElement(BaseElement):
            pass

        base0, base1 = TestElement(), TestElement()

        self.assertNotEqual(base0, base1)

        base0.alias = base1.alias = 1

        self.assertEqual(base0, base1)

        base0.alias = base1.alias = None

        self.assertNotEqual(base0, base1)

        base0.name = base1.name = 1

        self.assertEqual(base0, base1)

        base0.name = base1.name = None

        self.assertNotEqual(base0, base1)

        base0.uuid = base1.uuid

        self.assertEqual(base0, base1)

        base0.uuid = None

        self.assertNotEqual(base0, base1)

if __name__ == '__main__':
    main()
