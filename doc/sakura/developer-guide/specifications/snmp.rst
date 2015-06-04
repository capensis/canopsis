===========================================
Snmp: Engine and tools to process SNMP trap
===========================================

.. module:: canopsis.snmp

Description
===========

Functional
----------

The SNMP contain an engine that analyse the trap sent by `snmp2canopsis`,
search a corresponding rule to the trap and generate a new canopsis compliant
event.

Trap definitions
----------------

`snmp2canopsis` already preprocess the snmp trap to send a semi-complete event
that contains theses fields:

- `snmp_trap_oid`: OID of the SNMP trap
- `snmp_version`: version of SNMP used for the trap ("1" or "2c")
- `snmp_vars`: OID/values used in the trap
- `snmp_timeticks` for SNMP v1

Rule definitions
----------------

The rules are managed by the :class:`~canopsis.snmp.manager.SnmpManager`.
A rule is defined as

- a OID as an identifier
- a required `state` field
- a required `resource` field
- a optional `component` field
- a optional `output` field


Processing
----------

When a trap is received, the Snmp engine will search a rule according to the
`snmp_trap_oid` field.

If no rule are linked to the oid:
- the original event will have a new field named `snmp_trap_match` set to False

If a rule is find:
- Translated event will have a new field named `snmp_trap_match` set to True
- Translating event consists in compiling rules templates for rules fields. the template compilation context is set with both mib objects oids and event snmp_vars values (see section below: mib object translation)
- Translated event is sent to the normal canopsis event processing if translation succeeded
During the processing on the fields, if any errors happen:
- the original event will have a new field named `snmp_trap_errors`, containing a list of the errors.

Mib object translation
----------------------

First, rule fields may contain an Mib module object reference, for instance the following template works for a Nagios module and objects definitions

.. code-block:: javascript

   //A rule sample
   {
      ...
      component: 'customcomponent_{{ nSvcEvent }}'
      ...
   }

As the rule's module is ``NAGIOS-NOTIFY-MIB``, the UI editor shows all names availables in mib collection where ``nodetype = 'notification'``.
Selecting a name for a module let the UI find all the selected mib document **objects**. These objects are then available to edit templates.

By knowing in the rule witch module and witch name are used, it is possible to build the template compilation context. When the engine meets a rule for the current trap, it is able to get the mib information and objects related to this rule. When object list is retrieved, the engine search for object oids by building a document id like ``modulename::objectname`` where in database a document shoud live and be like:

.. code-block:: javascript

   //A rule sample
   {
      _id: 'module_name::object_name',
      'oid': 'x.y.z'
   }

When the object **oid** is retrieved, the template context is set with the the following information:

.. code-block:: python

   # Arbitrary values
   object_name = 'object_name'
   object_oid = 'x.y.z'

   # Template context building
   template_context_value = event['snmp_vars'][object_oid]
   template_context[object_name] = template_context_value

This way, the template context should looks like in our case

.. code-block:: javascript

   {
      'nSvcEvent': 'componentinfo',
      ...

   }

Then the translated event will have a component value equal to **customcomponent_componentinfo**

Mibs
----

The module :mod:`~canopsis.snmp.mibs` goal is to store and query SNMP
notifications and objects. It requires the `smitools` package on your system,
or the binary `smidump`.

Please note this database is NOT a complete MIBS database, it contains only the
necessary informations for the UI.

It's all about key/value again, as notifications and objects have key.

Import new MIB
~~~~~~~~~~~~~~

