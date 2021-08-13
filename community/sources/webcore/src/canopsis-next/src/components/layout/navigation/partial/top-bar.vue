<template lang="pug">
  v-toolbar.top-bar.primary(
    :height="$config.TOP_BAR_HEIGHT",
    dense,
    fixed,
    app
  )
    v-toolbar-side-icon.ml-2.white--text(
      v-if="isShownGroupsSideBar",
      data-test="groupsSideBarButton",
      @click="$emit('toggleSideBar')"
    )
    v-layout(v-else, fill-height, align-center)
      app-logo.canopsis-logo.mr-2
      v-layout.version.ml-1(fill-height, align-end)
        active-sessions-count(badgeColor="secondary")
        app-version
    v-toolbar-title.white--text.font-weight-regular(v-if="appTitle") {{ appTitle }}
    healthcheck-chips-list
    portal-target(:name="$constants.PORTALS_NAMES.additionalTopBarItems")
    v-toolbar-items
      top-bar-exploitation-menu
      top-bar-administration-menu
      top-bar-notifications-menu
      top-bar-user-menu
    groups-top-bar(v-if="isShownGroupsTopBar", slot="extension")
</template>

<script>
import { authMixin } from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

import HealthcheckChipsList from '@/components/other/healthcheck/healthcheck-chips-list.vue';

import AppLogo from './app-logo.vue';
import AppVersion from './app-version.vue';
import ActiveSessionsCount from './active-sessions-count.vue';
import GroupsTopBar from './groups-top-bar/groups-top-bar.vue';
import TopBarExploitationMenu from './top-bar-exploitation-menu.vue';
import TopBarAdministrationMenu from './top-bar-administration-menu.vue';
import TopBarNotificationsMenu from './top-bar-notifications-menu.vue';
import TopBarUserMenu from './top-bar-user-menu.vue';

/**
 * Component for the top bar of the application
 *
 * @event toggleSideBar#click
 */
export default {
  components: {
    AppLogo,
    AppVersion,
    ActiveSessionsCount,
    GroupsTopBar,
    TopBarExploitationMenu,
    TopBarAdministrationMenu,
    TopBarNotificationsMenu,
    TopBarUserMenu,
    HealthcheckChipsList,
  },
  mixins: [
    authMixin,
    entitiesInfoMixin,
  ],
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

  & /deep/ .active-sessions-count {
    position: absolute;
    top: 0;
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
  & /deep/ .v-toolbar__content {
    padding: 0;
  }

  & /deep/ .v-toolbar__extension {
    padding: 0;
  }
}
</style>
