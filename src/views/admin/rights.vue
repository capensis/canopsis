<template lang="pug">
  v-container.fixedTable
    table.table.table-bordered
      thead
        tr
          th
          th
          th(v-for="role in roles", :key="`role-header-${role._id}`") {{ role._id }}
      tbody
        template(v-for="(actions, groupKey) in groupedActions")
          tr(:key="`actions-group-title-${groupKey}`")
            td.group-title(:rowspan="actions.length + 1") {{ groupKey }}
          tr(v-for="action in actions", :key="`action-title-${action._id}`")
            td.action-title {{ action._id }}
            td.action-value(v-for="role in roles", :key="`role-action-${role._id}`")
              v-checkbox(
              v-for="(checkbox, index) in getCheckboxes(role, action)",
              :key="`role-${role._id}-action-${action._id}-checkbox-${index}`",
              v-bind="checkbox.bind",
              v-on="checkbox.on"
              )
    v-btn(v-show="hasChanges", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import get from 'lodash/get';
import isEmpty from 'lodash/isEmpty';
import isUndefined from 'lodash/isUndefined';

import entitiesActionMixin from '@/mixins/entities/action';
import entitiesRoleMixin from '@/mixins/entities/role';

export default {
  mixins: [entitiesActionMixin, entitiesRoleMixin],
  data() {
    return {
      pending: true,
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return !isEmpty(this.changedRoles);
    },

    groupedActions() {
      return this.actions.reduce((acc, action) => {
        if (action.id.startsWith('view') || action.id.startsWith('userview')) {
          acc.view.push(action);
        } else if (action.id.startsWith('models')) {
          acc.technical.push(action);
        } else {
          acc.business.push(action);
        }

        return acc;
      }, { business: [], view: [], technical: [] });
    },

    getCheckboxValue() {
      return (role, action, rightMask = 1) => {
        const checkSum = get(role, `rights.${action._id}.checksum`, 0);
        const changedCheckSum = get(this.changedRoles, `${role._id}.${action._id}`);

        const currentCheckSum = isUndefined(changedCheckSum) ? checkSum : changedCheckSum;
        const actionRightType = currentCheckSum & rightMask;

        return actionRightType === rightMask;
      };
    },

    getCheckboxes() {
      const { USERS_RIGHTS_MASKS, USERS_ACTIONS_TYPES } = this.$constants;

      return (role, action) => {
        if (action.type) {
          let masks = [];

          if (action.type === USERS_ACTIONS_TYPES.crud) {
            masks = ['create', 'read', 'update', 'delete'];
          }

          if (action.type === USERS_ACTIONS_TYPES.rw) {
            masks = ['read', 'update', 'delete'];
          }

          return masks.map((userRightMaskKey) => {
            const userRightMask = USERS_RIGHTS_MASKS[userRightMaskKey];

            return {
              bind: {
                inputValue: this.getCheckboxValue(role, action, userRightMask),
              },
              on: {
                change: value => this.changeCheckboxValue(value, role, action, userRightMask),
              },
            };
          });
        }

        return [
          {
            bind: {
              inputValue: this.getCheckboxValue(role, action),
            },
            on: {
              change: value => this.changeCheckboxValue(value, role, action),
            },
          },
        ];
      };
    },
  },
  mounted() {
    this.fetchActionsList({
      params: { limit: 10000 },
    });

    this.fetchRolesList({
      params: { limit: 10000 },
    });
  },
  methods: {
    submit() {
      // console.log(this.changedRoles);
    },

    changeCheckboxValue(value, role, action, rightType) {
      const currentCheckSum = get(role.rights, `${action._id}.checksum`, 0);
      const factor = value ? 1 : -1;

      if (!this.changedRoles[role._id]) {
        this.$set(this.changedRoles, role._id, {});
      }

      if (!isUndefined(this.changedRoles[role._id][action._id])) {
        const nextCheckSum = !rightType ?
          Number(value) : this.changedRoles[role._id][action._id] + (factor * rightType);

        if (currentCheckSum === nextCheckSum) {
          this.$delete(this.changedRoles[role._id], action._id);
        } else {
          this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      } else {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  $cellHeight: 20px;
  $cellWidth: 100px;
  $cellPadding: 5px;

  $cellsWide: 5;
  $cellsHigh: 15;

  .fixedTable {
    .table {
      background-color: white;
      width: 100%;

      tr {
        td, th {
          vertical-align: top;
          min-width: $cellWidth;
        }

        .group-title {
          width: 100px;
        }

        .action-title {
          width: 200px;
        }

        .action-value {
          text-align: center;
        }
      }

      & /deep/ {
        .v-input {
          margin: 0;

          .v-input--selection-controls__input {
            display: block;
            margin: auto;
          }
        }
      }
    }
  }
</style>
