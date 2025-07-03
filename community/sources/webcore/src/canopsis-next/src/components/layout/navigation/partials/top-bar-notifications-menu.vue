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
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { DEFAULT_NOTIFICATION_TOP_BAR_LIMIT, NOTIFICATIONS_PAGE_TABS_KEYS, ROUTES_NAMES } from '@/constants';

import { useComponentInstance } from '@/hooks/vue';
import { useInfo } from '@/hooks/store/modules/info';
import { useNotifications } from '@/hooks/store/modules/notification';
import { usePendingHandler } from '@/hooks/query/pending';

import TopBarNotificationsMenuItem from './top-bar-notifications-menu-item.vue';

export default {
  components: {
    TopBarNotificationsMenuItem,
  },
  setup() {
    const instance = useComponentInstance();
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

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const { data, meta } = await fetchNotificationsListWithoutStore({
        params: {
          page: 1,
          limit: notificationDisplayCount.value ?? DEFAULT_NOTIFICATION_TOP_BAR_LIMIT,
        },
      });

      notifications.value = data ?? [];
      totalCount.value = meta?.total_count ?? 0;
    });

    const updateNotifications = ({ data, total_count: newTotalCount } = {}) => {
      notifications.value = data ?? [];
      totalCount.value = newTotalCount ?? 0;
    };

    const joinNotificationsRoom = () => {
      if (instance?.$socket) {
        instance.$socket.join(SOCKET_ROOMS.notifications).addListener(updateNotifications);
      }
    };

    const leaveNotificationsRoom = () => {
      if (instance?.$socket) {
        instance.$socket.leave(SOCKET_ROOMS.notifications).removeListener(updateNotifications);
      }
    };

    onMounted(() => {
      joinNotificationsRoom();
      fetchList();
    });

    onBeforeUnmount(leaveNotificationsRoom);

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
