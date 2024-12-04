<template>
  <v-container class="admin-rights">
    <c-page-header />
    <v-card class="position-relative">
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
import { keyBy } from 'lodash';
import { computed, ref, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { permissionsToTreeview } from '@/helpers/entities/permissions/list';
import { roleToForm } from '@/helpers/entities/role/form';

import { useRole } from '@/hooks/store/modules/role';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePermissions } from '@/hooks/store/modules/permissions';

import PermissionsTable from '@/components/other/permission2/permissions-table.vue';

export default {
  components: { PermissionsTable },
  setup() {
    const originalRoles = ref([]);
    const rolesById = ref({});
    const permissions = ref([]);

    const { fetchPermissionsListWithoutStore } = usePermissions();
    const { fetchRolesListWithoutStore } = useRole();

    const resetRolesById = () => {
      rolesById.value = keyBy(originalRoles.value.map(roleToForm), '_id');
    };

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const [rolesResponse, permissionsResponse] = await Promise.all([
        fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT, with_flags: true } }),
        fetchPermissionsListWithoutStore({ params: { limit: MAX_LIMIT } }),
      ]);

      originalRoles.value = rolesResponse.data;
      permissions.value = permissionsResponse.data;
      resetRolesById();
    });

    const treeviewPermissions = computed(() => permissionsToTreeview(permissions.value));
    const roles = computed(() => Object.values(rolesById.value));

    const changeRole = () => {
    };

    onMounted(fetchList);

    return {
      pending,
      roles,
      treeviewPermissions,
      changeRole,
    };
  },
};
</script>
