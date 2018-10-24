Canopsis Statistics Engine
==========================

This document describes the statistics generation process.

References
----------

List of referenced functional requirements:

> -   FR::Alarm &lt;FR\_\_Alarm&gt;
> -   FR::Engine &lt;FR\_\_Engine&gt;
> -   FR::Webservice &lt;FR\_\_Webservice&gt;

Updates
-------

Contents
--------

### Description

Stats are a feature of canopsis producing metrics about alarms
&lt;FR\_\_Alarm&gt; activity and user sessions durations.

Stats **CAN** displayed on frontend thanks to a proprietary brick :
brick-statstable.

Information details about stats metrics is available in this brick.

### Stats engine

The stats engine has currently only a beat processing. On this beat, it
computes stats from `*_alarm` and `session` mongo collections, storing
datas in influxdb.

Those condensed data can be retrieved much faster, and quickly filtered
via authors or custom tags.

#### Custom tags

Alarms &lt;FR\_\_Alarm&gt; can be tagged via `extra_fields`. Those extra
fields are exported as tags in influxdb.

### Webservice

There is a stats webservice serving data to UI via the brick statstable.

In order to enable it (not enabled by default), add the line `stats=1`
in the `webservices` section of `PREFIX/etc/webserver.conf` and restart
webserver.

Functional Tests
----------------

Functionnal tests have been implemented to check most of produced
metrics. But since there is currently no place in canopsis deployment
workflow to insert it, this work has not been pushed.

<div class="admonition warning">

**TODO:** As it is difficult to simulate user authentication and
keepalive workflow, there are currently no functionnal tests for session
duration metrics.

</div>
