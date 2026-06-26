<template>
  <v-app-bar
    :height="TOP_BAR_HEIGHT"
    :class="{ 'top-bar--with-extension': isShownGroupsTopBar }"
    :style="topBarStyle"
    class="top-bar bg-primary"
    app
  >
    <v-app-bar-nav-icon
      v-if="isShownGroupsSideBar && !$route.meta.simpleNavigation"
      class="ml-0 text-white"
      size="large"
      @click="$emit('toggle:sideBar')"
    />

    <c-layout
      v-else
      class="top-bar__brand fill-height"
      align="center"
    >
      <app-logo class="canopsis-logo mr-2" />

      <c-layout
        class="version ml-1 fill-height"
        align="end"
      >
        <logged-users-count badge-color="secondary" />
        <app-version />
      </c-layout>
    </c-layout>

    <top-bar-title :title="appTitle" />

    <healthcheck-chips-list
      v-if="isProVersion && hasAccessToHealthcheckStatus"
      class="top-bar__healthcheck-chips"
    />

    <v-spacer />

    <c-layout
      v-if="!$route.meta.simpleNavigation"
      class="top-bar__menus fill-height"
      align="center"
      gap="0"
    >
      <top-bar-exploitation-menu />
      <top-bar-administration-menu />
      <top-bar-notifications-menu />
      <top-bar-user-menu />
    </c-layout>

    <template
      v-if="isShownGroupsTopBar"
      #extension
    >
      <groups-top-bar />
    </template>
  </v-app-bar>
</template>

<script>
import { computed } from 'vue'
import { useDisplay } from 'vuetify'

import { MEDIA_QUERIES_BREAKPOINTS, SIDE_BAR_WIDTH, TOP_BAR_HEIGHT } from '@/config'
import { USER_PERMISSIONS } from '@/constants'

import { useAuth, useCanPermission } from '@/hooks/auth'
import { useInfo } from '@/hooks/store/modules/info'

import HealthcheckChipsList from '@/components/other/healthcheck/partials/healthcheck-chips-list.vue'

import AppLogo from './app-logo.vue'
import AppVersion from './app-version.vue'
import LoggedUsersCount from './logged-users-count.vue'
import GroupsTopBar from './groups-top-bar/groups-top-bar.vue'
import TopBarExploitationMenu from './top-bar-exploitation-menu.vue'
import TopBarAdministrationMenu from './top-bar-administration-menu.vue'
import TopBarNotificationsMenu from './top-bar-notifications-menu.vue'
import TopBarUserMenu from './top-bar-user-menu.vue'
import TopBarTitle from './top-bar-title.vue'

export default {
  components: {
    AppLogo,
    AppVersion,
    HealthcheckChipsList,
    LoggedUsersCount,
    GroupsTopBar,
    TopBarExploitationMenu,
    TopBarAdministrationMenu,
    TopBarNotificationsMenu,
    TopBarUserMenu,
    TopBarTitle,
  },

  props: {
    sideBarOpen: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['toggle:sideBar'],

  setup(props) {
    const { isShownGroupsTopBar, isShownGroupsSideBar } = useAuth()
    const { isProVersion, appTitle } = useInfo()
    const { hasAccess: hasAccessToHealthcheckStatus } = useCanPermission(USER_PERMISSIONS.technical.healthcheckStatus)
    const { width } = useDisplay()

    const topBarStyle = computed(() => {
      if (!props.sideBarOpen || width.value < MEDIA_QUERIES_BREAKPOINTS.t) {
        return undefined
      }

      return {
        left: `${SIDE_BAR_WIDTH}px`,
        width: `calc(100% - ${SIDE_BAR_WIDTH}px)`,
      }
    })

    return {
      TOP_BAR_HEIGHT,
      appTitle,
      hasAccessToHealthcheckStatus,
      isProVersion,
      isShownGroupsSideBar,
      isShownGroupsTopBar,
      topBarStyle,
    }
  },
}
</script>

<style lang="scss" scoped>
.canopsis-logo {
  max-height: 80%;
  margin-left: 1em;
}

.version {
  color: white;
  font-size: 0.7em;
  position: relative;

  :deep(.logged-users-count) {
    left: -8px;
  }
}

.top-bar {
  &.top-bar--with-extension :deep(.v-toolbar__content) {
    z-index: 1;
  }

  :deep(.v-toolbar__content),
  :deep(.v-toolbar__extension) {
    padding: 0;
  }

  :deep(.v-toolbar__content) {
    height: 48px !important;
  }

  &__brand {
    flex: 0 0 auto;
    min-width: 250px;
  }

  &__healthcheck-chips {
    flex: 1 1 auto;
    min-width: 0;
  }

  &__menus {
    flex: 0 0 auto;
    height: 48px;
    min-width: 0;

    :deep(.v-btn) {
      align-self: stretch;
    }
  }
}
</style>
