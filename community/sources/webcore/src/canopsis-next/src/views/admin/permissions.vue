<template>
  <v-container class="admin-rights">
    <c-page-header />
    <v-card class="position-relative">
      <c-progress-overlay :pending="pending" />
      <v-tabs v-model="activeTab" fixed-tabs>
        <template v-for="tab in treeviewPermissions">
          <v-tab :key="tab._id" :href="`#${tab._id}`">
            {{ $t(`permission.title.${tab.name}`) }}
          </v-tab>
          <v-tab-item :key="`${tab._id}-item`" :value="tab._id">
            <permissions-table
              :treeview-permissions="tab.children"
              :roles="roles"
              @input="changeRole"
            />
          </v-tab-item>
        </template>
      </v-tabs>
    </v-card>
    <v-layout
      v-show="hasChanges"
      class="submit-button mt-3 gap-2"
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

    <permissions-fab-btn @refresh="fetchList" />
  </v-container>
</template>

<script>
import { keyBy, omit, isEqual, filter, cloneDeep } from 'lodash';
import { computed, ref, set, onMounted } from 'vue';

import { API_USER_PERMISSIONS_ROOT_GROUPS, MAX_LIMIT, MODALS, ROLE_TYPES } from '@/constants';

import { permissionsToTreeview } from '@/helpers/entities/permissions/list';
import { formToRolePermissions, roleToPermissionForm } from '@/helpers/entities/role/form';
import { mapIds } from '@/helpers/array';

import { useModals } from '@/hooks/modals';
import { useRole } from '@/hooks/store/modules/role';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePermissions } from '@/hooks/store/modules/permissions';

import PermissionsTable from '@/components/other/permission/permissions-table.vue';
import PermissionsFabBtn from '@/components/other/permission/permissions-fab-btn.vue';

export default {
  components: { PermissionsTable, PermissionsFabBtn },
  setup() {
    const activeTab = ref();
    const originalRolesById = ref({});
    const rolesById = ref({});
    const changedRoles = ref({});
    const permissions = ref([]);

    const modals = useModals();

    const { fetchPermissionsListWithoutStore } = usePermissions();
    const { fetchRolesListWithoutStore, bulkUpdateRolePermissions } = useRole();

    const resetRolesById = () => {
      rolesById.value = cloneDeep(originalRolesById.value);
      changedRoles.value = {};
    };

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const [rolesResponse, permissionsResponse] = await Promise.all([
        fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT, with_flags: true } }),
        fetchPermissionsListWithoutStore({ params: { limit: MAX_LIMIT } }),
      ]);

      originalRolesById.value = keyBy(Object.values(rolesResponse.data).map(roleToPermissionForm), '_id');
      permissions.value = permissionsResponse.data;
      resetRolesById();
    });

    const treeviewPermissions = computed(() => permissionsToTreeview(permissions.value));
    const isApiPermissionsTab = computed(() => API_USER_PERMISSIONS_ROOT_GROUPS.includes(activeTab.value));
    const uiRoles = computed(() => filter(Object.values(rolesById.value), ['type', ROLE_TYPES.ui]));
    const apiRoles = computed(() => filter(Object.values(rolesById.value), ['type', ROLE_TYPES.api]));
    const roles = computed(() => (isApiPermissionsTab.value ? apiRoles.value : uiRoles.value));
    const hasChanges = computed(() => (
      Object.entries(changedRoles.value).some(([roleId, rolePermissions]) => (
        Object.keys(rolePermissions).some(permissionId => (
          !isEqual(
            rolesById.value[roleId].permissions[permissionId],
            originalRolesById.value[roleId].permissions[permissionId],
          )
        ))
      ))
    ));

    const setChangedRolePermission = (roleId, permissionId) => {
      if (!changedRoles.value[roleId]) {
        set(changedRoles.value, roleId, { [permissionId]: true });
      }

      if (!changedRoles.value[roleId][permissionId]) {
        set(changedRoles.value[roleId], permissionId, true);
      }
    };

    const setChangedRolePermissions = (roleId, permissionsIds) => (
      permissionsIds.forEach(permissionId => setChangedRolePermission(roleId, permissionId))
    );

    const changeRole = (value, role, permission, action) => {
      if (!changedRoles.value[role._id]) {
        set(changedRoles.value, role._id, { [permission._id]: true });
      }

      if (!changedRoles.value[role._id][permission._id]) {
        set(changedRoles.value, permission._id, true);
      }

      if (permission.allChildren) {
        const allChildrenIds = mapIds(permission.allChildren);
        const newPermissions = value
          ? {
            ...role.permissions,
            ...permission.allChildren.reduce((acc, { _id: id, actions }) => {
              acc[id] = actions;

              return acc;
            }, {}),
          }
          : omit(role.permissions, allChildrenIds);

        setChangedRolePermissions(role._id, allChildrenIds);

        set(rolesById.value[role._id], 'permissions', newPermissions);

        return;
      }

      setChangedRolePermission(role._id, permission._id);

      const currentActions = role.permissions[permission._id] ?? [];

      const newActions = value
        ? [...currentActions, action]
        : currentActions.filter(currentAction => currentAction !== action);

      if (!newActions.length) {
        set(rolesById.value[role._id], 'permissions', omit(role.permissions, permission._id));

        return;
      }

      set(rolesById.value[role._id].permissions, permission._id, newActions);
    };

    const updateRoles = async () => {
      const rolesForUpdate = Object.keys(changedRoles.value)
        .map(roleId => rolesById.value[roleId] && formToRolePermissions(rolesById.value[roleId]))
        .filter(Boolean);

      await bulkUpdateRolePermissions({ data: rolesForUpdate });

      return fetchList();
    };

    const submit = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: updateRoles,
      },
    });

    const cancel = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: resetRolesById,
      },
    });

    onMounted(fetchList);

    return {
      activeTab,
      pending,
      roles,
      treeviewPermissions,
      hasChanges,
      changeRole,
      fetchList,
      submit,
      cancel,
    };
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
</style>
