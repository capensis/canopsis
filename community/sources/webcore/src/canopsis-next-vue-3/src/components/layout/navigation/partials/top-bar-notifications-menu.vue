<template>
  <v-menu
    :disabled="pending"
    :min-width="300"
    :close-on-content-click="false"
    origin="top center"
  >
    <template #activator="{ props: activatorProps }">
      <v-btn
        class="top-bar-notifications-menu__button text-white"
        variant="text"
        v-bind="activatorProps"
      >
        {{ t('common.notification', 2) }}

        <c-circle-badge
          class="top-bar-notifications-menu__badge ml-2 text-primary"
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
        <v-list-item-title class="text-center">
          <strong class="font-italic text-medium-emphasis">
            {{ t('notifications.topBar.noNotifications') }}
          </strong>
        </v-list-item-title>
      </v-list-item>

      <c-layout
        class="pa-4"
        justify="center"
      >
        <v-btn
          :disabled="!canOpenNotifications"
          color="primary"
          @click="showAllNotifications"
        >
          {{ t('notifications.topBar.seeAll') }}
        </v-btn>
      </c-layout>
    </v-list>
  </v-menu>
</template>

<script>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { SOCKET_ROOMS } from '@/config'
import { DEFAULT_NOTIFICATION_TOP_BAR_LIMIT, NOTIFICATIONS_PAGE_TABS_KEYS, ROUTES_NAMES } from '@/constants'

import { useI18n } from '@/hooks/i18n'
import { useSocketRoom } from '@/hooks/socket'
import { useInfo } from '@/hooks/store/modules/info'
import { useNotifications } from '@/hooks/store/modules/notification'
import { usePendingHandler } from '@/hooks/query/pending'

import TopBarNotificationsMenuItem from './top-bar-notifications-menu-item.vue'

export default {
  components: { TopBarNotificationsMenuItem },

  setup() {
    const router = useRouter()
    const { t } = useI18n()
    const { notificationDisplayCount } = useInfo()
    const { fetchNotificationsListWithoutStore } = useNotifications()

    const linkToSeeAllNotifications = {
      name: ROUTES_NAMES.notifications,
      params: { tabId: NOTIFICATIONS_PAGE_TABS_KEYS.instructionsToApprove },
    }

    const notifications = ref([])
    const totalCount = ref(0)

    const canOpenNotifications = computed(() => router.hasRoute(ROUTES_NAMES.notifications))

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const { data, meta } = await fetchNotificationsListWithoutStore({
        params: {
          page: 1,
          limit: notificationDisplayCount.value || DEFAULT_NOTIFICATION_TOP_BAR_LIMIT,
        },
      })

      notifications.value = data ?? []
      totalCount.value = meta?.total_count ?? 0
    })

    const updateNotifications = ({ data, total_count: newTotalCount } = {}) => {
      notifications.value = data ?? []
      totalCount.value = newTotalCount ?? 0
    }

    const showAllNotifications = () => {
      if (canOpenNotifications.value) {
        router.push(linkToSeeAllNotifications)
      }
    }

    useSocketRoom({ room: SOCKET_ROOMS.notifications, listener: updateNotifications })

    onMounted(fetchList)

    return {
      t,
      notifications,
      pending,
      totalCount,
      canOpenNotifications,
      showAllNotifications,
    }
  },
}
</script>

<style lang="scss" scoped>
.top-bar-notifications-menu {
  &__button {
    --v-btn-height: 48px;

    border-radius: 0;
    font-size: 14px;
    font-weight: 500;
    height: 48px;
    letter-spacing: 0;
    line-height: normal;
    min-height: 0;
    padding: 0 16px;
    text-transform: uppercase;

    :deep(.v-btn__content) {
      height: 100%;
    }
  }

  &__badge {
    background: white;
    color: rgb(var(--v-theme-primary));
    font-weight: 700;
  }
}
</style>
