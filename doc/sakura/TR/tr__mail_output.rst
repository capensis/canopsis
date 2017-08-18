.. _TR__Mail_Output:

=================
Mail Output Media
=================

This document specifies the Mail output media.

References
==========

 - :ref:`FR::Task-Handling <FR__Task>`
 - :ref:`FR::Output-Media <FR__Output>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/06", "0.1", "Document creation", ""

Contents
========

Description
-----------

This output media will send a E-Mail, configured in its :ref:`task <FR__Task__Job>` configuration.

Task configuration
------------------

It **MUST** contains:

 - a sender, composed of:
    - a user
    - a group
    - a sender
 - one or more recipients
 - a subject (it **MAY** be a template)
 - a body (it **MAY** be a template)
 - zero or more attachments
 - a SMTP host/port
 - an indication telling if the mail to send is in HTML

Task handling
-------------

It **MUST** render the subject/body templates with the provided job context as data source.

Test Cases
----------

Case: Everything went good
~~~~~~~~~~~~~~~~~~~~~~~~~~

It **SHOULD** return the following informations:

 - an error code, equal to ``0``
 - a message explaining that there was no error


Case: Invalid configuration
~~~~~~~~~~~~~~~~~~~~~~~~~~~

If there is a missing field, or an invalid field in the task received, it **SHOULD**
return the following informations:

 - an error code, superior to ``0``
 - a message explaining the error

Case: Can't connect to SMTP host
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

It **SHOULD** return the following informations:

 - an error code, superior to ``0``
 - a message explaining the error

Case: Can't send E-Mail
~~~~~~~~~~~~~~~~~~~~~~~

It **SHOULD** return the following informations:

 - an error code, superior to ``0``
 - a message explaining the error
