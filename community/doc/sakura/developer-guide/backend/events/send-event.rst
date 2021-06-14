.. _dev-backend-event-bash:

Send event command
==================

With the ``send_event`` command you can send events (via webservice) to Canopsis.

Authentication key
------------------

Firstly, you must choose a dedicated user for your script.

Secondly, get user's authentication key in web UI :

 * ``Your Account -> Profile -> Authentication key``

Examples
--------

Now you can use command:

.. code-block:: bash

    # echo '{
    >    "connector_name": "vtom",
    >    "event_type": "log",
    >    "source_type": "resource",
    >    "component": "NOM_de_la_machine",
    >    "resource": "NOM_du_JOB",
    >    "state": 0,
    >    "output": "MESSAGE",
    >    "display_name": "DISPLAY_NAME"
    > }' | send_event -s <CANOPSIS>:8082 -a <AUTH_KEY>

**Note:** you can specify ``timestamp`` field for specify unix timestamp
of event.

For more informations about event format you can read: `Event
specification <../../specifications/event.html>`__
