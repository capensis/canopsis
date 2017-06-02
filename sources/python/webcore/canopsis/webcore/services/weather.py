# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

from canopsis.common.ws import route
from canopsis.context_graph.manager import ContextGraph
from canopsis.alerts.reader import AlertsReader

context_manager = ContextGraph()
alarm_manager = AlertsReader()


def exports(ws):

    @route(
        ws.application.get,
        name='weather/get/selector',
        payload=[
            'selector_filter'
        ]
    )
    def get_selector(
            selector_filter
    ):
        selector_filter['type'] = 'selector'
        selector_list = context_manager.get_entities(query=selector_filter)
        
        ret_val = []
        for i in selector_list:
            tmp = {}
            f = open('/home/tgosselin/fichierdelog', 'a')
            f.write('{0}'.format(alarm_manager.get(filter_={'d':i['_id']})))
            f.close()
            tmp_alarm = alarm_manager.get(filter_={'d' : i['_id']})['alarms']

            tmp['entity_id'] = i['_id']
            try:
                tmp['criticity'] = i['infos']['criticity']
            except:
                tmp['criticity'] = ''
            try:
                tmp['org'] = i['infos']['org']
            except:
                tmp['org'] = ''
            tmp['sla_text'] = 'ca arrive avec les sla'
            tmp['display_name'] = i['name']
            if tmp_alarm != []:
                tmp['state'] = tmp_alarm[0]['v']['state']
                tmp['status'] = tmp_alarm[0]['v']['status']
                tmp['snooze'] = tmp_alarm[0]['v']['snooze']
                tmp['ack'] = tmp_alarm[0]['v']['ack']
            tmp['pbehavior'] = []
            tmp['linklist'] = []
            ret_val.append(tmp)
        return ret_val
