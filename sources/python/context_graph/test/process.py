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


def create_event(conn, conn_name,  comp=None, res=None):
    event = {"connector": conn,
             "connector_name": conn_name,
             }
    if comp is None:
        event["component"] = comp

    if res is None:
        event["resource"] = res
    return event


class Test(TestCase):

    def setUp(self):
        setattr(process, 'LOGGER', Logger())

    def test_update_case_6(self):
        entities_t1 = [{'_id': 'comp_1', 'type': 'component', 'impact': [], 'depends': []}, {'_id': 're_1', 'type':
                                                                                             'resource', 'impact': [],
                                                                                             'depends':[]}]
        entities_t2 = [{'_id': 'conn_1', 'type': 'connector'}, {
            '_id': 'comp_1', 'type': 'component'}, {'_id': 're_1', 'type': 'resource'}]
        ids = {'re_id': 're_1', 'comp_id': 'comp_1', 'conn_id': 'conn_1'}
        self.assertEquals(process.update_case6(entities_t2, ids), 1)
        self.assertEquals(process.update_case6(entities_t2, ids), 0)

    def test_update_case_4(self):
        res_id = "re_1"
        conn_id = "conn_id"
        comp_id = "comp_id"

        event = create_event(conn_id, conn_id, comp_id, res_id)

        process.cache_re.add(res_id)
        process.cache_comp.add(comp_id)
        process.cache_conn.add(conn_id)

        case, ids = process.prepare_update(event)

        self.assertEqual(case, 4)
        self.assertListEqual(ids, {})

        # check cache state
        expected_cache_re = set()
        expected_cache_re.add(res_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id)
        self.assertItemsEqual(process.cache_re, expected_cache_re)
        self.assertItemsEqual(process.cache_comp, expected_cache_comp)
        self.assertItemsEqual(process.cache_conn, expected_cache_conn)




if __name__ == '__main__':
    main()
