
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..update import Update


class UpdateTest(UTCase):

    def test_init(self):

        values = {None: None}

        update = Update(name='test', values=values)

        self.assertEqual('test', update.name)
        self.assertEqual(values, update.values)

if __name__ == '__main__':
    main()
