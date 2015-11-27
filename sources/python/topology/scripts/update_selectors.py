from canopsis.context.manager import Context
from canopsis.topology.manager import TopologyManager

context = Context()
tm = TopologyManager()

storage_name = 'ctx_storage'


def update_selectors():
    selectors = context.find(_type='selector')
    selector_ids_to_remove = []
    to_remove = '/selector/selector'
    to_put = '/selector/canopsis'
    to_remove_len = len(to_remove)
    # for all context selectors
    for selector in selectors:
        # update connector
        selector['connector'] = 'canopsis'
        # get old id
        old_id = selector['_id']
        # add old id in list to remove
        selector_ids_to_remove.append(old_id)
        # get new id
        _id = '{0}{1}'.format(to_put, old_id[to_remove_len:])
        # remove old id from selector
        context[storage_name].put_element(_id=_id, element=selector)
        # get selector topology nodes
        node_id = '/selector/canopsis/engine/{0}'.format(
            selector[Context.NAME].replace('_', ' ')
        )
        elts = tm.get_elts(info={'entity': node_id})
        for elt in elts:
            elt.info['entity'] = _id
            elt.save(manager=tm)
    # remove all old selectors
    context[storage_name].remove_elements(ids=selector_ids_to_remove)

if __name__ == '__main__':
    update_selectors()
