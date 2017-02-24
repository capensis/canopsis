from unittest import main, TestCase
import canopsis.context_graph.process

class Test(TestCase):

    def create_event(res, comp, conn_name, conn):
        pass

    def test_update_case_6(self):
        entities_t1 = [{'_id': 'comp_1','type': 'component', 'impact': [], 'depends': []}, {'_id': 're_1','type':
                                                                                            'resource', 'impact':[],
                                                                                            'depends':[]}]
        entities_t2 = [{'_id': 'conn_1','type': 'connector'}, {'_id': 'comp_1','type': 'component'}, {'_id': 're_1','type': 'resource'}]
        ids = {'re_id': 're_1', 'comp_id': 'comp_1', 'conn_id': 'conn_1'}
        self.assertEquals(process.update_case_6(entities_t2, ids), 1)
        self.assertEquals(process.update_case_6(entities_t2, ids), 0)


if __name__ == '__main__':
    main()
