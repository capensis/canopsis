Authentication with ldap
========================

Allow plain authentication
--------------------------

**Warning:** For use external authentication password must be
transmitted in plain format to the API. You must configure HTTPS for
more security.

Edit ``var/www/global_options.js`` and add this in ``global_options``:

::

    auth_plain: true

Configure LDAP binding
----------------------

For configure LDAP binding, connect with ``root`` access on webUI and
goto ``Build > Accounts > Ldap``.

Ex: Standard LDAP (OpenLDAP ...)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

::

    URI: ldap://<YOUR_SERVER>
    Base DN: ou=People,dc=canopsis,dc=org
    User DN: uid=%s,ou=People,dc=canopsis,dc=org
    Domain:
    User filter: (&(uid=%s)(objectclass=inetOrgPerson))
    Lastname field: sn
    Firstname field: givenName
    Mail field: mail

Ex: Active Directory
~~~~~~~~~~~~~~~~~~~~

::

    URI: ldap://<YOUR_SERVER>
    Base DN: cn=Users,dc=canopsis,dc=org
    User DN: 
    Domain: canopsis.org
    User filter: (&(objectClass=person)(cn=%s))
    Lastname field: sn
    Firstname field: givenName
    Mail field: mail
