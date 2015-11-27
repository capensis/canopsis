#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

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
