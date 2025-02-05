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
import { ref, onMounted } from 'vue';

import {
  useRolePermissionActions,
  useRolePermissionFetching,
} from '@/components/other/permission/hooks/role-permission';

import PermissionsTable from '@/components/other/permission/permissions-table.vue';
import PermissionsFabBtn from '@/components/other/permission/permissions-fab-btn.vue';

export default {
  components: { PermissionsTable, PermissionsFabBtn },
  setup() {
    const activeTab = ref();

    const {
      pending,
      roles,
      treeviewPermissions,
      hasChanges,
      resetRolesById,
      updateRoles,
      changeRole,
      fetchList,
    } = useRolePermissionFetching({ activeTab });

    const { submit, cancel } = useRolePermissionActions({ updateRoles, resetRolesById });

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
