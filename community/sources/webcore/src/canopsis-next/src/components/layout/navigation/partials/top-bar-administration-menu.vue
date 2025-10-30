<template>
  <v-menu
    v-if="administrationGroupedLinks.length"
    content-class="top-bar-menu__content topbar-menu-administration__content"
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
          v-if="group.title"
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
          v-if="index < administrationGroupedLinks.length - 1"
          :key="`${group.title}-divider`"
        />
      </template>
    </v-list>
  </v-menu>
</template>

<script>
import { computed } from 'vue';

import { USER_PERMISSIONS, ROUTES_NAMES, GROUPED_USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import { useTopBarMenu } from './hooks/top-bar-menu';
import { useMaintenanceActions } from './hooks/maintenance-actions';
import TopBarMenuLink from './top-bar-menu-link.vue';

export default {
  name: 'TopBarAdministrationMenu', // We need it for recursive
  components: { TopBarMenuLink },
  setup() {
    const { t, tc } = useI18n();
    const { showToggleMaintenanceModeModal } = useMaintenanceActions();

    const permissionsWithDefaultType = [
      USER_PERMISSIONS.technical.engine,
      USER_PERMISSIONS.technical.healthcheck,
      USER_PERMISSIONS.technical.kpi,
      USER_PERMISSIONS.technical.maintenance,
      USER_PERMISSIONS.technical.eventsRecord,
      USER_PERMISSIONS.technical.templateTesting,
    ];

    const { prepareLinks } = useTopBarMenu({ permissionsWithDefaultType });

    const accessLinks = [
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

    const maintenanceLinks = [
      {
        icon: '$vuetify.icons.build_circle',
        permission: USER_PERMISSIONS.technical.maintenance,
        handler: showToggleMaintenanceModeModal,
      },
      {
        route: { name: ROUTES_NAMES.adminPlanning },
        icon: 'event_note',
        permission: GROUPED_USER_PERMISSIONS.planning,
      },
    ];

    const communicationsLinks = [
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

    const generalLinks = [

      {
        route: { name: ROUTES_NAMES.adminHealthcheck },
        icon: '$vuetify.icons.alt_route',
        permission: USER_PERMISSIONS.technical.healthcheck,
      },
      {
        route: { name: ROUTES_NAMES.adminKPI },
        icon: 'stacked_bar_chart',
        permission: USER_PERMISSIONS.technical.kpi,
      },
      {
        route: { name: ROUTES_NAMES.adminEventsRecords },
        icon: '$vuetify.icons.mark_unread_chat_alt',
        permission: USER_PERMISSIONS.technical.eventsRecord,
      },
      {
        route: { name: ROUTES_NAMES.adminTemplateTesting },
        icon: '$vuetify.icons.play_circle',
        permission: USER_PERMISSIONS.technical.templateTesting,
      },
    ];

    const generalLinksWithChildren = [
      {
        icon: 'ticket',
        title: t('layout.topbar.customObjects'),
        links: [
          {
            route: { name: ROUTES_NAMES.adminExternalAuthTokens },
            icon: 'security',
            permission: USER_PERMISSIONS.technical.externalAuthTokens,
          },
          {
            route: { name: ROUTES_NAMES.adminExternalDataTables },
            icon: '$vuetify.icons.database_outlined',
            permission: USER_PERMISSIONS.technical.exploitation.externalDataTable,
          },
          {
            route: { name: ROUTES_NAMES.adminEntityInfosProperties },
            icon: 'info',
            permission: USER_PERMISSIONS.technical.exploitation.entityInfoProperty,
          },
        ],
      },
      {
        route: { name: ROUTES_NAMES.adminParameters },
        icon: 'settings',
        permission: USER_PERMISSIONS.technical.parameters,
      },
      {
        route: { name: ROUTES_NAMES.adminMaps },
        icon: 'edit_location',
        permission: USER_PERMISSIONS.technical.map,
      },
      {
        route: { name: ROUTES_NAMES.adminTags },
        icon: 'local_offer',
        permission: USER_PERMISSIONS.technical.tag,
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

    const administrationGroupedLinks = computed(() => {
      const groupedLinks = [
        {
          title: t('common.access'),
          links: accessLinks,
        },
        {
          title: t('common.maintenance'),
          links: maintenanceLinks,
        },
        {
          title: tc('common.communication', 2),
          links: communicationsLinks,
        },
        {
          links: generalLinks,
        },
        {
          links: generalLinksWithChildren,
        },
      ];

      return groupedLinks.reduce((acc, group) => {
        const links = prepareLinks(group.links);

        if (links.length) {
          acc.push({ links, title: group.title });
        }

        return acc;
      }, []);
    });

    return {
      administrationGroupedLinks,
    };
  },
};
</script>

<style lang="scss">
.topbar-menu-administration__content {
  .v-subheader {
    text-transform: uppercase;
    font-weight: 500 !important;
    font-size: 1rem !important;
    line-height: 1rem !important;
    height: 40px;
  }

  .v-list-item {
    padding-left: 24px;
  }
}
</style>
