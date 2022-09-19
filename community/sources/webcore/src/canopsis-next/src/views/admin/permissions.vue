<template lang="pug">
  v-container.admin-rights
    c-page-header
    div.position-relative
      v-fade-transition
        v-layout.white.progress(v-show="pending", column)
          v-progress-circular(indeterminate, color="primary")
      v-tabs(v-if="hasReadAnyRoleAccess", fixed-tabs, slider-color="primary")
        template(v-for="(permissions, groupKey) in groupedPermissions")
          v-tab(:key="`tab-${groupKey}`") {{ groupKey }}
          v-tab-item.white(:key="`tab-item-${groupKey}`")
            permissions-table-wrapper(
              :permissions="permissions",
              :roles="preparedRoles",
              :changed-roles="changedRoles",
              :disabled="!hasUpdateAnyPermissionAccess",
              :sort-by="getSortBy(groupKey)",
              @change="changeCheckboxValue"
            )
    v-layout.submit-button.mt-3(v-show="hasChanges")
      v-btn.primary.ml-3(@click="submit") {{ $t('common.submit') }}
      v-btn(@click="cancel") {{ $t('common.cancel') }}
    permissions-fab-buttons(@refresh="fetchList")
</template>

<script>
import { get, isEmpty, sortBy } from 'lodash';

import { MAX_LIMIT, MODALS } from '@/constants';

import { getGroupedPermissions } from '@/helpers/permission';
import { roleToForm, formToRole } from '@/helpers/forms/role';

import { authMixin } from '@/mixins/auth';
import { entitiesPermissionsMixin } from '@/mixins/entities/permission';
import { entitiesRoleMixin } from '@/mixins/entities/role';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { entitiesPlaylistMixin } from '@/mixins/entities/playlist';
import { permissionsTechnicalRoleMixin } from '@/mixins/permissions/technical/role';
import { permissionsTechnicalPermissionMixin } from '@/mixins/permissions/technical/permission';

import PermissionsTableWrapper from '@/components/other/permission/admin/permissions-table-wrapper.vue';
import PermissionsFabButtons from '@/components/other/permission/admin/permissions-fab-buttons.vue';

export default {
  components: {
    PermissionsTableWrapper,
    PermissionsFabButtons,
  },
  mixins: [
    authMixin,
    entitiesPermissionsMixin,
    entitiesRoleMixin,
    entitiesViewGroupMixin,
    entitiesPlaylistMixin,
    permissionsTechnicalRoleMixin,
    permissionsTechnicalPermissionMixin,
  ],
  data() {
    return {
      pending: false,
      groupedPermissions: {
        business: [],
        view: [],
        technical: [],
        api: [],
      },
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return !isEmpty(this.changedRoles);
    },

    preparedRoles() {
      return sortBy(this.roles, [({ _id: name }) => name.toLowerCase()]).map(roleToForm);
    },

    allViews() {
      return this.groups.reduce((acc, { views }) => acc.concat(views), []);
    },

    viewsPriorityById() {
      return this.allViews.reduce((acc, view, index) => {
        acc[view._id] = index;

        return acc;
      }, {});
    },

    playlistPriorityById() {
      return this.playlists.reduce((acc, playlist, index) => {
        acc[playlist._id] = index;

        return acc;
      }, {});
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    getPlaylistSort(permission) {
      return this.playlistPriorityById[permission._id] + this.allViews.length;
    },

    getViewSort(permission) {
      return this.viewsPriorityById[permission._id] ?? this.getPlaylistSort(permission);
    },

    // eslint-disable-next-line consistent-return
    getSortBy(key) {
      if (key === 'view') {
        return this.getViewSort;
      }
    },

    /**
     * Clear changed roles
     *
     * @returns void
     */
    clearChangedRoles() {
      this.changedRoles = {};
    },

    /**
     * Show the confirmation modal window for clearing a changedRoles
     *
     * @returns void
     */
    cancel() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.clearChangedRoles,
        },
      });
    },

    /**
     * Show the confirmation modal window for submitting a changedRoles
     *
     * @returns void
     */
    submit() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.updateRoles,
        },
      });
    },

    /**
     * Send request for update changedRoles and fetchCurrentUser if it needed
     *
     * @returns void
     */
    async updateRoles() {
      try {
        this.pending = true;

        await Promise.all(Object.entries(this.changedRoles).map(([roleId, permissions]) => {
          const role = this.getRoleById(roleId);
          const roleForm = roleToForm(role);

          roleForm.permissions = {
            ...roleForm.permissions,
            ...permissions,
          };

          return this.updateRole({
            id: roleId,
            data: formToRole(roleForm),
          });
        }));

        /**
         * If current user role changed
         */
        if (this.changedRoles[this.currentUser.role._id]) {
          await this.fetchCurrentUser();
        }

        await this.fetchRolesList({ params: { limit: MAX_LIMIT } });

        this.$popups.success({ text: this.$t('success.default') });

        this.clearChangedRoles();
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      } finally {
        this.pending = false;
      }
    },

    getNextActions(actions, newAction, value) {
      return value
        ? [...actions, newAction]
        : actions.filter(roleAction => roleAction !== newAction);
    },

    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {Object} role
     * @param {Object} permission
     * @param {string} action
     */
    changeCheckboxValue(value, role, permission, action) {
      const changedRole = this.changedRoles[role._id];
      const changedPermissionActions = get(this.changedRoles, [role._id, permission._id]);

      const currentActions = get(role, ['permissions', permission._id], []);

      /**
       * If we have changes for role and permission
       */
      if (changedPermissionActions) {
        const nextActions = this.getNextActions(changedPermissionActions, action, value);

        const isEqualOriginPermission = currentActions.length === nextActions.length
          && currentActions.every(nextAction => nextActions.includes(nextAction));

        if (isEqualOriginPermission) {
          if (Object.keys(changedRole).length === 1) {
            this.$delete(this.changedRoles, role._id);
          } else {
            this.$delete(changedRole, permission._id);
          }
        } else {
          this.$set(changedRole, permission._id, nextActions);
        }

        if (isEmpty(this.changedRoles[role._id])) {
          this.$delete(this.changedRoles, role._id);
        }

        return;
      }

      const nextActions = this.getNextActions(currentActions, action, value);

      if (changedRole) {
        /**
           * If we have changes for role but we don't have changes for permission
           */
        this.$set(changedRole, permission._id, nextActions);
      } else {
        /**
           * If we don't have changes for role
           */
        this.$set(this.changedRoles, role._id, { [permission._id]: nextActions });
      }
    },

    /**
     * Fetch permissions and roles lists
     *
     * @returns void
     */
    async fetchList() {
      this.pending = true;

      const [{ data: permissions }] = await Promise.all([
        this.fetchPermissionsListWithoutStore({ params: { limit: MAX_LIMIT } }),
        this.fetchRolesList({ params: { limit: MAX_LIMIT } }),
      ]);

      this.groupedPermissions = getGroupedPermissions(permissions, this.allViews, this.playlists);

      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
  .submit-button {
    position: sticky;
    bottom: 10px;
  }

  .admin-rights /deep/ {
    .v-window__container--is-active th {
      position: relative;
      top: 0;
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
</style>
