
"""conf file driver UTs."""

from canopsis.utils.ut import UTCase

from unittest import main

from ..read import Read


class ReadTest(UTCase):

    def test_init_default(self):

        read = Read()

        self.assertIsNone(read.select())
        self.assertIsNone(read.offset())
        self.assertIsNone(read.limit())
        self.assertIsNone(read.orderby())
        self.assertIsNone(read.groupby())
        self.assertFalse(read.join())

    def test_init(self):

        select = 'select'
        offset = 1
        limit = 2
        orderby = 'orderby'
        groupby = 'groupby'
        join = 'join'

        read = Read(
            select=(select,), offset=offset, limit=limit, orderby=(orderby,),
            groupby=(groupby,), join=join
        )

        self.assertEqual((select,), read.select())
        self.assertEqual(offset, read.offset())
        self.assertEqual(limit, read.limit())
        self.assertEqual((orderby,), read.orderby())
        self.assertEqual((groupby,), read.groupby())
        self.assertEqual(join, read.join())

    def test_init_error(self):

        select = 0
        offset = '1'
        limit = '2'
        orderby = 0
        groupby = 0
        join = 0

        #self.assertRaises(TypeError, Read, select=select)
        self.assertRaises(TypeError, Read, offset=offset)
        self.assertRaises(TypeError, Read, limit=limit)
        #self.assertRaises(TypeError, Read, orderby=orderby)
        #self.assertRaises(TypeError, Read, groupby=groupby)
        #&self.assertRaises(TypeError, Read, join=join)

    def test_chaining(self):

        select = 'select'
        offset = 1
        limit = 2
        orderby = 'orderby'
        groupby = 'groupby'
        join = 'join'

        read = Read()

        readbis = read.select(select).offset(offset).limit(limit)
        readbis = readbis.orderby(orderby).groupby(groupby).join(join)

        self.assertIs(read, readbis)

        self.assertEqual((select,), read.select())
        self.assertEqual(offset, read.offset())
        self.assertEqual(limit, read.limit())
        self.assertEqual((orderby,), read.orderby())
        self.assertEqual((groupby,), read.groupby())
        self.assertEqual(join, read.join())

if __name__ == '__main__':
    main()
