Mail Output Media
=================

This document specifies the Mail output media.

References
----------

> -   FR::Task-Handling &lt;FR\_\_Task&gt;
> -   FR::Output-Media &lt;FR\_\_Output&gt;

Updates
-------

Contents
--------

### Description

This output media will send a E-Mail, configured in its
task &lt;FR\_\_Task\_\_Job&gt; configuration.

### Task configuration

It **MUST** contains:

> -   
>
>     a sender, composed of:
>
>     :   -   a user
>         -   a group
>         -   a sender
>
> -   one or more recipients
> -   a subject (it **MAY** be a template)
> -   a body (it **MAY** be a template)
> -   zero or more attachments
> -   a SMTP host/port
> -   an indication telling if the mail to send is in HTML

### Task handling

It **MUST** render the subject/body templates with the provided job
context as data source.

### Test Cases

#### Case: Everything went good

It **SHOULD** return the following informations:

> -   an error code, equal to `0`
> -   a message explaining that there was no error

#### Case: Invalid configuration

If there is a missing field, or an invalid field in the task received,
it **SHOULD** return the following informations:

> -   an error code, superior to `0`
> -   a message explaining the error

#### Case: Can't connect to SMTP host

It **SHOULD** return the following informations:

> -   an error code, superior to `0`
> -   a message explaining the error

#### Case: Can't send E-Mail

It **SHOULD** return the following informations:

> -   an error code, superior to `0`
> -   a message explaining the error

