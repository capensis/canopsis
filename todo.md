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
  * sources/python/alerts/canopsis/alerts/manager.py
  * sources/python/alerts/canopsis/alerts/manager.py
  * sources/python/context/test/manager.py
  * sources/python/context/scripts/migration_id_name.py
  * **DONE** sources/python/downtime/canopsis/downtime/process.py
  * **DONE** sources/python/downtime/canopsis/downtime/selector.py
  * sources/python/downtime/test/process.py
  * sources/python/engines/canopsis/engines/context.py
  * sources/python/webcore/canopsis/webcore/services/context.py

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
