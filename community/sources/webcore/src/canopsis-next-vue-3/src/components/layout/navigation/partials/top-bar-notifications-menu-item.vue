<template>
  <v-list-item>
    <div class="top-bar-notifications-menu-item">
      <v-list-item-title class="mb-2">
        <strong>{{ typeMessage }}</strong>
      </v-list-item-title>

      <v-list-item-title class="mb-2">
        <router-link v-if="canOpenNotification" :to="link">
          <strong class="text-primary">
            {{ notification.rule?.name }}
          </strong>
        </router-link>
        <strong v-else class="text-primary">
          {{ notification.rule?.name }}
        </strong>
      </v-list-item-title>

      <v-list-item-subtitle v-if="notification.author?.display_name" class="mb-2">
        <strong>By: {{ notification.author.display_name }}</strong>
      </v-list-item-subtitle>

      <div v-if="notification.comment" class="mb-2">
        <p>{{ notification.comment }}</p>
      </div>

      <v-divider />
    </div>
  </v-list-item>
</template>

<script>
import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { NOTIFICATIONS_PAGE_TABS_KEYS_BY_TYPE, ROUTES_NAMES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    notification: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const router = useRouter();
    const { t } = useI18n();

    const canOpenNotification = computed(() => router.hasRoute(ROUTES_NAMES.notifications));

    const link = computed(() => ({
      name: ROUTES_NAMES.notifications,
      params: {
        tabId: NOTIFICATIONS_PAGE_TABS_KEYS_BY_TYPE[props.notification.type],
      },
      query: {
        id: props.notification?.rule?._id,
      },
    }));

    const typeMessage = computed(() => t(`notifications.topBar.types.${props.notification.type}`));

    return {
      canOpenNotification,
      link,
      typeMessage,
    };
  },
};
</script>

<style lang="scss" scoped>
.top-bar-notifications-menu-item {
  width: 100%;
}
</style>
