.. _admin-manage-services:

Services management
===================

.. toctree::
   :maxdepth: 2

   init
   admin-howto


supervisord behavior enhencement
================================

It is possible to run amqp2engines service in parallel mode setting the ``etc/hypcontrol.conf`` file like this:

   MIDDLEWARE=(collectd amqp2engines*)

This will increase the speed of start stop restart operation on amqp2engines service via `hypcontrol` command. These service operation are now available with

``service amqp2engines [mstart|mstop|mrestart]``

This made possible thanks to the sypervisord wildcard python lib bundled with canopsis.