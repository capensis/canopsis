<template>
  <v-menu
    v-if="administrationGroupedLinks.length"
    bottom
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        v-on="on"
        text
      >
        {{ $t('common.administration') }}
      </v-btn>
    </template>
    <v-list class="py-0">
      <template v-for="(group, index) in administrationGroupedLinks">
        <v-subheader
          class="text-subtitle-1"
          :key="`${group.title}-title`"
          @click.stop=""
        >
          {{ group.title }}
        </v-subheader>
        <top-bar-menu-link
          class="top-bar-administration-menu-link"
          v-for="link in group.links"
          :key="link.title"
          :link="link"
        />
        <v-divider
          v-if="index &lt; administrationGroupedLinks.length - 1"
          :key="`${group.title}-divider`"
        />
      </template>
    </v-list>
  </v-menu>
</template>

<script>
import { USERS_PERMISSIONS, ROUTES_NAMES } from '@/constants';

import { layoutNavigationTopBarMenuMixin } from '@/mixins/layout/navigation/top-bar-menu';
import { maintenanceActionsMixin } from '@/mixins/maintenance/maintenance-actions';

import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  components: { TopBarMenuLink },
  mixins: [layoutNavigationTopBarMenuMixin, maintenanceActionsMixin],
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
      return [
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
        {
          icon: '$vuetify.icons.build_circle',
          permission: USERS_PERMISSIONS.technical.maintenance,
          handler: this.showToggleMaintenanceModeModal,
        },
        {
          route: { name: ROUTES_NAMES.adminTags },
          icon: 'local_offer',
          permission: USERS_PERMISSIONS.technical.tag,
        },
        {
          route: { name: ROUTES_NAMES.adminHealthcheck },
          icon: '$vuetify.icons.alt_route',
          permission: USERS_PERMISSIONS.technical.healthcheck,
        },
      ];
    },

    permissionsWithDefaultType() {
      return [
        USERS_PERMISSIONS.technical.engine,
        USERS_PERMISSIONS.technical.healthcheck,
        USERS_PERMISSIONS.technical.kpi,
        USERS_PERMISSIONS.technical.maintenance,
      ];
    },
  },
};
</script>
