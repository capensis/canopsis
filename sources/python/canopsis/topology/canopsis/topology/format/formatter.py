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

from pymongo import Connection, MongoClient
import json
from bson.json_util import dumps
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
import logging


class Formatter(object):
    """docstring for Formatter"""
    TOPOIDS = ('conns', 'nodes', 'root', 'crecord_name')
    EVENT_TYPE = ('operator', 'check', 'selector', 'topology')
    SOURCE_TYPE = ('resource', 'component')
    TYPE = ('event_type', 'source_type')
    OPERATOR_ID = ('Cluster', 'Worst State', 'And', 'Or', 'Best State')
    STATE = (0, 1, 2, 3, 4)
    CONTEXT = ('type', 'connector', 'connector_name', 'component', 'resource')
    TOPALL = "{'crecord_type':'topology'}"
    TOPOQ = "{{'crecord_type':'topology', 'crecord_name':{!r}}}"
    # The collections name
    NAMESPACES = ('objectv1', 'eventsv1')
    MONGO_PORT = 27017
    MONGO_URL = "localhost"

    def __init__(self, query):
        self.qr = query
        self.data = self.loads()
        # logger is initialized by Formatter
        self.logger = logging.getLogger('formatter')

    @classmethod
    def storage_connection(self, namespace):
        '''
        Get the cconnection to canopsis DataBase.

        :param namespace: specify teh collection name.

        :return: a connection.
        :rtype: Mongo Connection.
        '''
        connection = get_storage(
            namespace=namespace,
            account=Account(
                user="root",
                group="root"
            )
        ).get_backend()
        return connection

    def get_data(self, kind=0, q=None):
        '''
         Access MongoDB and load topology or events data.

         :param kink: specify which request should be running in the DataBase.

         :return: a cursor of topology or events.
         :rtype: Cursor of elements dictionnary or empty dictionnary.
        '''
        if kind == 0:
            query = self.TOPOQ.format(self.qr)
            # Format String
            json_acceptable = query.replace("'", "\"")
            query = json.loads(json_acceptable)
            cursor = self.storage_connection(self.NAMESPACES[0]).find(query)
        else:
            query = q
            json_acceptable = query.replace("'", "\"")
            query = json.loads(json_acceptable)
            cursor = self.storage_connection(self.NAMESPACES[1]).find(query)
        return cursor

    @classmethod
    def get_all_comp_ids(self):
        '''
        '''
        # Format String
        json_acceptable = self.TOPALL.replace("'", "\"")
        query = json.loads(json_acceptable)
        cursor = self.storage_connection(self.NAMESPACES[0]).find(query)
        res = []
        for comp in list(cursor):
            res.append(str(comp.get('crecord_name')))
        return res

    def connection(self, kind=0, q=None):
        '''
         Access MongoDB and load topology or events data.

         :param kink: specify which request should be running in the DataBase.

         :return: a cursor of topology or events.
         :rtype: Cursor of elements dictionnary or NoneType.
        '''
        if self.username is None:
            connection = Connection()
            db = connection['canopsis']
        else:
            connection = MongoClient(self.MONGO_URL, self.MONGO_PORT)
            db = connection["canopsis"]
            # Do the authentication
            db.authenticate(self.username, self.password)
        if kind == 0:
            query = self.qr
            # Format string
            json_acceptable = query.replace("'", "\"")
            query = json.loads(json_acceptable)
            cursor = db.objectv1.find(query)
        else:
            query = q
            # Format string
            json_acceptable = query.replace("'", "\"")
            query = json.loads(json_acceptable)
            cursor = db.eventsv1.find(query)
        connection.close()
        return cursor

    def loads(self, kind=0, q=None):
        '''
         Serealize cursor data into JSON.
         :param kind: specify which request should be running in the DataBase.

         :return: a dictionnary of elements.
         :rtype: dictionnary.
        '''
        root_id = ""
        crecord_name = ""
        _str = self.dump(kind, q)
        comp_json = json.loads(_str)
        if len(comp_json) > 0:
            # catch exception here
            res = comp_json[0]
            # Retreive the root_id and the crecord_name
            root_id = res.get(self.TOPOIDS[2])
            crecord_name = res.get(self.TOPOIDS[3])
            # Change the root_id with the crecord_name
            _str = _str.replace(str(root_id), str(crecord_name))
            # load the new dict
            comp_json = json.loads(_str)
            res = comp_json[0]
        else:
            res = {}
        return res

    def dump(self, kind=0, q=None):
        '''
        '''
        return dumps(self.get_data(kind, q))

    def print_keys(self):
        tdata = self.data
        for k, v in tdata.iteritems():
            self.logger.info(u'Key: {} - Value: {}'.format(k, v))

    def get_value(self, value):
        '''
         Get elements of the topology.

         :return: dictionnary or list of elements.
         :rtype: dictionnary or list of dictionnary.
        '''
        return self.data.get(value)

    def get_comp_graph(self):
        '''
        '''
        return self.data.get(self.TOPOIDS[0])

    def get_components(self):
        '''
        '''
        return self.data.get(self.TOPOIDS[1])

    def get_component_keys(self):
        '''
         Get the list of distict component inside the topology.

         :return: a list of component.
         :rtype: List.
        '''
        return self.data.get(self.TOPOIDS[1]).keys()

    def get_event_type(self, kind=0):
        '''
        '''
        event_comp = {}
        ops_list = []
        chk_list = []
        sel_list = []
        components = {}
        if kind == 0:
            components = self.get_components()
        else:
            components = self.comp_formatter()
        for d in self.get_components().keys():
            if self.iequal(
                    components.get(d).get(self.TYPE[0]),
                    self.EVENT_TYPE[0]
            ):
                tmp_dict = {}
                tmp_dict[d] = components.get(d)
                ops_list.append(tmp_dict)
                event_comp[self.EVENT_TYPE[0]] = ops_list

            if self.iequal(
                    components.get(d).get(self.TYPE[0]),
                    self.EVENT_TYPE[1]
            ):
                tmp_dict = {}
                tmp_dict[d] = components.get(d)
                chk_list.append(tmp_dict)
                event_comp[self.EVENT_TYPE[1]] = chk_list

            if self.iequal(
                    components.get(d).get(self.TYPE[0]),
                    self.EVENT_TYPE[2]
            ):
                tmp_dict = {}
                tmp_dict[d] = components.get(d)
                sel_list.append(tmp_dict)
                event_comp[self.EVENT_TYPE[2]] = sel_list

            if self.iequal(
                    components.get(d).get(self.TYPE[0]),
                    self.EVENT_TYPE[3]
            ):
                tmp_dict = {}
                tmp_dict[d] = components.get(d)
                sel_list.append(tmp_dict)
                event_comp[self.EVENT_TYPE[3]] = sel_list

        return event_comp

    def get_source_type(self):
        '''
         Retreive all components classify by source.

         :return: a dictionnary of source classify by type.
         :rtype: dictionnary.
        '''
        source_type = {}
        resr_list = []
        comp_list = []
        for d in self.get_components().keys():
            if self.iequal(
                    self.get_components().get(d).get(self.TYPE[1]),
                    self.SOURCE_TYPE[0]
            ):
                resr_list.append(self.get_components().get(d))
                source_type[self.SOURCE_TYPE[0]] = resr_list

            if self.iequal(
                    self.get_components().get(d).get(self.TYPE[1]),
                    self.SOURCE_TYPE[1]
            ):
                comp_list.append(self.get_components().get(d))
                source_type[self.SOURCE_TYPE[1]] = comp_list

        return source_type

    def match_operator(self, comps):
        '''
        '''
        clt_list = []
        wst_list = []
        and_list = []
        or_list = []
        bes_list = []
        operators = {}
        components = comps
        for comp in components.get(self.EVENT_TYPE[0]):
            if self.iequal(
                    comp.values()[0].get('component'),
                    self.OPERATOR_ID[0]
            ):
                clt_list.append(comp)
                operators[self.OPERATOR_ID[0]] = clt_list

            if self.iequal(
                    comp.values()[0].get('component'),
                    self.OPERATOR_ID[1]
            ):
                wst_list.append(comp)
                operators[self.OPERATOR_ID[1]] = wst_list

            if self.iequal(
                    comp.values()[0].get('component'),
                    self.OPERATOR_ID[2]
            ):
                and_list.append(comp)
                operators[self.OPERATOR_ID[2]] = and_list

            if self.iequal(
                    comp.values()[0].get('component'),
                    self.OPERATOR_ID[3]
            ):
                or_list.append(comp)
                operators[self.OPERATOR_ID[3]] = or_list

            if self.iequal(
                    comp.values()[0].get('component'),
                    self.OPERATOR_ID[4]
            ):
                bes_list.append(comp)
                operators[self.OPERATOR_ID[4]] = bes_list

        return operators

    def get_operators(self):
        '''
         Classify components by operator.

         :return: a list of components for this kind.
         :rtype: list.
        '''
        return self.get_event_type().keys()

    def operator_components(self):
        '''
        '''
        return self.get_event_type().get(self.EVENT_TYPE[0])

    def check_components(self):
        '''
        '''
        return self.get_event_type().get(self.EVENT_TYPE[1])

    def selector_components(self):
        '''
        '''
        return self.get_event_type().get(self.EVENT_TYPE[2])

    def topology_components(self):
        '''
        '''
        return self.get_event_type().get(self.EVENT_TYPE[3])

    def is_context_compatible(self, elt):
        '''
        Verify if component has all context variables.
        '''
        comp = self.get_components().get(elt)
        for ctx in self.CONTEXT:
            if comp.get(ctx) is None:
                return False
        return True

    def get_connector_name(self):
        '''
        Loads the context data from events collections.
        '''
        return self.loads(1)

    def diff(self, newList):
        '''
        return the difference between two lists.
        '''
        return list(set([i for i in self.CONTEXT]) - set(newList))

    def query_generator(self, comp=None):
        '''
        '''
        start = "{"
        end = "}"
        data = ""
        quote = "'"
        comma = ","
        separator = ":"
        missing_ctx = []

        if comp is None:
            top = self.get_event_type().get(self.EVENT_TYPE[3])[0]
        else:
            top = self.get_components().get(comp)

        if top.get(self.CONTEXT[0]) is not None:
            missing_ctx.append(self.CONTEXT[0])
            data += '{0}{1}{0}{2}{0}{3}{0}{4}'.format(
                quote,
                self.CONTEXT[0],
                separator,
                top.get(self.CONTEXT[0]),
                comma
            )

        if top.get(self.CONTEXT[1]) is not None:
            missing_ctx.append(self.CONTEXT[1])
            data += '{0}{1}{0}{2}{0}{3}{0}{4}'.format(
                quote,
                self.CONTEXT[1],
                separator,
                top.get(self.CONTEXT[1]),
                comma
            )

        if top.get(self.CONTEXT[2]) is not None:
            missing_ctx.append(self.CONTEXT[2])
            data += '{0}{1}{0}{2}{0}{3}{0}{4}'.format(
                quote,
                self.CONTEXT[2],
                separator,
                top.get(self.CONTEXT[2]),
                comma
            )

        if top.get(self.CONTEXT[3]) is not None:
            missing_ctx.append(self.CONTEXT[3])
            data += '{0}{1}{0}{2}{0}{3}{0}{4}'.format(
                quote,
                self.CONTEXT[3],
                separator,
                top.get(self.CONTEXT[3]),
                comma
            )

        if top.get(self.CONTEXT[4]) is not None:
            missing_ctx.append(self.CONTEXT[4])
            data += '{0}{1}{0}{2}{0}{3}{0}{4}'.format(
                quote,
                self.CONTEXT[4],
                separator,
                top.get(self.CONTEXT[4]),
                comma
            )

        query = start + data.rstrip(comma) + end
        return query, self.diff(missing_ctx)

    def comp_formatter(self):
        '''
        '''
        if self.get_components() is not None:
            comps = self.get_components().copy()
        else:
            self.logger.info(
                u'Topology {} does not exist in database...'.format(self.qr))
            comps = {}
        for c in self.get_component_keys():
            q, lst = self.query_generator(c)
            # Loads the context information
            res = self.loads(1, q)
            if len(res) != 0:
                for d in lst:
                    if self.iequal(d, 'type'):
                        comps.get(c)[unicode(d)] = unicode(
                            comps.get(c).get(self.TYPE[0]))

                    else:
                        comps.get(c)[unicode(d)] = res.get(unicode(d))

            if len(res) == 0:
                self.logger.info(u'Component whithout event data {}'.format(c))
        return comps

    def get_root(self):
        '''
        Get the topology root.
        '''
        root_id = self.get_root_id()
        return self.data.get(self.TOPOIDS[1]).get(root_id)

    def get_root_id(self):
        '''
        Get the components root ID.
        '''
        return self.data.get(self.TOPOIDS[2])

    def get_operator_data(self, op_dict):
        '''
        Get Component form items.
        '''
        return op_dict.values().get('form').get('items')

    def iequal(self, a, b):
        '''
        Verify equality with case sensitive.
        :param a : first element.
        :param b: second element.

        :return: a boolean.
        :rtype: Boolean.
        '''
        try:
            return a.lower() == b.lower()
        except AttributeError as e:
            return a == b


if __name__ == '__main__':
    pass
