from __future__ import unicode_literals

from unittest import main, TestCase
import canopsis.context_graph.process as process

from canopsis.context_graph.manager import ContextGraph

context_graph_manager = ContextGraph()


class Logger(object):

    def debug(self, log):
        print("DEBUG : {0}".format(log))

    def info(self, log):
        print("INFO : {0}".format(log))

    def warning(self, log):
        print("WARNING : {0}".format(log))

    def critical(self, log):
        print("CRITICAL : {0}".format(log))

    def error(self, log):
        print("ERROR : {0}".format(log))


def create_event(conn, conn_name,  comp=None, res=None):
    event = {"connector": conn, "connector_name": conn_name}
    if comp is not None:
        event["component"] = comp

    if res is not None:
        event["resource"] = res
    return event


def prepare_test_update_context(result):
    conn = None
    comp = None
    re = None

    for entity in result:
        entity["impact"] = sorted(entity["impact"])
        entity["depends"] = sorted(entity["depends"])

    for entity in result:
        if entity["type"] == "connector":
            conn = entity
        elif entity["type"] == "component":
            comp = entity
        elif entity["type"] == "resource":
            re = entity
    return conn, comp, re


class Test(TestCase):

    def setUp(self):
        setattr(process, 'LOGGER', Logger())

    def tearDown(self):
        process.cache.clear()

    def test_create_entity(self):
        id = "id_1"
        name = "name_1"
        etype = "entity type"
        depends = ["id_2", "id_3", "id_4", "id_5"]
        impacts = ["id_6", "id_7", "id_8", "id_9"]
        measurements = {"tag_1": "data_1", "tag_2": "data_2"}
        infos = {"info_1": "foo_1", "info_2": "bar_2"}

        ent = process.create_entity(id, name, etype, depends,
                                    impacts, measurements, infos)

        self.assertEqual(id, ent["_id"])
        self.assertEqual(name, ent["name"])
        self.assertEqual(etype, ent["type"])
        self.assertEqual(depends, ent["depends"])
        self.assertEqual(impacts, ent["impact"])
        self.assertEqual(measurements, ent["measurements"])
        self.assertEqual(infos, ent["infos"])

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
        self.assertDictEqual(re_res_1, {
                '_id': 're_1',
                'name': 're_1',
                'type': 'resource',
                'impact': ['comp_1'],
                'depends': ['conn_1']})
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
        #self.assertEquals(process.update_case3(entities_t1, ids), 0)
        #self.assertEquals(process.update_case3(entities_t2, ids), 1)

    def test_update_case_5(self):
        pass

    def test_update_case_6(self):
        entities_t1 = [{'_id': 'comp_1',
                        'type': 'component',
        self.assertDictEqual(re_res_1, {
                '_id': 're_1',
                'name': 're_1',
                'type': 'resource',
                'impact': ['comp_1'],
                'depends': ['conn_1']})
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
        #self.assertEquals(process.update_case6(entities_t1, ids), 1)
        #self.assertEquals(process.update_case6(entities_t2, ids), 0)

    def test_determine_presence(self):
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
            process.determine_presence(ids_test1, cache),
            (False, False, False))
        self.assertEqual(
            process.determine_presence(ids_test2, cache),
            (True, False, False))
        self.assertEqual(
            process.determine_presence(ids_test3, cache),
            (True, True, False))
        self.assertEqual(
            process.determine_presence(ids_test4, cache),
            (True, True, True))
        self.assertEqual(
            process.determine_presence(ids_test5, cache),
            (False, True, False))
        self.assertEqual(
            process.determine_presence(ids_test6, cache),
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
            process.determine_presence(ids_test1_none, cache),
            (False, False, None))
        self.assertEqual(
            process.determine_presence(ids_test2_none, cache),
            (True, False, None))
        self.assertEqual(
            process.determine_presence(ids_test3_none, cache),
            (True, True, None))
        self.assertEqual(
            process.determine_presence(ids_test4_none, cache),
            (True, True, None))
        self.assertEqual(
            process.determine_presence(ids_test5_none, cache),
            (False, True, None))
        self.assertEqual(
            process.determine_presence(ids_test6_none, cache),
            (False, True, None))

    def test_add_missing_ids(self):
        res_id = "re_id"
        comp_id = "comp_id"
        conn_id = "conn_id"

        ids = {"re_id": res_id,
               "comp_id": comp_id,
               "conn_id": conn_id}

        # check function behaviour for the connector
        process.add_missing_ids((True, False, False), ids)
        self.assertNotIn(conn_id, process.cache)
        process.cache.clear()

        process.add_missing_ids((False, False, False), ids)
        self.assertIn(conn_id, process.cache)
        process.cache.clear()

        with self.assertRaises(KeyError):
            process.add_missing_ids((False, True, True), {
                "re_id": res_id, "comp_id": comp_id})
        process.cache.clear()

        # check function behaviour for the component
        process.add_missing_ids((False, True, False), ids)
        self.assertNotIn(comp_id, process.cache)
        process.cache.clear()

        process.add_missing_ids((False, False, False), ids)
        self.assertIn(conn_id, process.cache)
        process.cache.clear()

        with self.assertRaises(KeyError):
            process.add_missing_ids((True, False, True), {
                "conn_id": conn_id, "re_id": res_id})
        process.cache.clear()

        # check function behaviour for the component
        process.add_missing_ids((False, False, True), ids)
        self.assertNotIn(res_id, process.cache)
        process.cache.clear()

        process.add_missing_ids((False, False, False), ids)
        self.assertIn(conn_id, process.cache)
        process.cache.clear()

        with self.assertRaises(KeyError):
            process.add_missing_ids((True, True, False), {
                "conn_id": conn_id, "comp_id": comp_id})
        process.cache.clear()

    def test_gen_ids(self):
        conn_id = "conn_id"
        conn_name = "conn_name"
        comp_id = "comp_id"
        re_id = "re_id"

        event = create_event(conn_id, conn_name, comp_id, re_id)
        event_re_none = create_event(conn_id, conn_name, comp_id, None)

        expected = {"comp_id": comp_id,
                    "conn_id": "{0}/{1}".format(conn_id, conn_name),
                    "re_id": "{0}/{1}".format(re_id, comp_id)}

        expected_re_none = {"comp_id": comp_id,
                            "conn_id": "{0}/{1}".format(conn_id, conn_name),
                            "re_id": None}

        self.assertEqual(process.gen_ids(event), expected)
        self.assertEqual(process.gen_ids(event_re_none), expected_re_none)

    def test_update_context_case1(self):
        conn_id = "conn_id"
        conn_name = "conn_name"
        comp_id = "comp_id"
        re_id = "re_id"

        event = create_event(conn_id, conn_name, comp_id, re_id)
        ids = process.gen_ids(event)

        expected_conn = process.create_entity(ids["conn_id"],
                                              ids["conn_id"],
                                              "connector",
                                              impact=sorted([ids["comp_id"],
                                                             ids["re_id"]]))

        expected_comp = process.create_entity(ids["comp_id"],
                                              ids["comp_id"],
                                              "component",
                                              depends=sorted([ids["conn_id"],
                                                              ids["re_id"]]))

        expected_re = process.create_entity(ids["re_id"],
                                            ids["re_id"],
                                            "resource",
                                            impact=[ids["comp_id"]],
                                            depends=[ids["conn_id"]])

        res = process.update_context_case1(ids)

        result_conn, result_comp, result_re = prepare_test_update_context(res)

        self.assertDictEqual(expected_comp, result_comp)
        self.assertDictEqual(expected_conn, result_conn)
        self.assertDictEqual(expected_re, result_re)

    def test_update_context_case1_re_none(self):
        conn_id = "conn_id"
        conn_name = "conn_name"
        comp_id = "comp_id"

        event = create_event(conn_id, conn_name, comp_id)
        ids = process.gen_ids(event)

        expected_conn = process.create_entity(ids["conn_id"],
                                              ids["conn_id"],
                                              "connector",
                                              impact=[ids["comp_id"]])

        expected_comp = process.create_entity(ids["comp_id"],
                                              ids["comp_id"],
                                              "component",
                                              depends=[ids["conn_id"]])

        res = process.update_context_case1_re_none(ids)
        result_conn, result_comp, result_re = prepare_test_update_context(res)

        self.assertDictEqual(expected_comp, result_comp)
        self.assertDictEqual(expected_conn, result_conn)

    def test_update_context_case2(self):
        conn_id = "conn_id"
        conn_name = "conn_name"
        comp_id = "comp_id"
        re_id = "re_id"

        ids = process.gen_ids(create_event(conn_id, conn_name, comp_id, re_id))

        expected_conn = process.create_entity(ids["conn_id"],
                                              ids["conn_id"],
                                              "connector",
                                              impact=sorted([ids["comp_id"],
                                                             ids["re_id"]]))

        expected_comp = process.create_entity(ids["comp_id"],
                                              ids["comp_id"],
                                              "component",
                                              depends=sorted([ids["conn_id"],
                                                              ids["re_id"]]))

        expected_re = process.create_entity(ids["re_id"],
                                            ids["re_id"],
                                            "resource",
                                            impact=[ids["comp_id"]],
                                            depends=[ids["conn_id"]])

        conn = process.create_entity("{0}/{1}".format(conn_id, conn_name),
                                     "{0}/{1}".format(conn_id, conn_name),
                                     "connector",
                                     impact=[],
                                     depends=[])

        res = process.update_context_case2(ids, [conn])
        result_conn, result_comp, result_re = prepare_test_update_context(res)

        self.assertDictEqual(expected_conn, result_conn)
        self.assertDictEqual(expected_comp, result_comp)
        self.assertDictEqual(expected_re, result_re)

    def test_update_context_case2_re_none(self):
        conn_id = "conn_id"
        conn_name = "conn_name"
        comp_id = "comp_id"

        ids = process.gen_ids(create_event(conn_id, conn_name, comp_id))

        expected_conn = process.create_entity(ids["conn_id"],
                                              ids["conn_id"],
                                              "connector",
                                              impact=sorted([ids["comp_id"],
                                                             ids["re_id"]]))

        expected_comp = process.create_entity(ids["comp_id"],
                                              ids["comp_id"],
                                              "component",
                                              depends=sorted([ids["conn_id"],
                                                              ids["re_id"]]))

        conn = process.create_entity("{0}/{1}".format(conn_id, conn_name),
                                     "{0}/{1}".format(conn_id, conn_name),
                                     "connector",
                                     impact=[],
                                     depends=[])

        res = process.update_context_case2(ids, [conn])
        result_conn, result_comp, result_re = prepare_test_update_context(res)

        self.assertDictEqual(expected_conn, result_conn)
        self.assertDictEqual(expected_comp, result_comp)

    def test_update_context_case3(self):
        conn_id = "conn_id"
        conn_name = "conn_name"
        comp_id = "comp_id"
        re_id = "re_id"

        ids = process.gen_ids(create_event(conn_id, conn_name, comp_id, re_id))

        expected_conn = process.create_entity(ids["conn_id"],
                                              ids["conn_id"],
                                              "connector",
                                              impact=sorted([ids["comp_id"],
                                                             ids["re_id"]]))

        expected_comp = process.create_entity(ids["comp_id"],
                                              ids["comp_id"],
                                              "component",
                                              depends=sorted([ids["conn_id"],
                                                              ids["re_id"]]))

        expected_re = process.create_entity(ids["re_id"],
                                            ids["re_id"],
                                            "resource",
                                            impact=[ids["comp_id"]],
                                            depends=[ids["conn_id"]])

        conn = process.create_entity("{0}/{1}".format(conn_id, conn_name),
                                     "{0}/{1}".format(conn_id, conn_name),
                                     "connector",
                                     impact=[comp_id],
                                     depends=[])

        comp = process.create_entity(comp_id,
                                     comp_id,
                                     "component",
                                     impact=[],
                                     depends=["{0}/{1}".format(conn_id,
                                                               conn_name)])

        res = process.update_context_case3(ids, [conn, comp])
        result_conn, result_comp, result_re = prepare_test_update_context(res)

        self.assertDictEqual(expected_conn, result_conn)
        self.assertDictEqual(expected_comp, result_comp)
        self.assertDictEqual(expected_re, result_re)

    def test_update_context_case6(self):
        ids1 = {
            're_id': None,
            'conn_id': 'conn_1',
            'comp_id': 'comp_1'
        }
        ids2 = {
            're_id': 're_1',
            'conn_id': 'conn_1',
            'comp_id': 'comp_1'
        }
        in_db_1 = [
            {
                '_id': 're_1',
                'name': 're_1',
                'type': 'resource',
                'impact': ['comp_1'],
                'depends': []
            },
            {
                '_id': 'comp_1',
                'name': 'comp_1',
                'type': 'component',
                'impact': [],
                'depends': ['re_1']}]
        in_db_2 = [{
            '_id': 'comp_1',
            'name': 'comp_1',
            'type': 'component',
            'impact': [],
            'depends': []}]
        res_1 = process.update_context_case6(ids1, in_db_1) 
        res_2 = process.update_context_case6(ids2, in_db_2)


        comp_res_1 = None
        conn_res_1 = None
        re_res_1 = None
        comp_res_2 = None
        conn_res_2 = None
        re_res_2 = None
        for i in res_1:
            if i['type'] == 'component':
                comp_res_1= i
            if i['type'] == 'resource':
                re_res_1= i
            if i['type'] == 'connector':
                conn_res_1= i
        for i in res_2:
            if i['type'] == 'component':
                comp_res_2 = i
            if i['type'] == 'resource':
                re_res_2 = i
            if i['type'] == 'connector':
                conn_res_2 = i

        for i in comp_res_1:
            if isinstance(comp_res_1[i], list):
                comp_res_1[i] = sorted(comp_res_1[i])
        for i in conn_res_1:
            if isinstance(conn_res_1[i], list):
                conn_res_1[i] = sorted(conn_res_1[i])
        for i in re_res_1:
            if isinstance(re_res_1[i], list):
                re_res_1[i] = sorted(re_res_1[i])
        for i in comp_res_2:
            if isinstance(comp_res_2[i], list):
                comp_res_2[i] = sorted(comp_res_2[i])
        for i in conn_res_2:
            if isinstance(conn_res_2[i], list):
                conn_res_2[i] = sorted(conn_res_2[i])
        self.assertDictEqual(comp_res_1, {
                '_id': 'comp_1',
                'name': 'comp_1',
                'type': 'component',
            'impact': [],
            'depends': sorted(['re_1', 'conn_1'])})
        print(re_res_1, {
                '_id': 're_1',
                'name': 're_1',
                'type': 'resource',
                'impact': ['comp_1'],
                'depends': ['conn_1']})
        self.assertDictEqual(re_res_1, {
                '_id': 're_1',
                'name': 're_1',
                'type': 'resource',
                'impact': ['comp_1'],
                'depends': ['conn_1']})
        self.assertDictEqual(conn_res_1, {
                '_id': 'conn_1',
                'name': 'conn_1',
                'type': 'connector',
                'impact': sorted(['comp_1', 're_1']),
                'depends': []})

        self.assertDictEqual(comp_res_1, {
                '_id': 'comp_1',
                'name': 'comp_1',
                'type': 'component',
                'impact': [],
                'depends': sorted(['conn_1'])})
        self.assertEqual(res_2, None)
        self.assertDictEqual(conn_res_1, {
                '_id': 'conn_1',
                'name': 'conn_1',
                'type': 'connector',
                'impact': sorted(['comp_1']),
                'depends': []})

if __name__ == '__main__':
    main()
