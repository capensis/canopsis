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
  * sources/python/engines/canopsis/engines/eventstore.py
  * sources/python/serie/test/manager.py
  * sources/python/vevent/canopsis/vevent/process.py
  * sources/python/perfdata/canopsis/perfdata/process.py
  * sources/python/linklist/canopsis/linklist/ctxpropreg.py

  * sources/python/alerts/canopsis/alerts/manager.py
  * sources/python/webcore/canopsis/webcore/services/context.py
  * sources/python/engines/canopsis/engines/context.py
  * sources/python/downtime/canopsis/downtime/process.py
  * sources/python/downtime/canopsis/downtime/selector.py
  * sources/python/downtime/test/process.py
  * sources/python/perfdata/canopsis/perfdata/manager.py
  * sources/python/linklist/canopsis/entitylink/manager.py
  * sources/python/linklist/canopsis/engines/task_linklist.py
  * sources/python/linklist/canopsis/linklist/ctxpropreg.py
  * sources/python/context/scripts/migration_id_name.py
  * sources/python/context/test/manager.py
  * sources/python/alerts/canopsis/alerts/manager.py

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
