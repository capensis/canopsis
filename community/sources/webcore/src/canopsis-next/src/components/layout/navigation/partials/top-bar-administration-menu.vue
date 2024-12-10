<template>
  <v-menu
    v-if="administrationGroupedLinks.length"
    content-class="administration-menu__content"
    bottom
    offset-y
  >
    <template #activator="{ on }">
      <v-btn
        class="white--text"
        text
        v-on="on"
      >
        {{ $t('common.administration') }}
      </v-btn>
    </template>
    <v-list class="py-0">
      <template v-for="(group, index) in administrationGroupedLinks">
        <v-subheader
          :key="`${group.title}-title`"
          class="text-subtitle-1"
          @click.stop=""
        >
          {{ group.title }}
        </v-subheader>
        <top-bar-menu-link
          v-for="link in group.links"
          :key="link.title"
          :link="link"
          class="top-bar-administration-menu-link"
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
import { USER_PERMISSIONS, ROUTES_NAMES, GROUPED_USER_PERMISSIONS } from '@/constants';

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
          permission: USER_PERMISSIONS.technical.permission,
        },
        {
          route: { name: ROUTES_NAMES.adminRoles },
          icon: 'supervised_user_circle',
          permission: USER_PERMISSIONS.technical.role,
        },
        {
          route: { name: ROUTES_NAMES.adminUsers },
          icon: 'people',
          permission: USER_PERMISSIONS.technical.user,
        },
      ];
    },

    administrationCommunicationsLinks() {
      return [
        {
          route: { name: ROUTES_NAMES.adminBroadcastMessages },
          icon: '$vuetify.icons.bullhorn',
          permission: USER_PERMISSIONS.technical.broadcastMessage,
        },
        {
          route: { name: ROUTES_NAMES.adminPlaylists },
          icon: 'playlist_play',
          permission: USER_PERMISSIONS.technical.playlist,
        },
      ];
    },

    administrationGeneralLinks() {
      return [
        {
          route: { name: ROUTES_NAMES.adminEventsRecords },
          icon: '$vuetify.icons.mark_unread_chat_alt',
          permission: USER_PERMISSIONS.technical.eventsRecord,
        },
        {
          route: { name: ROUTES_NAMES.adminParameters },
          icon: 'settings',
          permission: USER_PERMISSIONS.technical.parameters,
        },
        {
          route: { name: ROUTES_NAMES.adminPlanning },
          icon: 'event_note',
          permission: GROUPED_USER_PERMISSIONS.planning,
        },
        {
          route: { name: ROUTES_NAMES.adminRemediation },
          icon: 'assignment',
          permission: GROUPED_USER_PERMISSIONS.remediation,
        },
        {
          route: { name: ROUTES_NAMES.adminKPI },
          icon: 'stacked_bar_chart',
          permission: USER_PERMISSIONS.technical.kpi,
        },
        {
          route: { name: ROUTES_NAMES.adminMaps },
          icon: 'edit_location',
          permission: USER_PERMISSIONS.technical.map,
        },
        {
          icon: '$vuetify.icons.build_circle',
          permission: USER_PERMISSIONS.technical.maintenance,
          handler: this.showToggleMaintenanceModeModal,
        },
        {
          route: { name: ROUTES_NAMES.adminTags },
          icon: 'local_offer',
          permission: USER_PERMISSIONS.technical.tag,
        },
        {
          route: { name: ROUTES_NAMES.adminHealthcheck },
          icon: '$vuetify.icons.alt_route',
          permission: USER_PERMISSIONS.technical.healthcheck,
        },
        {
          route: { name: ROUTES_NAMES.adminStorageSettings },
          icon: '$vuetify.icons.storage',
          permission: USER_PERMISSIONS.technical.storageSettings,
        },
        {
          route: { name: ROUTES_NAMES.adminStateSettings },
          icon: 'add_alert',
          permission: USER_PERMISSIONS.technical.stateSetting,
        },
      ];
    },

    permissionsWithDefaultType() {
      return [
        USER_PERMISSIONS.technical.engine,
        USER_PERMISSIONS.technical.healthcheck,
        USER_PERMISSIONS.technical.kpi,
        USER_PERMISSIONS.technical.maintenance,
        USER_PERMISSIONS.technical.eventsRecord,
      ];
    },
  },
};
</script>

<style lang="scss">
.administration-menu__content {
  max-height: 95vh;

  .v-avatar {
    border-radius: unset;
  }
}
</style>
