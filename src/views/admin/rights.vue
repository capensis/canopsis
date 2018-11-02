<template lang="pug">
  v-container.admin-rights
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.rights') }}
    div.progress-wrapper
      v-fade-transition
        v-layout.white.progress(v-show="pending", column)
          v-progress-circular(indeterminate, color="primary")
      v-expansion-panel
        v-expansion-panel-content(
        v-for="(actions, groupKey) in groupedActions",
        :key="groupKey",
        lazy,
        lazyWithUnmount,
        ripple
        )
          div(slot="header") {{ groupKey }}
          v-card
            v-card-text
              table.table
                thead
                  tr
                    th
                    th(v-for="role in roles", :key="`role-header-${role._id}`") {{ role._id }}
                tbody
                  tr(v-for="action in actions", :key="`action-title-${action._id}`")
                    td.action-title {{ action._id }}
                    td.action-value(v-for="role in roles", :key="`role-action-${role._id}`")
                      v-checkbox-functional(
                      v-for="(checkbox, index) in getCheckboxes(role, action)",
                      :key="`role-${role._id}-action-${action._id}-checkbox-${index}`",
                      v-bind="checkbox.bind",
                      v-on="checkbox.on"
                      )
    v-layout(v-show="hasChanges")
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
      v-btn(@click="cancel") {{ $t('common.cancel') }}
</template>

<script>
import get from 'lodash/get';
import isEmpty from 'lodash/isEmpty';
import isUndefined from 'lodash/isUndefined';
import transform from 'lodash/transform';

import { generateRoleRightByChecksum } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal/modal';
import entitiesActionMixin from '@/mixins/entities/action';
import entitiesRoleMixin from '@/mixins/entities/role';

export default {
  mixins: [authMixin, modalMixin, entitiesActionMixin, entitiesRoleMixin],
  data() {
    return {
      pending: false,
      groupedActions: { business: [], view: [], technical: [] },
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return !isEmpty(this.changedRoles);
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
  async mounted() {
    this.pending = true;

    await this.fetchList();

    this.pending = false;
  },
  methods: {
    clearChangedRoles() {
      this.changedRoles = {};
    },

    cancel() {
      this.showModal({
        name: this.$constants.MODALS.confirmation,
        config: {
          action: this.clearChangedRoles,
        },
      });
    },

    submit() {
      this.showModal({
        name: this.$constants.MODALS.confirmation,
        config: {
          action: this.updateRoles,
        },
      });
    },

    async updateRoles() {
      this.pending = true;

      await Promise.all(Object.keys(this.changedRoles).map((roleId) => {
        const changedRoleActions = this.changedRoles[roleId];
        const role = this.getRoleById(roleId);

        const newRights = transform(
          changedRoleActions,
          (acc, value, key) => acc[key] = generateRoleRightByChecksum(value),
        );

        return this.createRole({ data: { ...role, rights: { ...role.rights, ...newRights } } });
      }));

      /**
       * If current user role changed
       */
      if (this.changedRoles[this.currentUser.role]) {
        await this.fetchCurrentUser();
      }

      this.clearChangedRoles();

      this.pending = false;
    },

    changeCheckboxValue(value, role, action, rightType) {
      const currentCheckSum = get(role, ['rights', action._id, 'checksum'], 0);
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

    async fetchList() {
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
  },
};
</script>

<style lang="scss" scoped>
  .admin-rights {
    & /deep/ .v-expansion-panel__body {
      overflow: auto;
    }
  }

  .table {
    background-color: white;
    width: 100%;

    tr {
      td, th {
        vertical-align: top;
        padding: 5px;
      }
    }

    & /deep/ {
      .v-input {
        margin: 0;
      }
    }
  }

  .progress {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    opacity: .4;
    z-index: 1;

    & /deep/ .v-progress-circular {
      top: 50%;
      left: 50%;
      margin-top: -16px;
      margin-left: -16px;
    }
  }

  .progress-wrapper {
    position: relative;
  }
</style>
