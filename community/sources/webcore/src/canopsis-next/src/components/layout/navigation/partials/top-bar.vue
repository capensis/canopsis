<template>
  <v-app-bar
    :height="$config.TOP_BAR_HEIGHT"
    class="top-bar primary"
    dense
    fixed
    app
  >
    <v-app-bar-nav-icon
      v-if="isShownGroupsSideBar && !$route.meta.simpleNavigation"
      class="ml-0 white--text"
      @click="$emit('toggle:sideBar')"
    />
    <v-layout
      v-else
      fill-height
      align-center
    >
      <app-logo class="canopsis-logo mr-2" />
      <v-layout
        class="version ml-1"
        fill-height
        align-end
      >
        <logged-users-count badge-color="secondary" />
        <app-version />
      </v-layout>
    </v-layout>
    <top-bar-title :title="appTitle" />
    <healthcheck-chips-list v-if="isProVersion && hasAccessToHealthcheckStatus" />
    <v-spacer v-else />
    <portal-target :name="$constants.PORTALS_NAMES.additionalTopBarItems" />
    <v-toolbar-items v-if="!$route.meta.simpleNavigation">
      <top-bar-exploitation-menu />
      <top-bar-administration-menu />
      <top-bar-notifications-menu />
      <top-bar-user-menu />
    </v-toolbar-items>
    <template
      v-if="isShownGroupsTopBar"
      #extension=""
    >
      <groups-top-bar />
    </template>
  </v-app-bar>
</template>

<script>
import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import HealthcheckChipsList from '@/components/other/healthcheck/partials/healthcheck-chips-list.vue';

import AppLogo from './app-logo.vue';
import AppVersion from './app-version.vue';
import LoggedUsersCount from './logged-users-count.vue';
import GroupsTopBar from './groups-top-bar/groups-top-bar.vue';
import TopBarExploitationMenu from './top-bar-exploitation-menu.vue';
import TopBarAdministrationMenu from './top-bar-administration-menu.vue';
import TopBarNotificationsMenu from './top-bar-notifications-menu.vue';
import TopBarUserMenu from './top-bar-user-menu.vue';
import TopBarTitle from './top-bar-title.vue';

/**
 * Component for the top bar of the application
 *
 * @event toggleSideBar#click
 */
export default {
  components: {
    HealthcheckChipsList,
    AppLogo,
    AppVersion,
    LoggedUsersCount,
    GroupsTopBar,
    TopBarExploitationMenu,
    TopBarAdministrationMenu,
    TopBarNotificationsMenu,
    TopBarUserMenu,
    TopBarTitle,
  },
  mixins: [
    authMixin,
    entitiesInfoMixin,
  ],
  computed: {
    hasAccessToHealthcheckStatus() {
      return this.checkAccess(USERS_PERMISSIONS.technical.healthcheckStatus);
    },
  },
};
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

  & ::v-deep .logged-users-count {
    left: -8px;
  }
}

.brand {
  display: flex;
  align-items: center;
  margin: 0;
  width: 250px;
  height: 100%;

  img {
    margin: auto;
  }
}

.top-bar {
  & ::v-deep .v-toolbar__content {
    padding: 0;
  }

  & ::v-deep .v-toolbar__extension {
    padding: 0;
  }
}
</style>
