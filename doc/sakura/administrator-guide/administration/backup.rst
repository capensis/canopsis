.. _Backup:

Backup
======

.. warning:: TODO: Update this page, Now backup is in ``gz`` format ... /!\\**

Cold Backup
-----------

::

    $ hypcontrol stop
    $ exit

(as *root* user)

::

    # cd /opt/
    # tar cfz canopsis.tgz canopsis  --exclude=canopsis/var/lib/mongodb/journal

Hot backup
----------

Canopsis has two default backup tasks :

-  Config
-  Mongo

Default backups directory is ``~/var/backups``.

There is no backup history : any new backup will erase the previous one.

Database: Mongodb
~~~~~~~~~~~~~~~~~

Make a zip archive of ``mongodump``.

To launch task:

::

    runtask task_backup mongo

Configuration files
~~~~~~~~~~~~~~~~~~~

Make a zip archive of ``~/etc/``.

To launch task:

::

    runtask task_backup config

Restore backup
--------------

Canopsis can be easily restored with two commands.

***You must restore configuration before the database.***

Configuration
~~~~~~~~~~~~~

In order to restore your configuration files, use the
``configrestore``\ tool.

Stop Canopsis first :

::

    hypcontrol stop

Restore configurations :

::

    cd ~/var/backups
    configrestore backup_config.zip

Database: Mongodb
~~~~~~~~~~~~~~~~~

Run Canopsis first.

::

    hypcontrol start

You have to unzip archive and give it to mongorestore with ``--drop``
option. Drop option will purge mongo before restore it. Without it, you
will have duplicate entries.

::

    cd ~/var/backups

    unzip backup_mongodb.zip 
    Archive:  backup_mongodb.zip
      inflating: backup_mongodb/canopsis/perfdata.fs.chunks.bson  
      inflating: backup_mongodb/canopsis/perfdata.bson  
      inflating: backup_mongodb/canopsis/system.indexes.bson  
      inflating: backup_mongodb/canopsis/events_log.bson  
      inflating: backup_mongodb/canopsis/events.bson  
      inflating: backup_mongodb/canopsis/cache.bson  
      inflating: backup_mongodb/canopsis/object.bson  

    mongorestore --drop backup_mongodb

``Mongorestore`` will display information about restored data.

And restart every services.

::

    hypcontrol restart
