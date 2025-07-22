<template>
  <v-menu
    :disabled="pending"
    :min-width="300"
    :close-on-content-click="false"
    bottom
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        text
        v-on="on"
      >
        {{ $tc('common.notification', 2) }}
        <c-circle-badge
          :class="badgeTextClasses"
          :outlined="!totalCount"
          class="ml-2"
          color="white"
        >
          <v-progress-circular
            v-if="pending"
            size="12"
            width="1"
            indeterminate
          />
          <span v-else>
            {{ totalCount }}
          </span>
        </c-circle-badge>
      </v-btn>
    </template>
    <v-list class="py-0">
      <top-bar-notifications-menu-item
        v-for="notification in notifications"
        :key="notification._id"
        :notification="notification"
      />
      <v-list-item v-if="!notifications.length">
        <v-list-item-content>
          <v-list-item-title class="text-center">
            <strong class="font-italic grey--text">{{ $t('notifications.topBar.noNotifications') }}</strong>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-layout class="pa-4" justify-center>
        <v-btn
          :to="linkToSeeAllNotifications"
          color="primary"
        >
          {{ $t('notifications.topBar.seeAll') }}
        </v-btn>
      </v-layout>
    </v-list>
  </v-menu>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { DEFAULT_NOTIFICATION_TOP_BAR_LIMIT, NOTIFICATIONS_PAGE_TABS_KEYS, ROUTES_NAMES } from '@/constants';

import { useSocketRoom } from '@/hooks/socket';
import { useInfo } from '@/hooks/store/modules/info';
import { useNotifications } from '@/hooks/store/modules/notification';
import { usePendingHandler } from '@/hooks/query/pending';

import TopBarNotificationsMenuItem from './top-bar-notifications-menu-item.vue';

export default {
  components: {
    TopBarNotificationsMenuItem,
  },
  setup() {
    const { notificationDisplayCount } = useInfo();
    const { fetchNotificationsListWithoutStore } = useNotifications();

    const linkToSeeAllNotifications = {
      name: ROUTES_NAMES.notifications,
      params: { tabId: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove },
    };

    const notifications = ref([]);
    const totalCount = ref(0);

    const badgeTextClasses = computed(() => ({
      'primary--text': !!totalCount.value,
      'white--text': !totalCount.value,
    }));

    /**
     * Fetches notifications list from the API and updates the component state
     */
    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const { data, meta } = await fetchNotificationsListWithoutStore({
        params: {
          page: 1,
          limit: notificationDisplayCount.value || DEFAULT_NOTIFICATION_TOP_BAR_LIMIT,
        },
      });

      notifications.value = data ?? [];
      totalCount.value = meta?.total_count ?? 0;
    });

    /**
     * Updates notifications data and total count from socket events
     *
     * @param {Object} params - The update parameters
     * @param {Array} [params.data] - Array of notification objects
     * @param {number} [params.total_count] - Total count of notifications
     */
    const updateNotifications = ({ data, total_count: newTotalCount } = {}) => {
      notifications.value = data ?? [];
      totalCount.value = newTotalCount ?? 0;
    };

    useSocketRoom(SOCKET_ROOMS.notifications, updateNotifications);

    onMounted(fetchList);

    return {
      linkToSeeAllNotifications,
      notifications,
      pending,
      totalCount,
      badgeTextClasses,
    };
  },
};
</script>
