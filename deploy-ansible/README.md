# deploy-ansible

Deploy Canopsis with Ansible

## Requirements

- Ansible 2.8.7 - Use a virtualenv and `pip install "ansible==2.8.7"`

## Role Variables

See `playbook/roles/canopsis/defaults/main.yml`, then `playbook/group_vars/all.yml`.

Also, use variables from vendored roles before doing the same task yourself in the `canopsis` role.

## License

BSD

## Author Information

Capensis/Canopsis 

## Vendoring

See roles imported in `playbook/vendored_roles`.

To update those roles:

```sh
ansible-galaxy install -r requirements.yml -p playbook/vendored_roles/
```

**NEVER CHANGE ANYTHING MANUALLY IN THESE ROLES.** If you do so, any fix or change **will** be wiped out in an upgrade.

### Why vendoring

Simplicity :

 * Builds are simpler
 * `ansible-galaxy` takes care of doing the job of upgrading if required, and is not specific to one SCM (unlike git submodules or subrepos)
 * Keep track of modifications made to external roles so we can check regressions
