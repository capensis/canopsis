<template lang="pug">
  v-container
    v-expansion-panel
      expansion-panel-content(v-for="(actions, groupKey) in groupedActions", :key="groupKey", lazy, ripple)
        div(slot="header") {{ groupKey }}
        v-card
          v-card-text
            .fixedTable
              table.table
                thead
                  tr
                    th
                    th(v-for="role in roles", :key="`role-header-${role._id}`") {{ role._id }}
                tbody
                  tr(v-for="action in actions", :key="`action-title-${action._id}`")
                    td.action-title {{ action._id }}
                    td.action-value(v-for="role in roles", :key="`role-action-${role._id}`")
                      input(
                      type="checkbox",
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

import VirtualList from 'vue-virtual-scroll-list';

import entitiesActionMixin from '@/mixins/entities/action';
import entitiesRoleMixin from '@/mixins/entities/role';

import ExpansionPanelContent from './expansion-panel.vue';

export default {
  components: { VirtualList, ExpansionPanelContent },
  mixins: [entitiesActionMixin, entitiesRoleMixin],
  data() {
    return {
      pending: true,
      groupedActions: { business: [], view: [], technical: [] },
      models: {},
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return Object.keys(this.changedRoles).length;
    },

    getCheckboxValue() {
      return (role, action, rightMask = 1) => {
        const checkSum = get(role, ['rights', action._id, 'checksum'], 0);
        const changedCheckSum = get(this.changedRoles, [role._id, action._id]);

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
                checked: this.getCheckboxValue(role, action, userRightMask),
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
              checked: this.getCheckboxValue(role, action),
            },
            on: {
              change: value => this.changeCheckboxValue(value, role, action),
            },
          },
        ];
      };
    },
  },
  async mounted() {
    const [{ data: actions }] = await Promise.all([
      this.fetchActionsListWithoutStore({ params: { limit: 10000 } }),
      this.fetchRolesList({ params: { limit: 10000 } }),
    ]);

    this.groupedActions = actions.reduce((acc, action) => {
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
  methods: {
    submit() {
      //       console.log(this.changedRoles);
    },

    changeCheckboxValue(event, role, action, rightType) {
      const currentCheckSum = get(role, ['rights', action._id, 'checksum'], 0);
      const { checked: value } = event.target;
      const factor = value ? 1 : -1;

      /**
       * If we don't have changes for role
       */
      if (!this.changedRoles[role._id]) {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles, role._id, { [action._id]: nextCheckSum });

        /**
         * If we have changes for role but we don't have changes for action
         */
      } else if (isUndefined(this.changedRoles[role._id][action._id])) {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles[role._id], action._id, nextCheckSum);

        /**
         * If we have changes for role and for action
         */
      } else {
        const nextCheckSum = !rightType ?
          Number(value) : this.changedRoles[role._id][action._id] + (factor * rightType);

        if (currentCheckSum === nextCheckSum) {
          if (Object.keys(this.changedRoles[role._id]).length === 1) {
            this.$delete(this.changedRoles, role._id);
          } else {
            this.$delete(this.changedRoles[role._id], action._id);
          }
        } else {
          this.$set(this.changedRoles[role._id], action._id, nextCheckSum);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  $cellWidth: 100px;

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
