.. _admin-practices-partition:

Partition table
===============

Using a separated partition for ``/opt`` is recommanded.
It must have at least **20 GB** free in order to let MongoDB start.

This is the commonly used partition table:

-  ``/boot`` 500MB Ext3
-  ``SWAP`` 4GB
-  LVM **vg0**

   -  **lv\_root**: ``/`` 20GB Ext3
   -  **lv\_log**: ``/var/log`` 5GB Ext3
   -  **lv\_app**: ``/opt/`` Max Ext3
