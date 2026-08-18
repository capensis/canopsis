<template>
  <div class="admin-rights">
    <c-page-header />

    <v-card class="admin-rights__card ma-4 mt-0 px-4 pb-4">
      <c-progress-overlay :pending="pending" />

      <v-tabs
        v-model="activeTab"
        align-tabs="center"
        color="primary"
        fixed-tabs
      >
        <v-tab v-for="tab in treeviewPermissions" :key="tab._id" :value="tab._id">
          {{ $t(`permission.title.${tab.name}`) }}
        </v-tab>
      </v-tabs>

      <v-window v-model="activeTab" class="admin-rights__window">
        <v-window-item v-for="tab in treeviewPermissions" :key="`${tab._id}-item`" :value="tab._id">
          <permissions-table-wrapper
            :roles="tab._id === apiTabId ? apiRoles : uiRoles"
            :search-depth="searchDepthByTabId[tab._id]"
            :treeview-permissions="tab.children"
            @update:permission="changeRole"
          />
        </v-window-item>
      </v-window>
    </v-card>

    <c-layout v-show="hasChanges" class="sticky-bottom-buttons mt-3" gap="2">
      <v-btn class="ml-3" color="primary" @click="submit">
        {{ $t('common.submit') }}
      </v-btn>

      <v-btn @click="cancel">
        {{ $t('common.cancel') }}
      </v-btn>
    </c-layout>

    <permissions-fab-btn @refresh="fetchList" />
  </div>
</template>

<script>
import { onMounted, ref, watch } from 'vue';

import { USER_PERMISSIONS_GROUPS } from '@/constants';

import {
  useRolePermissionActions,
  useRolePermissionFetching,
} from '@/components/other/permission/hooks/role-permission';

import PermissionsFabBtn from '@/components/other/permission/permissions-fab-btn.vue';
import PermissionsTableWrapper from '@/components/other/permission/permissions-table-wrapper.vue';

export default {
  components: { PermissionsTableWrapper, PermissionsFabBtn },

  setup() {
    const activeTab = ref();
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
      apiTabId,
    } = useRolePermissionFetching();

    const { submit, cancel } = useRolePermissionActions({ updateRoles, resetRolesById });

    watch(treeviewPermissions, ([firstTab]) => {
      if (!activeTab.value && firstTab) {
        activeTab.value = firstTab._id;
      }
    });

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

<style lang="scss" scoped>
.admin-rights {
  &__card,
  &__window {
    overflow: visible;
  }

  &__window {
    :deep(.v-window__container),
    :deep(.v-window-item) {
      overflow: visible;
    }
  }
}
</style>
