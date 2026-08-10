<template>
  <div>
    <c-page-header />

    <v-card class="ma-4 mt-0">
      <v-tabs v-model="activeTab" align-tabs="center" color="primary">
        <v-tab v-if="hasReadAnyUserAccess" :value="USERS_TABS.users">
          {{ t('common.user', 2) }}
        </v-tab>
        <v-tab v-if="hasReadAnyShareTokenAccess" :value="USERS_TABS.shareTokens">
          {{ t('common.sharedTokens') }}
        </v-tab>
      </v-tabs>

      <v-card-text>
        <v-window v-model="activeTab">
          <v-window-item :value="USERS_TABS.users">
            <users />
          </v-window-item>
          <v-window-item :value="USERS_TABS.shareTokens">
            <div class="text-medium-emphasis pa-4">
              {{ t('common.sharedTokens') }}
            </div>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <c-fab-btn :has-access="hasCreateAccess" @refresh="refresh" @create="create">
      <span>{{ t('modals.createUser.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import Users from '@/components/other/users/users.vue';
import { MODALS, USERS_TABS, USER_PERMISSIONS } from '@/constants';
import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useUser } from '@/hooks/store/modules/user';

export default {
  components: {
    Users,
  },

  setup() {
    const modals = useModals();
    const { t } = useI18n();
    const { fetchUsersListWithPreviousParams, createUser } = useUser();

    const { hasCreateAccess: hasCreateAnyUserAccess, hasReadAccess: hasReadAnyUserAccess } =
      useCRUDPermissions(USER_PERMISSIONS.technical.user);

    const { hasReadAccess: hasReadAnyShareTokenAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.shareToken,
    );

    const activeTab = ref(USERS_TABS.users);

    const hasCreateAccess = computed(
      () =>
        ({
          [USERS_TABS.users]: hasCreateAnyUserAccess.value,
          [USERS_TABS.shareTokens]: false,
        })[activeTab.value],
    );

    const showCreateUserModal = () => {
      modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await createUser({ data });

            await fetchUsersListWithPreviousParams();
          },
        },
      });
    };

    const refresh = () => {
      switch (activeTab.value) {
        case USERS_TABS.users:
          fetchUsersListWithPreviousParams();
          break;
      }
    };

    const createByActiveTab = () => {
      switch (activeTab.value) {
        case USERS_TABS.users:
          showCreateUserModal();
          break;
      }
    };

    return {
      USERS_TABS,
      activeTab,
      hasCreateAccess,
      hasReadAnyUserAccess,
      hasReadAnyShareTokenAccess,
      t,
      create: createByActiveTab,
      refresh,
    };
  },
};
</script>
