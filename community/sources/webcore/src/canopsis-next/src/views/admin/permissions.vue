<template>
  <v-container class="admin-rights">
    <c-page-header />
    <v-card class="position-relative">
      <h1 v-if="hasChanges">
        HAS CHANGES
      </h1>
      <c-progress-overlay :pending="pending" />
      <v-tabs fixed-tabs>
        <template v-for="tab in treeviewPermissions">
          <v-tab :key="tab._id">
            {{ tab._id }}
          </v-tab>
          <v-tab-item :key="`${tab._id}-item`">
            <permissions-table
              :treeview-permissions="treeviewPermissions"
              :roles="roles"
              @input="changeRole"
            />
          </v-tab-item>
        </template>
      </v-tabs>
    </v-card>
  </v-container>
</template>

<script>
import { keyBy, omit, isEqual } from 'lodash';
import { computed, ref, set, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { permissionsToTreeview } from '@/helpers/entities/permissions/list';
import { roleToForm } from '@/helpers/entities/role/form';
import { mapIds } from '@/helpers/array';

import { useRole } from '@/hooks/store/modules/role';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePermissions } from '@/hooks/store/modules/permissions';

import PermissionsTable from '@/components/other/permission/permissions-table.vue';

export default {
  components: { PermissionsTable },
  setup() {
    const originalRolesById = ref([]);
    const rolesById = ref({});
    const changedRoles = ref({});
    const permissions = ref([]);

    const { fetchPermissionsListWithoutStore } = usePermissions();
    const { fetchRolesListWithoutStore } = useRole();

    const resetRolesById = () => {
      rolesById.value = keyBy(Object.values(originalRolesById.value).map(roleToForm), '_id');
      changedRoles.value = {};
    };

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const [rolesResponse, permissionsResponse] = await Promise.all([
        fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT, with_flags: true } }),
        fetchPermissionsListWithoutStore({ params: { limit: MAX_LIMIT } }),
      ]);

      originalRolesById.value = keyBy(rolesResponse.data, '_id');
      permissions.value = permissionsResponse.data;
      resetRolesById();
    });

    const treeviewPermissions = computed(() => permissionsToTreeview(permissions.value));
    const roles = computed(() => Object.values(rolesById.value));
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

    onMounted(fetchList);

    return {
      pending,
      roles,
      treeviewPermissions,
      hasChanges,
      changeRole,
    };
  },
};
</script>
