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
          :href="`#${$constants.USERS_TABS.users}`"
        >
          {{ $tc('common.user', 2) }}
        </v-tab>
        <v-tab
          v-if="hasReadAnyShareTokenAccess"
          :href="`#${$constants.USERS_TABS.shareTokens}`"
        >
          {{ $t('common.sharedTokens') }}
        </v-tab>
      </v-tabs>
      <v-tabs-items v-model="activeTab">
        <v-card-text>
          <v-tab-item :value="$constants.USERS_TABS.users">
            <users />
          </v-tab-item>
          <v-tab-item :value="$constants.USERS_TABS.shareTokens">
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
import { MODALS, USERS_TABS } from '@/constants';

import { entitiesUserMixin } from '@/mixins/entities/user';
import { entitiesShareTokenMixin } from '@/mixins/entities/share-token';
import { permissionsTechnicalUserMixin } from '@/mixins/permissions/technical/user';
import { permissionsTechnicalShareTokenMixin } from '@/mixins/permissions/technical/share-token';

import Users from '@/components/other/users/users.vue';
import ShareTokens from '@/components/other/share-token/share-tokens.vue';

export default {
  components: {
    Users,
    ShareTokens,
  },
  mixins: [
    entitiesUserMixin,
    entitiesShareTokenMixin,
    permissionsTechnicalUserMixin,
    permissionsTechnicalShareTokenMixin,
  ],
  data() {
    return {
      activeTab: USERS_TABS.users,
    };
  },
  computed: {
    hasCreateAccess() {
      return {
        [USERS_TABS.users]: this.hasCreateAnyUserAccess,
        [USERS_TABS.shareTokens]: false,
      }[this.activeTab];
    },
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case USERS_TABS.users:
          this.fetchUsersListWithPreviousParams();
          break;
        case USERS_TABS.shareTokens:
          this.fetchShareTokensListWithPreviousParams();
          break;
      }
    },

    create() {
      switch (this.activeTab) {
        case USERS_TABS.users:
          this.showCreateUserModal();
          break;
      }
    },

    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await this.createUserWithPopup({ data });

            await this.fetchUsersListWithPreviousParams();
          },
        },
      });
    },
  },
};
</script>
