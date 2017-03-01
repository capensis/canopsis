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


if __name__ == '__main__':
    main()
