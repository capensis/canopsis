<template lang="pug">
  v-container.admin-rights
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.rights') }}
    div.progress-wrapper
      v-fade-transition
        v-layout.white.progress(v-show="pending", column)
          v-progress-circular(indeterminate, color="primary")
      v-tabs(fixed-tabs)
        template(v-for="(rights, groupKey) in groupedRights")
          v-tab(:key="`tab-${groupKey}`") {{ groupKey }}
          v-tab-item(:key="`tab-item-${groupKey}`")
            v-card(v-if="hasReadAnyRoleAccess")
              v-card-text
                table.table
                  thead
                    tr
                      th
                      th(v-for="role in roles", :key="`role-header-${role._id}`") {{ role._id }}
                  tbody
                    tr(v-for="right in rights", :key="`right-title-${right._id}`")
                      td {{ right.desc }}
                      td(v-for="role in roles", :key="`role-right-${role._id}`")
                        v-checkbox-functional(
                        v-for="(checkbox, index) in getCheckboxes(role, right)",
                        :key="`role-${role._id}-right-${right._id}-checkbox-${index}`",
                        v-bind="checkbox.bind",
                        v-on="checkbox.on",
                        :disabled="!hasUpdateAnyActionAccess"
                        )
    v-layout(v-show="hasUpdateAnyActionAccess && hasChanges")
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
      v-btn(@click="cancel") {{ $t('common.cancel') }}
    .fab(v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess || hasCreateAnyActionAccess")
      v-layout(column)
        refresh-btn(@click="fetchRightsList")
        v-speed-dial(
        v-model="fab",
        direction="left",
        transition="slide-y-reverse-transition"
        )
          v-btn(slot="activator", color="primary", fab, v-model="fab")
            v-icon add
            v-icon close
          v-tooltip(v-if="hasCreateAnyUserAccess", top)
            v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateUserModal")
              v-icon people
            span {{ $t('modals.createUser.title') }}
          v-tooltip(v-if="hasCreateAnyRoleAccess", top)
            v-btn(slot="activator", fab, dark, small, color="deep-purple ", @click.stop="showCreateRoleModal")
              v-icon supervised_user_circle
            span {{ $t('modals.createRole.title') }}
          v-tooltip(v-if="hasCreateAnyActionAccess", top)
            v-btn(slot="activator", fab, dark, small, color="teal", @click.stop="showCreateRightModal")
              v-icon verified_user
            span {{ $t('modals.createRight.title') }}
</template>

<script>
import { get, isEmpty, isUndefined, transform } from 'lodash';
import flatten from 'flat';

import { MODALS, USERS_RIGHTS, USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES, NOT_COMPLETED_USER_RIGHTS_KEYS } from '@/constants';
import { generateRoleRightByChecksum } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import entitiesRightMixin from '@/mixins/entities/right';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';
import rightsTechnicalActionMixin from '@/mixins/rights/technical/action';

import RefreshBtn from '@/components/other/view/refresh-btn.vue';

export default {
  components: {
    RefreshBtn,
  },
  mixins: [
    authMixin,
    popupMixin,
    modalMixin,
    entitiesRightMixin,
    entitiesRoleMixin,
    entitiesViewGroupMixin,
    rightsTechnicalUserMixin,
    rightsTechnicalRoleMixin,
    rightsTechnicalActionMixin,
  ],
  data() {
    return {
      fab: false,
      pending: false,
      groupedRights: { business: [], view: [], technical: [] },
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return !isEmpty(this.changedRoles);
    },

    getCheckboxValue() {
      return (role, right, rightMask = 1) => {
        const checkSum = get(role, ['rights', right._id, 'checksum'], 0);
        const changedCheckSum = get(this.changedRoles, [role._id, right._id]);

        const currentCheckSum = isUndefined(changedCheckSum) ? checkSum : changedCheckSum;
        const rightType = currentCheckSum & rightMask;

        return rightType === rightMask;
      };
    },

    getCheckboxes() {
      return (role, right) => {
        if (right.type) {
          let masks = [];

          if (right.type === USERS_RIGHTS_TYPES.crud) {
            masks = ['create', 'read', 'update', 'delete'];
          }

          if (right.type === USERS_RIGHTS_TYPES.rw) {
            masks = ['read', 'update', 'delete'];
          }

          return masks.map((userRightMaskKey) => {
            const userRightMask = USERS_RIGHTS_MASKS[userRightMaskKey];

            return {
              bind: {
                inputValue: this.getCheckboxValue(role, right, userRightMask),
                label: userRightMaskKey,
              },
              on: {
                change: value => this.changeCheckboxValue(value, role, right, userRightMask),
              },
            };
          });
        }

        return [
          {
            bind: {
              inputValue: this.getCheckboxValue(role, right),
            },
            on: {
              change: value => this.changeCheckboxValue(value, role, right),
            },
          },
        ];
      };
    },
  },
  mounted() {
    this.fetchRightsList();
  },
  methods: {
    async fetchRightsList() {
      this.pending = true;

      await this.fetchList();

      this.pending = false;
    },

    showCreateUserModal() {
      this.showModal({
        name: MODALS.createUser,
      });
    },

    showCreateRoleModal() {
      this.showModal({
        name: MODALS.createRole,
      });
    },

    showCreateRightModal() {
      this.showModal({
        name: MODALS.createRight,
      });
    },

    clearChangedRoles() {
      this.changedRoles = {};
    },

    cancel() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: this.clearChangedRoles,
        },
      });
    },

    submit() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: this.updateRoles,
        },
      });
    },

    async updateRoles() {
      try {
        this.pending = true;

        await Promise.all(Object.keys(this.changedRoles).map((roleId) => {
          const changedRoleRights = this.changedRoles[roleId];
          const role = this.getRoleById(roleId);

          const newRights = transform(
            changedRoleRights,
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

        this.addSuccessPopup({ text: this.$t('success.default') });
        this.clearChangedRoles();

        this.pending = false;
      } catch (err) {
        this.addErrorPopup({ text: this.$t('errors.default') });
      }
    },

    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {Object} role
     * @param {Object} right
     * @param {number} rightType
     */
    changeCheckboxValue(value, role, right, rightType) {
      const currentCheckSum = get(role, ['rights', right._id, 'checksum'], 0);
      const factor = value ? 1 : -1;

      /**
       * If we don't have changes for role
       */
      if (!this.changedRoles[role._id]) {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles, role._id, { [right._id]: nextCheckSum });

        /**
         * If we have changes for role but we don't have changes for right
         */
      } else if (isUndefined(this.changedRoles[role._id][right._id])) {
        const nextCheckSum = !rightType ?
          Number(value) : currentCheckSum + (factor * rightType);

        this.$set(this.changedRoles[role._id], right._id, nextCheckSum);

        /**
         * If we have changes for role and for right
         */
      } else {
        const nextCheckSum = !rightType ?
          Number(value) : this.changedRoles[role._id][right._id] + (factor * rightType);

        if (currentCheckSum === nextCheckSum) {
          if (Object.keys(this.changedRoles[role._id]).length === 1) {
            this.$delete(this.changedRoles, role._id);
          } else {
            this.$delete(this.changedRoles[role._id], right._id);
          }
        } else {
          this.$set(this.changedRoles[role._id], right._id, nextCheckSum);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }
      }
    },

    /**
     * Fetch rights and roles lists
     *
     * @returns void
     */
    async fetchList() {
      const [{ data: rights }] = await Promise.all([
        this.fetchRightsListWithoutStore({ params: { limit: 10000 } }),
        this.fetchRolesList({ params: { limit: 10000 } }),
      ]);

      const allViews = this.groups.reduce((acc, { views }) => acc.concat(views), []);
      const allRightsIds = flatten(USERS_RIGHTS);
      const allTechnicalRightsIds = flatten(USERS_RIGHTS.technical);
      const allBusinessRightsIds = flatten(USERS_RIGHTS.business);

      this.groupedRights = rights.reduce((acc, right) => {
        const rightId = String(right._id);
        const view = allViews.find(({ _id }) => _id === rightId);

        if (view) {
          acc.view.push({
            ...right,

            desc: right.desc.replace(view._id, view.name),
          });
        } else if (Object.values(allTechnicalRightsIds).indexOf(rightId) !== -1) {
          acc.technical.push(right);
        } else if (
          Object.values(allBusinessRightsIds).indexOf(rightId) !== -1 ||
          NOT_COMPLETED_USER_RIGHTS_KEYS.some(userRightKey => rightId.startsWith(allRightsIds[userRightKey]))
        ) {
          acc.business.push(right);
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
