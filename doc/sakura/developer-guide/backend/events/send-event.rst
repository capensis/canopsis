Send\_event
===========

With ``send_event`` command you can send event (via webservice) to
Canopsis.

Authentication key
-------------------

Firstly, you must add one user in ``CPS_event_admin`` group.

Secondly, get user's ``Authentication key`` in web UI
(``Build -> Accounts -> right-click -> Authentication key``)

Get command
-----------

Two way, use directly in canopsis environment or download on github
project
`send\_event <https://raw.github.com/capensis/canopsis/freeze/sources/canotools/bin/send_event>`__.

(For
`window <https://github.com/capensis/canopsis/tree/freeze/sources/extra/powershell>`__)

Examples
--------

Now you can use command:

::

    echo '{
      "connector_name": "vtom",
      "event_type": "log",
      "source_type": "resource",
      "component": "NOM_de_la_machine",
      "resource": "NOM_du_JOB",
      "state": 0,
      "output": "MESSAGE",
      "display_name": "DISPLAY_NAME"
    }
    ' | send_event -s <CANOPSIS>:8082 -a <AUTH_KEY>

**Note:** you can specify ``timestamp`` field for specify unix timestamp
of event.

For more informations about event format you can read: `Event
specification <https://github.com/capensis/canopsis/wiki/Event-specification>`__

