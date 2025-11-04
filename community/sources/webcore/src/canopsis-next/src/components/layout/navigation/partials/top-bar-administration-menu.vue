<template>
  <top-bar-menu
    :title="$t('common.administration')"
    :links="administrationLinks"
    :permissions-with-default-type="permissionsWithDefaultType"
    content-class="topbar-menu-administration__content"
    without-sort
  />
</template>

<script>
import { computed } from 'vue';

import { USER_PERMISSIONS, ROUTES_NAMES, GROUPED_USER_PERMISSIONS } from '@/constants';

import { uid } from '@/helpers/uid';

import { useI18n } from '@/hooks/i18n';

import { useMaintenanceActions } from './hooks/maintenance-actions';
import TopBarMenu from './top-bar-menu.vue';

export default {
  components: { TopBarMenu },
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
      USER_PERMISSIONS.technical.viewImportExport,
      USER_PERMISSIONS.technical.notification.common,
    ];

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

    const customObjectsLinks = computed(() => [
      {
        route: { name: ROUTES_NAMES.adminCustomObjectsExternalAuthTokens },
        icon: 'security',
        permission: USER_PERMISSIONS.technical.externalAuthTokens,
      },
      {
        route: { name: ROUTES_NAMES.adminCustomObjectsExternalDataTables },
        icon: '$vuetify.icons.database_outlined',
        permission: USER_PERMISSIONS.technical.exploitation.externalDataTable,
      },
      {
        route: { name: ROUTES_NAMES.adminCustomObjectsEntityInfosProperties },
        icon: 'info',
        permission: USER_PERMISSIONS.technical.exploitation.entityInfoProperty,
      },
      {
        route: { name: ROUTES_NAMES.adminCustomObjectsIcons },
        icon: '$vuetify.icons.square_circle',
        title: tc('common.icon', 2),
        permission: USER_PERMISSIONS.technical.icon,
      },
      {
        route: { name: ROUTES_NAMES.adminCustomObjectsMaps },
        icon: 'edit_location',
        permission: USER_PERMISSIONS.technical.map,
      },
      {
        route: { name: ROUTES_NAMES.adminCustomObjectsTags },
        icon: 'local_offer',
        permission: USER_PERMISSIONS.technical.tag,
      },
    ]);

    const settingsLinks = computed(() => [
      {
        route: { name: ROUTES_NAMES.adminSettingsUserInterface },
        icon: 'computer',
        permission: USER_PERMISSIONS.technical.parameters,
      },
      {
        route: { name: ROUTES_NAMES.adminSettingsViewsImportExport },
        icon: 'import_export',
        permission: USER_PERMISSIONS.technical.viewImportExport,
      },
      {
        route: { name: ROUTES_NAMES.adminSettingsNotifications },
        icon: 'notifications',
        permission: USER_PERMISSIONS.technical.notification.common,
      },
      {
        route: { name: ROUTES_NAMES.adminSettingsCommentTemplates },
        icon: 'comment',
        permission: USER_PERMISSIONS.technical.commentTemplate,
      },
      {
        route: { name: ROUTES_NAMES.adminSettingsWidgetTemplates },
        icon: 'widgets',
        permission: USER_PERMISSIONS.technical.widgetTemplate,
      },
      {
        route: { name: ROUTES_NAMES.adminSettingsStateSettings },
        icon: 'add_alert',
        permission: USER_PERMISSIONS.technical.stateSetting,
      },
      {
        route: { name: ROUTES_NAMES.adminSettingsStorageSettings },
        icon: '$vuetify.icons.storage',
        permission: USER_PERMISSIONS.technical.storageSettings,
      },
    ]);

    const generalLinksWithChildren = [
      {
        icon: 'local_offer',
        title: t('layout.topbar.customObjects'),
        links: customObjectsLinks.value,
      },
      {
        icon: 'settings',
        title: t('common.settings'),
        links: settingsLinks.value,
      },
    ];

    const administrationLinks = computed(() => [
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
    ].reduce((acc, group) => {
      if (group.title) {
        acc.push({ title: group.title, header: true });
      }

      acc.push(
        ...group.links,

        { divider: true, key: uid() },
      );

      return acc;
    }, []));

    return {
      permissionsWithDefaultType,

      administrationLinks,
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
