
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..create import Create


class CreateTest(UTCase):

    def test_init(self):

        values = {None: None}

        create = Create(name='test', values=values)

        self.assertEqual('test', create.name)
        self.assertEqual(values, create.values)

if __name__ == '__main__':
    main()
