from __future__ import unicode_literals

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
    event = {"connector": conn, "connector_name": conn_name}
    if comp is not None:
        event["component"] = comp

    if res is not None:
        event["resource"] = res
    return event


class Test(TestCase):

    def setUp(self):
        setattr(process, 'LOGGER', Logger())

    def tearDown(self):
        process.cache_comp.clear()
        process.cache_re.clear()
        process.cache_conn.clear()

    def test_preprare_update_case_1(self):
        res_id = "re_id"
        conn_id = "conn_id"
        comp_id = "comp_id"

        event = create_event(conn_id, conn_id, comp=comp_id, res=res_id)

        case, ids = process.prepare_update(event)

        expected_ids = {'comp_id': comp_id,
                        're_id': res_id + "/" + comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        self.assertEqual(case, 1)
        self.assertDictEqual(ids, expected_ids)

        expected_cache_re = set()
        expected_cache_re.add(res_id + "/" + comp_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_2(self):
        res_id = "re_id"
        conn_id = "conn_id"
        comp_id = "comp_id"

        process.cache_conn.add(conn_id + "/" + conn_id)

        event = create_event(conn_id, conn_id, comp=comp_id, res=res_id)

        case, ids = process.prepare_update(event)

        expected_ids = {'comp_id': comp_id,
                        're_id': res_id + "/" + comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        self.assertDictEqual(ids, expected_ids)
        self.assertEqual(case, 2)

        expected_cache_re = set()
        expected_cache_re.add(res_id + "/" + comp_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_2_re_none(self):
        res_id = None
        conn_id = "conn_id"
        comp_id = "comp_id"

        process.cache_conn.add(conn_id + "/" + conn_id)

        event = create_event(conn_id, conn_id, comp=comp_id, res=res_id)

        case, ids = process.prepare_update(event)

        expected_ids = {'comp_id': comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        self.assertDictEqual(ids, expected_ids)
        self.assertEqual(case, 2)

        expected_cache_re = set()
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_3(self):
        res_id = "re_id"
        conn_id = "conn_id"
        comp_id = "comp_id"

        process.cache_comp.add(comp_id)
        process.cache_conn.add(conn_id + "/" + conn_id)

        event = create_event(conn_id, conn_id, comp_id, res_id)

        case, ids = process.prepare_update(event)

        expected_ids = {'comp_id': comp_id,
                        're_id': res_id + "/" + comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        self.assertEqual(case, 3)
        self.assertDictEqual(ids, expected_ids)

        expected_cache_re = set()
        expected_cache_re.add(res_id + "/" + comp_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_4(self):
        res_id = "re_1"
        conn_id = "conn_id"
        comp_id = "comp_id"

        event = create_event(conn_id, conn_id, comp_id, res_id)

        process.cache_re.add(res_id + "/" + comp_id)
        process.cache_comp.add(comp_id)
        process.cache_conn.add(conn_id + "/" + conn_id)

        case, ids = process.prepare_update(event)

        self.assertEqual(case, 4)
        self.assertDictEqual(ids, {})

        # check cache state
        expected_cache_re = set()
        expected_cache_re.add(res_id + "/" + comp_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_5(self):
        res_id = "re_id"
        conn_id = "conn_id"
        comp_id = "comp_id"

        expected_ids = {'comp_id': comp_id,
                        're_id': res_id + "/" + comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        event = create_event(conn_id, conn_id, comp_id, res_id)

        process.cache_comp.add(comp_id)

        case, ids = process.prepare_update(event)

        self.assertEqual(case, 5)
        self.assertDictEqual(ids, expected_ids)

        # check cache state
        expected_cache_re = set()
        expected_cache_re.add(res_id + "/" + comp_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_5_re_none(self):
        res_id = None
        conn_id = "conn_id"
        comp_id = "comp_id"

        expected_ids = {'comp_id': comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        event = create_event(conn_id, conn_id, comp_id, res_id)

        process.cache_comp.add(comp_id)

        case, ids = process.prepare_update(event)

        self.assertEqual(case, 5)
        self.assertDictEqual(ids, expected_ids)

        # check cache state
        expected_cache_re = set()
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_preprare_update_case_6(self):
        res_id = "re_id"
        conn_id = "conn_id"
        comp_id = "comp_id"

        expected_ids = {'comp_id': comp_id,
                        're_id': res_id + "/" + comp_id,
                        'conn_id': conn_id + "/" + conn_id}

        event = create_event(conn_id, conn_id, comp_id, res_id)

        process.cache_comp.add(comp_id)
        process.cache_re.add(res_id + "/" + comp_id)

        case, ids = process.prepare_update(event)

        self.assertEqual(case, 6)
        self.assertDictEqual(ids, expected_ids)

        # check cache state
        expected_cache_re = set()
        expected_cache_re.add(res_id + "/" + comp_id)
        expected_cache_comp = set()
        expected_cache_comp.add(comp_id)
        expected_cache_conn = set()
        expected_cache_conn.add(conn_id + "/" + conn_id)
        self.assertSetEqual(process.cache_re, expected_cache_re)
        self.assertSetEqual(process.cache_comp, expected_cache_comp)
        self.assertSetEqual(process.cache_conn, expected_cache_conn)

    def test_check_type(self):
        re_entity = {'_id': 'conn_1', 'type': 'resource'}
        con_entity = {'_id': 'conn_1', 'type': 'connector'}
        comp_entity = {'_id': 'conn_1', 'type': 'component'}

        self.assertTrue(process.check_type(con_entity, 'connector'))
        self.assertTrue(process.check_type(re_entity, 'resource'))
        self.assertTrue(process.check_type(comp_entity, 'component'))

        with self.assertRaises(TypeError):
            process.check_type(con_entity, "not_a_connector")
        with self.assertRaises(TypeError):
            process.check_type(comp_entity, "not_a_component")
        with self.assertRaises(TypeError):
            process.check_type(re_entity, "not_a_resource")

    def test_update_depends_links(self):
        e_1 = {
            '_id': 'comp_1',
            'type': 'component',
            'impact': [],
            'depends': []
        }
        e_2 = {
            '_id': 'conn_1',
            'type': 'connector',
            'impact': [],
            'depends': []
        }
        process.update_depends_links(e_1, e_2)
        self.assertTrue(e_2['_id'] in e_1['depends'])
        process.update_depends_links(e_1, e_2)
        self.assertTrue(e_1['depends'] == [e_2['_id']])

    def test_update_impact_links(self):
        e_1 = {
            '_id': 'comp_1',
            'type': 'component',
            'impact': [],
            'depends': []
        }
        e_2 = {
            '_id': 'conn_1',
            'type': 'connector',
            'impact': [],
            'depends': []
        }
        process.update_impact_links(e_1, e_2)
        self.assertTrue(e_2['_id'] in e_1['impact'])
        process.update_impact_links(e_1, e_2)
        print(e_1)
        self.assertTrue(e_1['impact'] == [e_2['_id']])

    def test_update_case_1(self):
        pass

    def test_update_case_2(self):
        pass

    def test_update_case_3(self):
        entities_t1 = [{'_id': 'comp_1',
                        'type': 'component',
                        'impact': [],
                        'depends': []},
                       {'_id': 'conn_1',
                        'type': 'connector',
                        'impact': [],
                        'depends': []}]
        entities_t2 = [{'_id': 'conn_1', 'type': 'connector'},
                       {'_id': 'comp_1', 'type': 'component'},
                       {'_id': 're_1', 'type': 'resource'}]
        ids = {'re_id': 're_1', 'comp_id': 'comp_1', 'conn_id': 'conn_1'}
        self.assertEquals(process.update_case3(entities_t1, ids), 0)
        self.assertEquals(process.update_case3(entities_t2, ids), 1)

    def test_update_case_5(self):
        pass

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

        def test_determine_presence(ids, data):
            """Determine the case with the list of id ids and the data as a set of ids.
            :param ids: a list of ids
            :parama data: a set of ids
            :return: a tuple with the case number and the ids related.
            """
            cache = set(['comp_1', 're_1', 'conn_1'])
            ids_test1 = {
                'comp_id': 'comp_2',
                're_id': 're_2',
                'conn_id': 'conn_2'}
            ids_test2 = {
                'conn_id': 'conn_1',
                'comp_id': 'comp_2',
                're_id': 're_2'}
            ids_test3 = {
                'conn_id': 'conn_1',
                'comp_id': 'comp_1',
                're_id': 're_2'}
            ids_test4 = {
                'conn_id': 'conn_1',
                'comp_id': 'comp_1',
                're_id': 're_1'}
            ids_test5 = {
                'comp_id': 'comp_1',
                're_id': 're_2',
                'conn_id': 'conn_2'}
            ids_test6 = {
                're_id': 're_1',
                'comp_id': 'comp_1',
                'conn_id': 'conn_2'}
            self.assertEqual(
                process.determine_case(ids_test1, cache),
                (False, False, False))
            self.assertEqual(
                process.determine_case(ids_test2, cache),
                (True, False, False))
            self.assertEqual(
                process.determine_case(ids_test3, cache),
                (True, True, False))
            self.assertEqual(
                process.determine_case(ids_test4, cache),
                (True, True, True))
            self.assertEqual(
                process.determine_case(ids_test5, cache),
                (False, True, False))
            self.assertEqual(
                process.determine_case(ids_test6, cache),
                (False, True, True))
            ids_test1_none = {
                'comp_id': 'comp_2',
                're_id': None,
                'conn_id': 'conn_2'}
            ids_test2_none = {
                'conn_id': 'conn_1',
                'comp_id': 'comp_2',
                're_id': None}
            ids_test3_none = {
                'conn_id': 'conn_1',
                'comp_id': 'comp_1',
                're_id': None}
            ids_test4_none = {
                'conn_id': 'conn_1',
                'comp_id': 'comp_1',
                're_id': None}
            ids_test5_none = {
                'comp_id': 'comp_1',
                're_id': None,
                'conn_id': 'conn_2'}
            ids_test6_none = {
                're_id': None,
                'comp_id': 'comp_1',
                'conn_id': 'conn_2'}
            self.assertEqual(
                process.determine_case(ids_test1_none, cache),
                (False, False, None))
            self.assertEqual(
                process.determine_case(ids_test2_none, cache),
                (True, False, None))
            self.assertEqual(
                process.determine_case(ids_test3_none, cache),
                (True, True, None))
            self.assertEqual(
                process.determine_case(ids_test4_none, cache),
                (True, True, None))
            self.assertEqual(
                process.determine_case(ids_test5_none, cache),
                (False, True, None))
            self.assertEqual(
                process.determine_case(ids_test6_none, cache),
                (False, True, None))


if __name__ == '__main__':
    main()
