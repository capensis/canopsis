<a href="http://www.canopsis.org" >
    <img src="https://github.com/capensis/canopsis/wiki/images/logo_canopsis.png"/>
</a>

[![pipeline status](https://git.canopsis.net/canopsis/canopsis/badges/develop/pipeline.svg)](https://git.canopsis.net/canopsis/canopsis/commits/develop)

## What is Canopsis ?

[Canopsis](http://canopsis.org) is an open-source [hypervisor](http://www.capensis.fr/solutions/hypervision/) whose goal is to <a href="https://github.com/capensis/canopsis/wiki/consolidation" target="_blank">aggregate/consolidate</a> information and events (containing <a href="https://github.com/capensis/canopsis/wiki/metrics">metrics</a> of different types such as performance, availability, etc.) coming from multiple sources to create a global solution for <a href="https://github.com/capensis/canopsis/wiki/Dashboard" target="_blank">monitoring</a> and <a href="https://github.com/capensis/canopsis/wiki/engines" target="_blank">administrating</a> resources.

Built to last on top of [proven Open Source technologies by and for all IT professionals](http://www.capensis.fr/solutions/supervision/). It is an event based architecture and it is modular by design. Plug your infrastructure tools like `Syslog`, `Nagios`, [`Shinken`](https://github.com/naparuba/shinken), `...` to [Canopsis](http://canopsis.org) and you're ready to go.

A <a href="https://github.com/capensis/canopsis/wiki/Glossary" target="_blank">Glossary</a> page is also given for better descriptions about canopsis concepts.

## How to try ?

You can try Canopsis on demo platform:
* Master branch (stable): http://sakura.canopsis.net
* Devel branch (unstable): http://sakura-devel.canopsis.net

## How to install ?

For a standalone version of Canopsis on CentOS 7 and Debian 8 (Jessie), you can use the [auto deploy](/STANDALONE-DEPLOY.md).

For a custom install, you can still use packages and follow [this documentation](/doc/docs/fr/guide_administrateur/package_install.md).

## How to use ?

To know more about Canopsis, have a look at <a href="https://canopsis.readthedocs.io" target="_blank">this documentation</a>

## Other links

* <a href="http://www.canopsis.org" target="_blank">Community</a>
* <a href="http://forums.monitoring-fr.org/index.php?board=127.0" target="_blank">Forum (french)</a>

## Tested dependencies

The following software versions have been tested to play nice with Canopsis: 


|Software  | Canopsis 2.3 | Canopsis  >=2.5.7 |
|----------|--------------|-------------------|
|Rabbit    | 3.6.9        | 3.6.9             |
|Erlang    | 19.3         | 19.3              |
|MongoDB   | <= 3.4.10    | 3.4.10            |
|InfluxDB  | 1.2.2        | 1.2.2             |
|Python    | 2.7.13       | 2.7.13            |
|Gunicorn  | 19.7.1       | 19.7.1            |
|HAProxy   | 1.5.18       | 1.5.18            |

