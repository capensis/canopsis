<template>
  <v-container class="admin-rights">
    <c-page-header />
    <v-card class="position-relative">
      <v-fade-transition>
        <v-layout
          v-show="pending"
          class="progress"
          column
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
      </v-fade-transition>
      <v-tabs
        v-if="hasReadAnyRoleAccess"
        slider-color="primary"
        fixed-tabs
      >
        <template v-for="(groupPermissions, groupKey) in preparedPermissionsGroups">
          <v-tab :key="`tab-${groupKey}`">
            {{ groupKey }}
          </v-tab>
          <v-tab-item :key="`tab-item-${groupKey}`">
            <permissions-table-wrapper
              :permissions="groupPermissions"
              :roles="preparedRoles"
              :changed-roles="changedRoles"
              :disabled="!hasUpdateAnyPermissionAccess"
              @change="changeCheckboxValue"
            />
          </v-tab-item>
        </template>
      </v-tabs>
    </v-card>
    <v-layout
      v-show="hasChanges"
      class="submit-button mt-3"
    >
      <v-btn
        class="ml-3"
        color="primary"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
      <v-btn @click="cancel">
        {{ $t('common.cancel') }}
      </v-btn>
    </v-layout>
    <permissions-fab-buttons @refresh="fetchList" />
  </v-container>
</template>

<script>
import { get, isEmpty, sortBy } from 'lodash';

import { MAX_LIMIT, MODALS } from '@/constants';

import { getGroupedPermissions } from '@/helpers/entities/permissions/list';
import { roleToForm, formToRole } from '@/helpers/entities/role/form';

import { authMixin } from '@/mixins/auth';
import { entitiesPermissionsMixin } from '@/mixins/entities/permission';
import { entitiesRoleMixin } from '@/mixins/entities/role';
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
    permissionsTechnicalRoleMixin,
    permissionsTechnicalPermissionMixin,
  ],
  data() {
    return {
      pending: false,
      permissions: [],
      changedRoles: {},
    };
  },
  computed: {
    hasChanges() {
      return !isEmpty(this.changedRoles);
    },

    groupedPermissions() {
      return getGroupedPermissions(this.permissions);
    },

    businessPermissionsGroup() {
      const { business } = this.groupedPermissions;

      return Object.entries(business).map(([key, groupPermissions]) => ({
        name: this.$tc(`permission.business.${key}`),
        permissions: this.prepareGroupPermissions(groupPermissions),
      }));
    },

    technicalPermissionsGroup() {
      const { technical } = this.groupedPermissions;

      return Object.entries(technical).map(([key, groupPermissions]) => ({
        name: this.$tc(`permission.technical.${key}`),
        permissions: this.prepareGroupPermissions(groupPermissions),
      }));
    },

    apiPermissionsGroup() {
      const { api } = this.groupedPermissions;

      return Object.entries(api).map(([key, groupPermissions]) => ({
        name: this.$tc(`permission.api.${key}`),
        permissions: this.prepareGroupPermissions(groupPermissions),
      }));
    },

    viewPermissionsGroup() {
      const { view, playlist } = this.groupedPermissions;

      const viewsPermissionsByGroupTitle = view.reduce((acc, permission) => {
        const { view_group: viewGroup, ...rest } = permission;

        if (!acc[viewGroup._id]) {
          acc[viewGroup._id] = {
            viewGroup,
            permissions: [],
          };
        }

        acc[viewGroup._id].permissions.push(rest);

        return acc;
      }, {});

      const viewPermissionsGroup = sortBy(Object.values(viewsPermissionsByGroupTitle), ['viewGroup.position'])
        .map(({ viewGroup, permissions: viewGroupPermissions }) => ({
          name: viewGroup.title,
          permissions: sortBy(viewGroupPermissions, ['view.position']).map(({ description, name, ...permission }) => ({
            name: description,
            ...permission,
          })),
        }));

      viewPermissionsGroup.push({
        name: this.$tc('common.playlist'),
        permissions: sortBy(playlist, ['playlist.name']).map(({ description, name, ...permission }) => ({
          name: description,
          ...permission,
        })),
      });

      return viewPermissionsGroup;
    },

    preparedPermissionsGroups() {
      return {
        business: this.businessPermissionsGroup,
        view: this.viewPermissionsGroup,
        technical: this.technicalPermissionsGroup,
        api: this.apiPermissionsGroup,
      };
    },

    preparedRoles() {
      return sortBy(this.roles, [({ _id: name }) => name.toLowerCase()])
        .map(role => ({
          ...roleToForm(role),

          editable: role.editable,
          deletable: role.deletable,
        }));
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    prepareGroupPermissions(groupPermissions) {
      const preparedPermissions = groupPermissions.map((permission) => {
        const messageKey = `permission.permissions.${permission._id}`;
        const { name, description } = this.$te(messageKey) ? this.$t(messageKey) : {};

        return ({
          ...permission,
          name: name ?? permission.description,
          description: description ?? '',
        });
      });

      /**
       * We are using order which one we've defined on the reduce accumulator initial value.
       * For not `number`/`number string` object keys ordering is staying like we define
       */
      return sortBy(preparedPermissions, ['name']);
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
     * Send request for update changedRoles
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

        await this.fetchRolesList({ params: { limit: MAX_LIMIT, with_flags: true } });

        this.$popups.success({ text: this.$t('success.default') });

        this.clearChangedRoles();
      } catch (err) {
        console.error(err);

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
        this.fetchRolesList({ params: { limit: MAX_LIMIT, with_flags: true } }),
      ]);

      this.permissions = permissions;

      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
  .submit-button {
    position: sticky;
    bottom: 10px;
    pointer-events: none;

    button {
      pointer-events: all;
    }
  }

  .admin-rights ::v-deep {
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

    & ::v-deep .v-progress-circular {
      top: 50%;
      left: 50%;
      margin-top: -16px;
      margin-left: -16px;
    }
  }
</style>
