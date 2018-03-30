Canopsis on boot
================

For start Canopsis on boot you can use
`/opt/canopsis/etc/init.d/canopsis` as init script.

Create symlink

```bash
ln -s /opt/canopsis/etc/init.d/canopsis /etc/init.d/
```

RHEL and CentOS
---------------

Add script to startup process

```bash
chkconfig --add canopsis
```

Check configuration

```bash
chkconfig | grep canopsis
```

Debian and Ubuntu
-----------------

Add script to startup process

```bash
update-rc.d canopsis defaults
```
