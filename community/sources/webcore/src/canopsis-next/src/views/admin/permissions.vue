<template>
  <div class="admin-rights">
    <c-page-header />
    <v-card class="ma-4 mt-0 px-4 pb-4">
      <c-progress-overlay :pending="pending" />
      <v-tabs v-model="activeTab" fixed-tabs>
        <template v-for="tab in treeviewPermissions">
          <v-tab :key="tab._id" :href="`#${tab._id}`">
            {{ $t(`permission.title.${tab.name}`) }}
          </v-tab>
          <v-tab-item :key="`${tab._id}-item`" :value="tab._id">
            <permissions-table-wrapper
              :treeview-permissions="tab.children"
              :roles="tab._id === apiTabId ? apiRoles : uiRoles"
              :search-depth="searchDepthByTabId[tab._id]"
              @input="changeRole"
            />
          </v-tab-item>
        </template>
      </v-tabs>
    </v-card>
    <v-layout
      v-show="hasChanges"
      class="sticky-bottom-buttons mt-3 gap-2"
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
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';

import { USER_PERMISSIONS_GROUPS } from '@/constants';

import {
  useRolePermissionActions,
  useRolePermissionFetching,
} from '@/components/other/permission/hooks/role-permission';

import PermissionsTableWrapper from '@/components/other/permission/permissions-table-wrapper.vue';
import PermissionsFabBtn from '@/components/other/permission/permissions-fab-btn.vue';

export default {
  components: { PermissionsTableWrapper, PermissionsFabBtn },
  setup() {
    const activeTab = ref();
    const apiTabId = USER_PERMISSIONS_GROUPS.api;
    const searchDepthByTabId = {
      [USER_PERMISSIONS_GROUPS.commonviews]: 1,
    };

    const {
      pending,
      uiRoles,
      apiRoles,
      treeviewPermissions,
      hasChanges,
      resetRolesById,
      updateRoles,
      changeRole,
      fetchList,
    } = useRolePermissionFetching();

    const { submit, cancel } = useRolePermissionActions({ updateRoles, resetRolesById });

    onMounted(fetchList);

    return {
      activeTab,

      apiTabId,
      searchDepthByTabId,

      pending,
      uiRoles,
      apiRoles,
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
