<template lang="pug">
  v-menu(v-show="administrationGroupedLinks.length", bottom, offset-y)
    template(#activator="{ on }")
      v-btn.white--text(v-on="on", flat) {{ $t('common.administration') }}
    v-list.py-0
      template(v-for="(group, index) in administrationGroupedLinks")
        v-subheader.subheading(:key="`${group.title}-title`", @click.stop="") {{ group.title }}
        top-bar-menu-link.top-bar-administration-menu-link(
          v-for="link in group.links",
          :key="link.title",
          :link="link"
        )
        v-divider(
          v-if="index < administrationGroupedLinks.length - 1",
          :key="`${group.title}-divider`"
        )
</template>

<script>
import { USERS_PERMISSIONS, ROUTES_NAMES } from '@/constants';

import { layoutNavigationTopBarMenuMixin } from '@/mixins/layout/navigation/top-bar-menu';

import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
  mixins: [layoutNavigationTopBarMenuMixin],
  computed: {
    administrationGroupedLinks() {
      const groupedLinks = [
        {
          title: this.$t('common.access'),
          links: this.administrationAccessLinks,
        },
        {
          title: this.$tc('common.communication', 2),
          links: this.administrationCommunicationsLinks,
        },
        {
          title: this.$t('common.general'),
          links: this.administrationGeneralLinks,
        },
      ];

      return groupedLinks.reduce((acc, group) => {
        const links = this.prepareLinks(group.links);

        if (links.length) {
          acc.push({ links, title: group.title });
        }

        return acc;
      }, []);
    },

    administrationAccessLinks() {
      return [
        {
          route: { name: ROUTES_NAMES.adminRights },
          icon: 'verified_user',
          permission: USERS_PERMISSIONS.technical.permission,
        },
        {
          route: { name: ROUTES_NAMES.adminRoles },
          icon: 'supervised_user_circle',
          permission: USERS_PERMISSIONS.technical.role,
        },
        {
          route: { name: ROUTES_NAMES.adminUsers },
          icon: 'people',
          permission: USERS_PERMISSIONS.technical.user,
        },
      ];
    },

    administrationCommunicationsLinks() {
      return [
        {
          route: { name: ROUTES_NAMES.adminBroadcastMessages },
          icon: '$vuetify.icons.bullhorn',
          permission: USERS_PERMISSIONS.technical.broadcastMessage,
        },
        {
          route: { name: ROUTES_NAMES.adminPlaylists },
          icon: 'playlist_play',
          permission: USERS_PERMISSIONS.technical.playlist,
        },
      ];
    },

    administrationGeneralLinks() {
      const links = [
        {
          route: { name: ROUTES_NAMES.adminParameters },
          icon: 'settings',
          permission: USERS_PERMISSIONS.technical.parameters,
        },
        {
          route: { name: ROUTES_NAMES.adminPlanning },
          icon: 'event_note',
          permission: USERS_PERMISSIONS.technical.planning,
        },
        {
          route: { name: ROUTES_NAMES.adminRemediation },
          icon: 'assignment',
          permission: USERS_PERMISSIONS.technical.remediation,
        },
        {
          route: { name: ROUTES_NAMES.adminKPI },
          icon: 'stacked_bar_chart',
          permission: USERS_PERMISSIONS.technical.kpi,
        },
        {
          route: { name: ROUTES_NAMES.adminMaps },
          icon: 'edit_location',
          permission: USERS_PERMISSIONS.technical.map,
        },
      ];

      const enginesLink = this.isProVersion
        ? {
          route: { name: ROUTES_NAMES.adminHealthcheck },
          icon: '$vuetify.icons.alt_route',
          permission: USERS_PERMISSIONS.technical.healthcheck,
        }
        : {
          route: { name: ROUTES_NAMES.adminEngines },
          icon: '$vuetify.icons.alt_route',
          permission: USERS_PERMISSIONS.technical.engine,
        };

      links.push(enginesLink);

      return links;
    },

    permissionsWithDefaultType() {
      return [
        USERS_PERMISSIONS.technical.engine,
        USERS_PERMISSIONS.technical.healthcheck,
        USERS_PERMISSIONS.technical.kpi,
      ];
    },
  },
};
</script>

<style lang="scss" scoped>
.top-bar-administration-menu-link ::v-deep span {
  margin-left: 8px;
}
</style>
