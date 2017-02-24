from unittest import main, TestCase
import canopsis.context_graph.process as process


class Logger(object):
    def debug(self, log):
        print("DEBUG : {0}".format(log))

    def info(self, log):
        print("INFO : {0}".format(log))

    def warning(self, log):
        print("WARNING : {0}".format(log))

    def critical(self, log):
        print("CRITICAL : {0}".format(log))

    def create_event(res, comp, conn_name, conn):
        pass


class Test(TestCase):

    def setUp(self):
        setattr(process, 'LOGGER', Logger())

    def test_update_case_6(self):
        entities_t1 = [{'_id': 'comp_1',
                        'type': 'component',
                        'impact': [],
                        'depends': []},
                       {'_id': 're_1',
                        'type': 'resource',
                        'impact': [],
                        'depends': []}]
        entities_t2 = [{'_id': 'conn_1', 'type': 'connector'},
                       {'_id': 'comp_1', 'type': 'component'},
                       {'_id': 're_1', 'type': 'resource'}]
        ids = {'re_id': 're_1', 'comp_id': 'comp_1', 'conn_id': 'conn_1'}
        self.assertEquals(process.update_case6(entities_t1, ids), 1)
        self.assertEquals(process.update_case6(entities_t2, ids), 0)

    def test_check_type(self):
        entities = {'_id': 'conn_1', 'type': 'connector'}
        self.assertTrue(process.check_type(entities, 'connector'))
        self.assertRaises(TypeError, process.check_type(entities, 'component'))


if __name__ == '__main__':
    main()
