<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <v-tabs
        v-model="activeTab"
        slider-color="primary"
        centered
      >
        <v-tab
          v-if="hasReadAnyUserAccess"
          :href="`#${USERS_TABS.users}`"
        >
          {{ $tc('common.user', 2) }}
        </v-tab>
        <v-tab
          v-if="hasReadAnyShareTokenAccess"
          :href="`#${USERS_TABS.shareTokens}`"
        >
          {{ $t('common.sharedTokens') }}
        </v-tab>
      </v-tabs>
      <v-tabs-items v-model="activeTab">
        <v-card-text>
          <v-tab-item :value="USERS_TABS.users">
            <users />
          </v-tab-item>
          <v-tab-item :value="USERS_TABS.shareTokens">
            <share-tokens />
          </v-tab-item>
        </v-card-text>
      </v-tabs-items>
    </v-card>
    <c-fab-btn
      :has-access="hasCreateAccess"
      @refresh="refresh"
      @create="create"
    >
      <span>{{ $t('modals.createUser.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS, USERS_TABS, USER_PERMISSIONS } from '@/constants';

import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useStoreModuleHooks } from '@/hooks/store';
import { useUser } from '@/hooks/store/modules/user';

import Users from '@/components/other/users/users.vue';
import ShareTokens from '@/components/other/share-token/share-tokens.vue';

export default {
  components: {
    Users,
    ShareTokens,
  },
  setup() {
    const modals = useModals();

    const {
      fetchUsersListWithPreviousParams,
      createUserWithPopup,
    } = useUser();

    const {
      hasCreateAccess: hasCreateAnyUserAccess,
      hasReadAccess: hasReadAnyUserAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.user);

    const {
      hasReadAccess: hasReadAnyShareTokenAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.shareToken);

    const { useActions: useShareTokenActions } = useStoreModuleHooks('shareToken');
    const { fetchShareTokensListWithPreviousParams } = useShareTokenActions({
      fetchShareTokensListWithPreviousParams: 'fetchListWithPreviousParams',
    });

    const activeTab = ref(USERS_TABS.users);

    const hasCreateAccess = computed(() => ({
      [USERS_TABS.users]: hasCreateAnyUserAccess.value,
      [USERS_TABS.shareTokens]: false,
    }[activeTab.value]));

    /**
     * Shows the modal for creating a new user.
     * After successful creation, refreshes the users list.
     */
    const showCreateUserModal = () => {
      modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await createUserWithPopup({ data });

            await fetchUsersListWithPreviousParams();
          },
        },
      });
    };

    /**
     * Refreshes the list based on the active tab.
     */
    const refresh = () => {
      switch (activeTab.value) {
        case USERS_TABS.users:
          fetchUsersListWithPreviousParams();
          break;
        case USERS_TABS.shareTokens:
          fetchShareTokensListWithPreviousParams();
          break;
      }
    };

    /**
     * Shows the create modal based on the active tab.
     */
    const create = () => {
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
      refresh,
      create,
    };
  },
};
</script>
