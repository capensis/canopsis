deploy-ansible
=========

Deploy Canopsis by ansible or mount it in a vagrantbox

Requirements
------------

- Ansible 2.4.3 - Use a virtualenv and `pip install "ansible<2.4.4"`
- Host using ansible as to get right on git.canopsis.com to fetch canopsis sources (ex: ssk_key)

### Optional

If you want to quickly test canopsis install those two component 

- Vagrant
- VirtualBox

Or run Ansible manually by customising the inventory.

Deploy Canopsis
---------------

```bash
cd deploy-ansible/<ansible roles directory>/deploy-ansible

cp <package> playbook/files/
```

Then edit `playbook/group_vars/all.yml` to change the package version if needed.

Vagrant:

```bash
vagrant up

# to run only ansible after initial deploy
vagrant up --provision
```

Manual:

```bash
# THIS RESETS ALL DATABASES
./install-standalone.sh -i inventory_custom

# this updates databases
./update-standalone.sh -i inventory_custom
```

* WebUI: http://localhost:8082
* RabbitMQ admin interface: http://localhost:15672
* RabbitMQ AMQP bus: amqp://localhost:5672
* MongoDB: mongo localhost -p 27017

Role Variables
--------------

See `playbook/roles/canopsis/defaults/main.yml`, then `playbook/group_vars/all.yml`.

Also, use variables from vendored roles before doing the same task yourself in the `canopsis` role.

Development
-----------

Checks dependencies and upgrades the virtualenv:

```
./install-dev.sh /mnt/canopsis/sources/canopsis
```

Do not checks for dependencies, only upgrade the virtualenv:

```
./install-dev-fast.sh /mnt/canopsis/sources/canopsis
```

Tags
----

Some tags you can use with `--tags tag,list`:

 * `cps_services`: run only tasks related to the state of canopsis services (engines, webserver...)
 * `cps_install_package`: this tag is **REQUIRED** when running ansible **AFTER** the package has been installed, aka when you want to deploy Canopsis from the packages's playbook.

Groups
------

 * `all`: the default group.
 * `go`: this group overrides some variables so you will not start engines that have been replaces by Go ones.

Dependencies
------------

- Ansible

Example Playbook
----------------

See the playbook/canopsis.yml file.

License
-------

BSD

Author Information
------------------

Capensis/Canopsis 

Vendoring
---------

See roles imported in `playbook/vendored_roles`.

To update those roles:

```bash
ansible-galaxy install -r requirements.yml -p playbook/vendored_roles/
```

**NEVER CHANGE ANYTHING MANUALLY IN THESE ROLES.** If you do so, any fix or change **will** be wiped out in an upgrade.

Why vendoring
-------------

Simplicity :

 * Builds are simpler
 * `ansible-galaxy` takes care of doing the job of upgrading if required, and is not specific to one SCM (unlike git submodules or subrepos)
 * Keep track of modifications made to external roles so we can check regressions