You can import new MIB from your system via command line::

    python -m canopsis.snmp.mibs -k /usr/share/mibs/ietf/*

A long output will be printed, and the end of the output should look like
this::

    Import summary
    - 415 notifications definitions
    - 22167 objects definitions
    - 1 error
      - /usr/share/mibs/ietf/ISIS-MIB: Invalid python generated from smidump

Pragmatically, you can do it with the MibsManager::

    from canopsis.snmp.mibs import MibsManager
    manager = MibsManager()
    manager.import_mibs("yourfilename")


Query the MIBS database
~~~~~~~~~~~~~~~~~~~~~~~

.. note::

    Except for MIB, all name search if found returns a oid.

Get the description of a MIB::

    $ python -m canopsis.snmp.mibs --query IF-MIB
    {u'_id': u'IF-MIB',
     u'contact': u'   Keith McCloghrie\nCisco Systems, Inc.\n170 West Tasman Drive\nSan Jose, CA  95134-1706\nUS\n\n408-526-5260\nkzm@cisco.com',
     u'description': u"The MIB module to describe generic objects for network\ninterface sub-layers.  This MIB is an updated version of\nMIB-II's ifTable, and incorporates the extensions defined in\nRFC 1229.",
     u'identity node': u'ifMIB',
     u'language': u'SMIv2',
     u'nodetype': u'module',
     u'organization': u'IETF Interfaces MIB Working Group',
     u'revisions': [{u'date': u'2000-06-14 00:00',
                     u'description': u'Clarifications agreed upon by the Interfaces MIB WG, and\npublished as RFC 2863.'},
                    {u'date': u'1996-02-28 21:55',
                     u'description': u'Revisions made by the Interfaces MIB WG, and published in\nRFC 2233.'},
                    {u'date': u'1993-11-08 21:55',
                     u'description': u'Initial revision, published as part of RFC 1573.'}]}

Get the OID of the linkDown notification::

    $ python -m canopsis.snmp.mibs --query IF-MIB::linkDown
    {u'_id': u'IF-MIB::linkDown', u'oid': u'1.3.6.1.6.3.1.1.5.3'}

Get the description of the linkDown OID::

    $ python -m canopsis.snmp.mibs --query 1.3.6.1.6.3.1.1.5.3
    {u'_id': u'1.3.6.1.6.3.1.1.5.3',
     u'description': u'A linkDown trap signifies that the SNMP entity, acting in\nan agent role, has detected that the ifOperStatus object for\none of its communication links is about to enter the down\nstate from some other state (but not from the notPresent\nstate).  This other state is indicated by the included value\nof ifOperStatus.',
     u'moduleName': u'IF-MIB',
     u'name': u'linkDown',
     u'nodetype': u'notification',
     u'objects': {u'ifAdminStatus': {u'module': u'IF-MIB',
                                     u'nodetype': u'object'},
                  u'ifIndex': {u'module': u'IF-MIB', u'nodetype': u'object'},
                  u'ifOperStatus': {u'module': u'IF-MIB', u'nodetype': u'object'}},
     u'oid': u'1.3.6.1.6.3.1.1.5.3',
     u'status': u'current'}

From an OID, we know which MIB is associated via `moduleName`, the `nodetype`
is a notification, so there is an additional field named `name`.
To get the informations of all the vars, follow the same pattern.

For example, `ifAdminStatus` module is `IF-MIB`, so we want to query for
`IF-MIB::ifAdminStatus`::

    $ python -m canopsis.snmp.mibs --query
    {u'_id': u'IF-MIB::ifAdminStatus', u'oid': u'1.3.6.1.2.1.2.2.1.7'}

Then get the information of this OID::

    $ python -m canopsis.snmp.mibs --query 1.3.6.1.2.1.2.2.1.7
    {u'_id': u'1.3.6.1.2.1.2.2.1.7',
     u'access': u'readwrite',
     u'description': u'The desired state of the interface.  The testing(3) state\nindicates that no operational packets can be passed.  When a\nmanaged system initializes, all interfaces start with\nifAdminStatus in the down(2) state.  As a result of either\nexplicit management action or per configuration information\nretained by the managed system, ifAdminStatus is then\nchanged to either the up(1) or testing(3) states (or remains\nin the down(2) state).',
     u'moduleName': u'IF-MIB',
     u'nodetype': u'column',
     u'oid': u'1.3.6.1.2.1.2.2.1.7',
     u'status': u'current',
     u'syntax': {u'type': {u'basetype': u'Enumeration',
                           u'down': {u'nodetype': u'namednumber',
                                     u'number': u'2'},
                           u'testing': {u'nodetype': u'namednumber',
                                        u'number': u'3'},
                           u'up': {u'nodetype': u'namednumber', u'number': u'1'}}}}
