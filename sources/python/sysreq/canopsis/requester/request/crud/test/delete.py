
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..delete import Delete


class DeleteTest(UTCase):

    def test_init_default(self):

        delete = Delete()

        self.assertIsNone(delete.names)

    def test_init(self):

        names = 1

        delete = Delete(names=names)

        self.assertEqual(names, delete.names)

if __name__ == '__main__':
    main()
