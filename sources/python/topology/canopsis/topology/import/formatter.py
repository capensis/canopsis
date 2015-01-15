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

from pymongo import Connection
import json
from bson.json_util import dumps


class Formatter(object):
   '''
      TODO
   '''
   TOPOIDS = ('conns', 'nodes')
   EVENT_TYPE = ('operator', 'check', 'selector', 'topology')
   SOURCE_TYPE = ('resource', 'component')
   TYPE = ('event_type', 'source_type')
   OPERATOR_ID = ('Cluster', 'Worst State', 'And', 'Or', 'Best State')
   STATE = (0, 1, 2, 3, 4)
   CONTEXT = ('type', 'connector', 'connector_name', 'component', 'resource')
   TOPOQ = "{'crecord_type':'topology', 'crecord_name':'internal'}" # Topology query



   def __init__(self, datasource=None):
      self.arg = datasource
      self.cursor = self.connection()
      self.data = self.loads()

   def connection(self, kind=0):
      '''
         Access MongoDB and load topology or events data.

         :param kink: specify which request should be running in the DataBase.

         :return: a cursor of topology or events.
         :rtype: Cursor of elements dictionnary or NoneType.
      '''
      connection = Connection()
      db = connection.canopsis
      if(kind==0):
         query = self.TOPOQ
         json_acceptable = query.replace("'", "\"") # Format string
         query = json.loads(json_acceptable)
         cursor = db.object.find(query)
      else:
         query = self.query_generator()[0]
         json_acceptable = query.replace("'", "\"") # Format string
         query = json.loads(json_acceptable)
         cursor = db.events.findOne(query)
      connection.close()
      return cursor

   def loads(self, kind=0):
      '''
         Serealize cursor data into JSON.
         :param kind: specify which request should be running in the DataBase.

         :return: a dictionnary of elements.
         :rtype: dictionnary.
      '''
      str = dumps(self.connection(kind))
      if len(json.loads(str)) > 0:
         res = json.loads(str)[0] # catch exception here
      else:
         res = {}
      return res

   def print_key(self):
      tdata = self.data
      for k, v in tdata.iteritems():
         print k, v

   def get_value(self, value):
      '''
         Get elements of the topology.

         :return: dictionnary or list of elements.
         :rtype: dictionnary or list of dictionnary.
      '''
      return self.data.get(value)

   def get_comp_graph(self):
      ''' 
         TODO
      '''
      return self.data.get(self.TOPOIDS[0])

   def get_value(self, value):
      '''
         TODO
      '''
      return self.data.get(value)

   def get_components(self):
      '''
         Retreives all components from object collections.

         :return: the list of components inside the topology.
         :rtype: dictionnary.
      '''
      return self.data.get(self.TOPOIDS[1])

   def get_comp_graph(self):
      """ TODO """
      return self.data.get(self.TOPOIDS[0])

   def get_component_keys(self):
      '''
         Get the list of distict component inside the topology.

         :return: a list of component.
         :rtype: List.
      '''
      return self.data.get(self.TOPOIDS[1]).keys()

   def get_event_type(self):
      '''
         Retreive all components classify by event.

         :return: a dictionnary of events classify their type.
         :rtype: dictionnary.
      '''
      event_comp = {}
      ops_list = []
      chk_list = []
      sel_list = []
      top_list = []
      for d in self.get_components().keys():
         if self.get_components().get(d).get(self.TYPE[0]) == self.EVENT_TYPE[0]:
            ops_list.append(self.get_components().get(d))
            event_comp[self.EVENT_TYPE[0]] = ops_list
         if self.get_components().get(d).get(self.TYPE[0]) == self.EVENT_TYPE[1]:
            chk_list.append(self.get_components().get(d))
            event_comp[self.EVENT_TYPE[1]] = chk_list
         if self.get_components().get(d).get(self.TYPE[0]) == self.EVENT_TYPE[2]:
            sel_list.append(self.get_components().get(d))
            event_comp[self.EVENT_TYPE[2]] = sel_list
         if self.get_components().get(d).get(self.TYPE[0]) == self.EVENT_TYPE[3]:
            top_list.append(self.get_components().get(d))
            event_comp[self.EVENT_TYPE[3]] = top_list
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
         if self.get_components().get(d).get(self.TYPE[1]) == self.SOURCE_TYPE[0]:
            resr_list.append(self.get_components().get(d))
            source_type[self.SOURCE_TYPE[0]] = resr_list
         if self.get_components().get(d).get(self.TYPE[1]) == self.SOURCE_TYPE[1]:
            comp_list.append(self.get_components().get(d))
            source_type[self.SOURCE_TYPE[1]] = comp_list
      return source_type

   def match_operator(self):
      '''
         TODO
      '''
      for comp in self.get_event_type().get(self.EVENT_TYPE[0]):
         if comp.get('label') == self.OPERATOR_ID[0] :
            print type(comp.get('form').get('items'))
         if comp.get('label') == self.OPERATOR_ID[1]:
            pass
         if comp.get('label') == self.OPERATOR_ID[2]:
            comp.get('form').get('items')
         if comp.get('label') == self.OPERATOR_ID[3]:
            comp.get('form').get('items')
         if comp.get('label') == self.OPERATOR_ID[4]:
            comp.get('form').get('items')

   def get_operators(self):
      '''
         Classify components by operator.

         :return: a list of components for this kind.
         :rtype: list.
      '''
      self.get_event_type().keys()

   def operator_components(self):
      '''TODO'''
      return self.get_event_type().get(self.EVENT_TYPE[0])

   def check_components(self):
      """ TODO """
      return self.get_event_type().get(self.EVENT_TYPE[1])

   def selector_components(self):
      """ TODO """
      return self.get_event_type().get(self.EVENT_TYPE[2])

   def topology_components(self):
      """ TODO """
      return self.get_event_type().get(self.EVENT_TYPE[3])

   def nodes_factory(self):
      """ TODO """
      nod_list = []
      for comp in self.get_component_keys():
         nod_list.append(Node(comp))
      return nod_list

   def get_node_instance(self,node):
      """ TODO """
      isnt = None
      for nondei in self.nodes_factory():
         if nondei.data == node:
            isnt = nondei
            break
      return isnt

   def is_context_compatible(self, elt):
      '''
        Verify if component has all context variables
      '''
      comp = self.get_components().get(elt)
      for ctx in self.CONTEXT:
         if comp.get(ctx) == None:
            return False
      return True

    def topology_format(self):
      '''
         TODO
      '''
      for data in self.get_event_type().get(self.EVENT_TYPE[3]):
         print data

   def get_connector_name(self):
      '''
         Loads the context data from events collections
      '''
      return self.loads(1)

   def diff(self, newList):
      '''
         return the diff between two lists
      '''
      return list(set([i for i in self.CONTEXT]) - set(newList))

   def query_generator(self):
      '''
         Generates query which will be executed
      '''
      start = "{"
      end = "}"
      data = ""
      quote = "'"
      comma = ","
      separator = ":"
      missing_ctx = []
      top = self.get_event_type().get(self.EVENT_TYPE[3])[0]
      if top.get(self.CONTEXT[0]) != None:
         missing_ctx.append(self.CONTEXT[0])
         data += quote + self.CONTEXT[0] + quote + separator + quote + top.get(self.CONTEXT[0]) + quote + comma
      if top.get(self.CONTEXT[1]) != None:
         missing_ctx.append(self.CONTEXT[1])
         data += quote + self.CONTEXT[1] + quote + separator + quote + top.get(self.CONTEXT[1]) + quote + comma
      if top.get(self.CONTEXT[2]) != None:
         missing_ctx.append(self.CONTEXT[2])
         data += quote + self.CONTEXT[2] + quote + separator + quote + top.get(self.CONTEXT[2]) + quote + comma
      if top.get(self.CONTEXT[3]) != None:
         missing_ctx.append(self.CONTEXT[3])
         data += quote + self.CONTEXT[3] + quote + separator + quote + top.get(self.CONTEXT[3]) + quote + comma
      query = start + data.rstrip(comma) + end
      return query, self.diff(missing_ctx)

   def query_generator(self, comp):
      '''
         Generates query which will be executed from a component
      '''
      start = "{"
      end = "}"
      data = ""
      quote = "'"
      comma = ","
      separator = ":"
      missing_ctx = []
      top = self.get_components().get(comp)
      if top.get(self.CONTEXT[0]) != None:
         missing_ctx.append(self.CONTEXT[0])
         data += quote + self.CONTEXT[0] + quote + separator + quote + top.get(self.CONTEXT[0]) + quote + comma
      if top.get(self.CONTEXT[1]) != None:
         missing_ctx.append(self.CONTEXT[1])
         data += quote + self.CONTEXT[1] + quote + separator + quote + top.get(self.CONTEXT[1]) + quote + comma
      if top.get(self.CONTEXT[2]) != None:
         missing_ctx.append(self.CONTEXT[2])
         data += quote + self.CONTEXT[2] + quote + separator + quote + top.get(self.CONTEXT[2]) + quote + comma
      if top.get(self.CONTEXT[3]) != None:
         missing_ctx.append(self.CONTEXT[3])
         data += quote + self.CONTEXT[3] + quote + separator + quote + top.get(self.CONTEXT[3]) + quote + comma
      query = start + data.rstrip(comma) + end
      return query, missing_ctx

   def comp_formatter(self):
      '''
         TODO
      '''
      comps = self.get_components()
      for c in self.get_component_keys():
         q, lst = self.query_generator(c)
         # Loads the context information
         res = self.loads(1)
         if(res == None):
            return
         for d in lst:
            if(d == 'type'):
               comps.get(c)[unicode(d)] = unicode(comps.get(c).get(self.TYPE[0]))
            else:
               if(res!=None):
                  comps.get(c)[unicode(d)] = res.get(unicode(d))
      return comps
