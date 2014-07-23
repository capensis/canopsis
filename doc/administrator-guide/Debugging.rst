Debugging
=========

Exporting report
----------------

Purpose
~~~~~~~

In order to simplify the debugging of **Canopsis**, a reporting tool has
been developed to collect useful informations about the running
instance, allowing us to diagnose more easily the problem.

Howto
~~~~~

The ``export_dbginfo`` tool, included inside the Canopsis environment
and provided by package *canotools*, generate a tarball containing the
following data :

-  services state (via supervisord) ;
-  RabbitMQ informations about :

   -  connections ;
   -  channels ;
   -  nodes ;
   -  users ;
   -  exchanges ;
   -  virtual hosts ;
   -  queues ;
   -  bindings ;
   -  permissions.

-  RabbitMQ overview and definitions ;
-  all logs located in ``~/var/log`` ;
-  processes (CPU/RAM used per process owned by the user ``canopsis``).

The generated tarball has the following name :
``~/tmp/canodebug.`date +"%F.%H-%M-%S"`.tar.xz``