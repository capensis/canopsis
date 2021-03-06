<template lang="pug">
  v-toolbar.top-bar.primary(
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
      app-logo.canopsisLogo.mr-2
      v-layout.version.ml-1(fill-height, align-end)
        active-sessions-count(badgeColor="secondary")
        app-version
    v-toolbar-title.white--text.font-weight-regular(v-if="appTitle") {{ appTitle }}
    v-spacer
    portal-target(:name="$constants.PORTALS_NAMES.additionalTopBarItems")
    v-toolbar-items
      v-menu(v-show="exploitationLinks.length", bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.exploitation') }}
        v-list.pb-0
          v-list-tile(v-for="(link, index) in exploitationLinks", :key="`exploitation-${index}`")
            v-list-tile-title
              router-link(:to="link.route")
                v-layout(justify-space-between)
                  span.black--text {{ link.text }}
                  v-icon.ml-2 {{ link.icon }}
      v-menu(v-show="administrationLinks.length", bottom, offset-y)
        v-btn.white--text(slot="activator", flat) {{ $t('common.administration') }}
        v-list.pb-0
          v-list-tile(v-for="(link, index) in administrationLinks", :key="`administration-${index}`")
            v-list-tile-title
              router-link(:to="link.route")
                v-layout(justify-space-between)
                  span.black--text {{ link.text }}
                  v-icon.ml-2 {{ link.icon }}
      v-menu(bottom, offset-y, offset-x)
        v-btn.white--text(
          slot="activator",
          data-test="userTopBarDropdownButton",
          flat
        ) {{ userName }}
        v-list.pb-0
          v-list-tile
            v-list-tile-content
              v-btn.ma-0.pa-1(data-test="userProfileButton", flat, @click.prevent="showEditUserModal")
                v-layout(align-center)
                  v-icon person
                  div.ml-2 {{ $t('user.seeProfile') }}
          v-list-tile
            v-list-tile-content
              v-btn.ma-0.pa-1.error--text(data-test="logoutButton", flat, @click.prevent="logoutHandler")
                v-layout(align-center)
                  v-icon exit_to_app
                  div.ml-2 {{ $t('common.logout') }}
    template(v-if="isShownGroupsTopBar", slot="extension")
      groups-top-bar
</template>

<script>
import { MODALS, USERS_RIGHTS } from '@/constants';

import { prepareUserByData } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import entitiesUserMixin from '@/mixins/entities/user';
import entitiesInfoMixin from '@/mixins/entities/info';

import AppLogo from './app-logo.vue';
import AppVersion from './app-version.vue';
import ActiveSessionsCount from './active-sessions-count.vue';
import GroupsTopBar from './groups-top-bar/groups-top-bar.vue';

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
  },
  mixins: [
    authMixin,
    entitiesUserMixin,
    entitiesInfoMixin,
  ],
  computed: {
    exploitationLinks() {
      const links = [
        {
          route: { name: 'exploitation-event-filter' },
          text: this.$t('eventFilter.title'),
          icon: 'list',
          right: USERS_RIGHTS.technical.exploitation.eventFilter,
        },
        {
          route: { name: 'exploitation-pbehaviors' },
          text: this.$t('common.pbehaviors'),
          icon: 'pause',
          right: USERS_RIGHTS.technical.exploitation.pbehavior,
        },
        {
          route: { name: 'exploitation-webhooks' },
          text: this.$t('common.webhooks'),
          icon: '$vuetify.icons.webhook',
          right: USERS_RIGHTS.technical.exploitation.webhook,
        },
        {
          route: { name: 'exploitation-snmp-rules' },
          text: this.$t('snmpRules.title'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.exploitation.snmpRule,
        },
        {
          route: { name: 'exploitation-actions' },
          text: this.$t('actions.title'),
          icon: 'linear_scale',
          right: USERS_RIGHTS.technical.exploitation.action,
        },
        {
          route: { name: 'exploitation-heartbeats' },
          text: this.$t('heartbeat.title'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.exploitation.heartbeat,
        },
        {
          route: { name: 'exploitation-dynamic-infos' },
          text: this.$t('dynamicInfo.title'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.exploitation.dynamicInfo,
        },
        {
          route: { name: 'exploitation-meta-alarm-rules' },
          text: this.$t('metaAlarmRule.title'),
          icon: 'list',
          right: USERS_RIGHTS.technical.exploitation.metaAlarmRule,
        },
      ];

      return links.filter(({ right }) =>
        this.checkAppInfoAccessByRight(right)
        && this.checkReadAccess(right));
    },

    administrationLinks() {
      const links = [
        {
          route: { name: 'admin-rights' },
          text: this.$t('common.rights'),
          icon: 'verified_user',
          right: USERS_RIGHTS.technical.action,
        },
        {
          route: { name: 'admin-users' },
          text: this.$t('common.users'),
          icon: 'people',
          right: USERS_RIGHTS.technical.user,
        },
        {
          route: { name: 'admin-roles' },
          text: this.$t('common.roles'),
          icon: 'supervised_user_circle',
          right: USERS_RIGHTS.technical.role,
        },
        {
          route: { name: 'admin-parameters' },
          text: this.$t('common.parameters'),
          icon: 'settings',
          right: USERS_RIGHTS.technical.parameters,
        },
        {
          route: { name: 'admin-broadcast-messages' },
          text: this.$t('common.broadcastMessages'),
          icon: '$vuetify.icons.bullhorn',
          right: USERS_RIGHTS.technical.broadcastMessage,
        },
        {
          route: { name: 'admin-playlists' },
          text: this.$t('common.playlists'),
          icon: 'playlist_play',
          right: USERS_RIGHTS.technical.playlist,
        },
        {
          route: { name: 'admin-planning-administration' },
          text: this.$t('common.planningAdministration'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.planning,
        },
        {
          route: { name: 'remediation-administration' },
          text: this.$t('common.remediation'),
          icon: 'assignment',
          right: USERS_RIGHTS.technical.remediation,
        },
      ];

      const linksWithDefaultRight = [
        {
          route: { name: 'admin-engines' },
          text: this.$t('common.engines'),
          icon: '$vuetify.icons.alt_route',
          right: USERS_RIGHTS.technical.engine,
        },
      ];

      return links
        .filter(({ right }) => this.checkAppInfoAccessByRight(right) && this.checkReadAccess(right))
        .concat(linksWithDefaultRight.filter(({ right }) => this.checkAccess(right)));
    },

    userName() {
      return this.currentUser.crecord_name || this.currentUser._id;
    },
  },
  methods: {
    showEditUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          title: this.$t('common.profile'),
          user: this.currentUser,
          onlyUserPrefs: true,
          action: async (data) => {
            await this.createUserWithPopup({ data: prepareUserByData(data, this.currentUser) });

            await this.fetchCurrentUser();
          },
        },
      });
    },
    logoutHandler() {
      return this.logout({ redirectTo: { name: 'login' } });
    },
  },
};
</script>

<style lang="scss" scoped>
  .canopsisLogo {
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

  a {
    text-decoration: none;
    color: inherit;
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
