<template>
  <v-list-item>
    <v-list-item-content>
      <v-list-item-title class="mb-2">
        <strong>{{ typeMessage }}</strong>
      </v-list-item-title>
      <v-list-item-title class="mb-2">
        <router-link :to="link">
          <strong class="blue--text">{{ notification.rule?.name }}</strong>
        </router-link>
      </v-list-item-title>
      <v-list-item-subtitle v-if="notification.author?.display_name" class="mb-2">
        <strong class="grey--text">By: {{ notification.author.display_name }}</strong>
      </v-list-item-subtitle>
      <v-list-item-subtitle v-if="notification.message">
        <p>{{ notification.message }}</p>
      </v-list-item-subtitle>
      <v-divider />
    </v-list-item-content>
  </v-list-item>
</template>

<script>
import { computed } from 'vue';

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
    const { t } = useI18n();

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
      link,
      typeMessage,
    };
  },
};
</script>
