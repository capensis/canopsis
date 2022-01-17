# deploy-ansible

Deploy Canopsis RPM packages with Ansible (complete the installation, add some external components).

As of early 2022, this is still required, but most of its content should be moved to RPM dependencies and RPM post-tasks (WIP).

## Requirements

- Ansible 2.8.7 - Use a virtualenv and `pip install "ansible==2.8.7"`
- Always keep it as simple as possible; don't hesitate breaking any Ansible "best practice" when it makes things simpler and much easier to maintain.

⚠️ **DON'T, DON'T, DON'T** ⚠️ UPGRADE ANSIBLE, NOT EVEN TO A MINOR VERSION. It's just a perpetual breaking machine.

## Handling role dependencies

If your role *really* needs some Python dependencies, add them to `/community/docker/build/pip-ansible.sh` inside this repository.

But see above: keep it as simple as possible.

## Role Variables

See `playbook/roles/canopsis/defaults/main.yml`, then `playbook/group_vars/all.yml`.

Also, use variables from vendored roles before doing the same task yourself in the `canopsis` role.

## License

See the license of each individual role within this playbook.

## Author Information

Capensis/Canopsis and the author of each individual role.
