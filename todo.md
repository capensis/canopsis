Context
=======

TODO
----
Top priority :
--------------
  * **WIP**  sources/python/context_graph/canopsis/context_graph/manager.py (write API)
  * **DONE** sources/python/webcore/canopsis/webcore/services/rest.py
  * **DONE** sources/python/ccalendar/canopsis/ccalendar/process.py
  * **DONE** sources/python/engines/canopsis/engines/event_filter.py
  * **DONE** sources/python/engines/canopsis/engines/eventstore.py
  * **DONE** sources/python/serie/test/manager.py
  * **DONE** sources/python/vevent/canopsis/vevent/process.py
  * **DONE** sources/python/perfdata/canopsis/perfdata/manager.py
  * **DONE** sources/python/perfdata/canopsis/perfdata/process.py

Bottom priority :
-----------------
  * **DONE** sources/python/linklist/canopsis/linklist/ctxpropreg.py
  * **DONE** sources/python/linklist/canopsis/entitylink/manager.py
  * **DONE** sources/python/linklist/canopsis/engines/task_linklist.py
  * **DONE** sources/python/alerts/canopsis/alerts/manager.py
  * **DONE** sources/python/context/test/manager.py
  * **DONE** sources/python/context/scripts/migration_id_name.py
  * **DONE** sources/python/downtime/canopsis/downtime/process.py
  * **DONE** sources/python/downtime/canopsis/downtime/selector.py
  * **DONE** sources/python/downtime/test/process.py
  * **DONE** sources/python/engines/canopsis/engines/context.py
  * **WIP**  sources/python/webcore/canopsis/webcore/services/context.py

Function to replace
-------------------
  * def context(self):
  * def context(self, value):
  * def get_entities(self, ids):
  * def iter_ids(self):
  * def get_entity(
  * def get_entity_by_id(self, _id, _type=None):
  * def get_event(self, entity, event_type='check', **kwargs):
  * def get_by_id(
  * def get(self, _type, names, context=None, extended=False):
  * def find(
  * def put(
  * def remove(
  * def get_entity_id(self, entity):
  * def get_entity_id_context_name(self, entity):
  * def unify_entities(self, entities, extended=False, cache=False):
  * def _configure(self, unified_conf, *args, **kwargs):


Function not used or only used inside the manager
-------------------------------------------------
  * def get_entity_context_and_name(self, entity):
  * def get_children(self, entity):
  * def clean(self, entity):
  * def accept_event_types(self):
  * def accept_event_types(self, value):
  * def get_name(self, entity_id, _type=None):

ContextGraph Unit test
----------------------
  * **SKIP** def add_comp(self, comp):
  * **SKIP** def add_re(self, re):
  * **SKIP** def add_conn(self, conn):
  * **DONE** def get_id(cls, event):
  * **DONE** def check_comp(self, comp_id):
  * **DONE** def get_event(self, entity, event_type='check', \**kwargs):
  * **DONE** def check_re(self, re_id):
  * **DONE** def check_conn(self, conn_id):
  * **DONE** def get_entities_by_id(self, id):
  * **DONE** def put_entities(self, entities):
  * **DONE** def get_all_entities_id(self):
  * **DONE** def check_links(self, conn_id, comp_id, re_id):
  * **DONE** def __update_dependancies(self, from_entity, status, dependancy_type):
  * **DONE** def update_entity(self, entity):
  * **DONE** def delete_entity(self, id_):
  * **DONE** def create_entity(self, entity):
  * **SKIP** def manage_comp_to_re_link(self, re_id, comp_id):
  * **SKIP** def manage_re_to_conn_link(self, conn_id, re_id):
  * **SKIP** def manage_comp_to_conn_link(self, conn_id, comp_id):
  * **SKIP** def _check_conn_comp_link(self, conn_id, comp_id):
  * **SKIP** def _check_conn_re_link(self, conn_id, re_id):
  * **SKIP** def _check_comp_re_link(self, comp_id, re_id):
  * def get_entities(self,
