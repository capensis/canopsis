Partition table
===============

-  ``/boot`` 500MB Ext3
-  ``SWAP`` 4GB
-  LVM **vg0**

   -  **lv\_root**: ``/`` 20GB Ext3
   -  **lv\_log**: ``/var/log`` 5GB Ext3
   -  **lv\_app**: ``/opt/`` Max Ext3