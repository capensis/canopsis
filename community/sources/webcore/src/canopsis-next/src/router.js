import Vue from 'vue';
import Router from 'vue-router';

import { ROUTER_MODE, ROUTER_ACCESS_TOKEN_KEY } from '@/config';
import {
  CRUD_ACTIONS,
  ROUTES_NAMES,
  ROUTES,
  OLD_ROUTES,
  USER_PERMISSIONS,
  GROUPED_USER_PERMISSIONS,
} from '@/constants';

import store from '@/store';

import { checkAppInfoAccessForRoute, checkUserAccessForRoute } from '@/helpers/router';

import Login from '@/views/login.vue';
import Error from '@/views/error.vue';

const Home = () => import(/* webpackChunkName: "Home" */ '@/views/home.vue');
const View = () => import(/* webpackChunkName: "View" */ '@/views/view.vue');
const ViewKiosk = () => import(/* webpackChunkName: "View" */ '@/views/view-kiosk.vue');
const Alarm = () => import(/* webpackChunkName: "Alarm" */ '@/views/alarm.vue');
const AdminPermissions = () => import(/* webpackChunkName: "Permission" */ '@/views/admin/permissions.vue');
const AdminUsers = () => import(/* webpackChunkName: "User" */ '@/views/admin/users.vue');
const AdminRoles = () => import(/* webpackChunkName: "Role" */ '@/views/admin/roles.vue');
const AdminBroadcastMessages = () => import(/* webpackChunkName: "BroadcastMessage" */ '@/views/admin/broadcast-messages.vue');
const AdminPlaylists = () => import(/* webpackChunkName: "Playlist" */ '@/views/admin/playlists.vue');
const AdminPlanning = () => import(/* webpackChunkName: "Planning" */ '@/views/admin/planning.vue');
const AdminHealthcheck = () => import(/* webpackChunkName: "Healthcheck" */ '@/views/admin/healthcheck.vue');
const AdminKPI = () => import(/* webpackChunkName: "KPI" */ '@/views/admin/kpi.vue');
const AdminEventsRecords = () => import(/* webpackChunkName: "EventsRecords" */ '@/views/admin/events-records.vue');
const AdminTemplateTesting = () => import(/* webpackChunkName: "TemplateTesting" */ '@/views/admin/template-testing.vue');
const AdminCustomObjectsExternalAuthTokens = () => import(/* webpackChunkName: "ExternalAuthTokens" */ '@/views/admin/custom-objects/external-auth-tokens.vue');
const AdminCustomObjectsEntityInfosProperties = () => import(/* webpackChunkName: "EntityInfosProperties" */ '@/views/admin/custom-objects/entity-infos-properties.vue');
const AdminCustomObjectsExternalDataTables = () => import(/* webpackChunkName: "ExternalDataTables" */ '@/views/admin/custom-objects/external-data-tables.vue');
const AdminCustomObjectsIcons = () => import(/* webpackChunkName: "Icons" */ '@/views/admin/custom-objects/icons.vue');
const AdminCustomObjectsMaps = () => import(/* webpackChunkName: "Maps" */ '@/views/admin/custom-objects/maps.vue');
const AdminCustomObjectsTags = () => import(/* webpackChunkName: "Tags" */ '@/views/admin/custom-objects/tags.vue');
const AdminCustomObjectsLlms = () => import(/* webpackChunkName: "Llms" */ '@/views/admin/custom-objects/llms.vue');
const AdminSettingsUserInterface = () => import(/* webpackChunkName: "UserInterface" */ '@/views/admin/settings/user-interface.vue');
const AdminSettingsViewsImportExport = () => import(/* webpackChunkName: "ViewsImportExport" */ '@/views/admin/settings/views-import-export.vue');
const AdminSettingsNotifications = () => import(/* webpackChunkName: "NotificationsSettings" */ '@/views/admin/settings/notifications-settings.vue');
const AdminSettingsCommentTemplates = () => import(/* webpackChunkName: "CommentTemplates" */ '@/views/admin/settings/comment-templates.vue');
const AdminSettingsWidgetTemplates = () => import(/* webpackChunkName: "WidgetTemplates" */ '@/views/admin/settings/widget-templates.vue');
const AdminSettingsStorageSettings = () => import(/* webpackChunkName: "Tags" */ '@/views/admin/settings/storage-settings.vue');
const AdminSettingsStateSettings = () => import(/* webpackChunkName: "Tags" */ '@/views/admin/settings/state-settings.vue');
const ExploitationPbehaviors = () => import(/* webpackChunkName: "Pbehavior" */ '@/views/exploitation/pbehaviors.vue');
const ExploitationEventFilters = () => import(/* webpackChunkName: "EventFilters" */ '@/views/exploitation/event-filters.vue');
const ExploitationSnmpRules = () => import(/* webpackChunkName: "SnmpRule" */ '@/views/exploitation/snmp-rules.vue');
const ExploitationDynamicInfos = () => import(/* webpackChunkName: "DynamicInfo" */ '@/views/exploitation/dynamic-infos.vue');
const ExploitationMetaAlarmRules = () => import(/* webpackChunkName: "MetaAlarmRule" */ '@/views/exploitation/meta-alarm-rules.vue');
const ExploitationScenarios = () => import(/* webpackChunkName: "Scenario" */ '@/views/exploitation/scenarios.vue');
const ExploitationIdleRules = () => import(/* webpackChunkName: "IdleRule" */ '@/views/exploitation/idle-rules.vue');
const ExploitationFlappingRules = () => import(/* webpackChunkName: "AlarmStatusRule" */ '@/views/exploitation/flapping-rules.vue');
const ExploitationResolveRules = () => import(/* webpackChunkName: "AlarmStatusRule" */ '@/views/exploitation/resolve-rules.vue');
const ExploitationDeclareTicketRules = () => import(/* webpackChunkName: "DeclareTicketRule" */ '@/views/exploitation/declare-ticket-rules.vue');
const ExploitationLinkRules = () => import(/* webpackChunkName: "LinkRule" */ '@/views/exploitation/link-rules.vue');
const ExploitationRemediation = () => import(/* webpackChunkName: "Remediation" */ '@/views/exploitation/remediation.vue');
const ProfilePatterns = () => import(/* webpackChunkName: "Pattern" */ '@/views/profile/patterns.vue');
const ProfileThemes = () => import(/* webpackChunkName: "Theme" */ '@/views/profile/themes.vue');
const Playlist = () => import(/* webpackChunkName: "Playlist" */ '@/views/playlist.vue');
const Notifications = () => import(/* webpackChunkName: "Notifications" */ '@/views/notification/notifications.vue');

Vue.use(Router);

const routes = [
  {
    path: ROUTES.login,
    name: ROUTES_NAMES.login,
    component: Login,
    meta: {
      hideHeader: true,
      requiresLogin: false,
    },
  },
  {
    path: ROUTES.home,
    name: ROUTES_NAMES.home,
    component: Home,
    meta: {
      requiresLogin: true,
    },
  },
  {
    path: ROUTES.view,
    name: ROUTES_NAMES.view,
    component: View,
    meta: {
      requiresLogin: true,
    },
    props: route => ({ id: route.params.id }),
  },
  {
    path: ROUTES.viewKiosk,
    name: ROUTES_NAMES.viewKiosk,
    component: ViewKiosk,
    meta: {
      simpleNavigation: true,
      requiresLogin: true,
      requiresPermission: {
        id: route => route.params.id,
      },
    },
    props: route => ({ id: route.params.id, tabId: route.params.tabId }),
  },
  {
    path: ROUTES.alarms,
    name: ROUTES_NAMES.alarms,
    component: Alarm,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.view,
      },
    },
    props: route => ({ id: route.params.id }),
  },
  {
    path: ROUTES.adminRights,
    name: ROUTES_NAMES.adminRights,
    component: AdminPermissions,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.permission,
      },
    },
  },
  {
    path: ROUTES.adminUsers,
    name: ROUTES_NAMES.adminUsers,
    component: AdminUsers,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.user,
      },
    },
  },
  {
    path: ROUTES.adminRoles,
    name: ROUTES_NAMES.adminRoles,
    component: AdminRoles,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.role,
      },
    },
  },
  {
    path: ROUTES.adminBroadcastMessages,
    name: ROUTES_NAMES.adminBroadcastMessages,
    component: AdminBroadcastMessages,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.broadcastMessage,
      },
    },
  },
  {
    path: ROUTES.adminPlaylists,
    name: ROUTES_NAMES.adminPlaylists,
    component: AdminPlaylists,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.playlist,
      },
    },
  },
  {
    path: ROUTES.adminPlanning,
    name: ROUTES_NAMES.adminPlanning,
    component: AdminPlanning,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: GROUPED_USER_PERMISSIONS.planning,
      },
    },
  },
  {
    path: ROUTES.adminHealthcheck,
    name: ROUTES_NAMES.adminHealthcheck,
    component: AdminHealthcheck,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.healthcheck,
      },
    },
  },
  {
    path: ROUTES.adminKPI,
    name: ROUTES_NAMES.adminKPI,
    component: AdminKPI,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.kpi,
      },
    },
  },
  {
    path: ROUTES.adminEventsRecords,
    name: ROUTES_NAMES.adminEventsRecords,
    component: AdminEventsRecords,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.eventsRecord,
      },
    },
  },
  {
    path: ROUTES.adminTemplateTesting,
    name: ROUTES_NAMES.adminTemplateTesting,
    component: AdminTemplateTesting,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.templateTesting,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsExternalAuthTokens,
    name: ROUTES_NAMES.adminCustomObjectsExternalAuthTokens,
    component: AdminCustomObjectsExternalAuthTokens,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.externalAuthTokens,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsExternalDataTables,
    name: ROUTES_NAMES.adminCustomObjectsExternalDataTables,
    component: AdminCustomObjectsExternalDataTables,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.externalDataTable,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsEntityInfosProperties,
    name: ROUTES_NAMES.adminCustomObjectsEntityInfosProperties,
    component: AdminCustomObjectsEntityInfosProperties,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.entityInfoProperty,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsIcons,
    name: ROUTES_NAMES.adminCustomObjectsIcons,
    component: AdminCustomObjectsIcons,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.icon,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsMaps,
    name: ROUTES_NAMES.adminCustomObjectsMaps,
    component: AdminCustomObjectsMaps,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.map,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsTags,
    name: ROUTES_NAMES.adminCustomObjectsTags,
    component: AdminCustomObjectsTags,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.tag,
      },
    },
  },
  {
    path: ROUTES.adminCustomObjectsLlms,
    name: ROUTES_NAMES.adminCustomObjectsLlms,
    component: AdminCustomObjectsLlms,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.llm,
      },
    },
  },
  {
    path: ROUTES.adminSettingsUserInterface,
    name: ROUTES_NAMES.adminSettingsUserInterface,
    component: AdminSettingsUserInterface,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.parameters,
      },
    },
  },
  {
    path: ROUTES.adminSettingsViewsImportExport,
    name: ROUTES_NAMES.adminSettingsViewsImportExport,
    component: AdminSettingsViewsImportExport,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.viewImportExport,
      },
    },
  },
  {
    path: ROUTES.adminSettingsNotifications,
    name: ROUTES_NAMES.adminSettingsNotifications,
    component: AdminSettingsNotifications,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.notification.common,
      },
    },
  },
  {
    path: ROUTES.adminSettingsCommentTemplates,
    name: ROUTES_NAMES.adminSettingsCommentTemplates,
    component: AdminSettingsCommentTemplates,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.commentTemplate,
      },
    },
  },
  {
    path: ROUTES.adminSettingsWidgetTemplates,
    name: ROUTES_NAMES.adminSettingsWidgetTemplates,
    component: AdminSettingsWidgetTemplates,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.widgetTemplate,
      },
    },
  },
  {
    path: ROUTES.adminSettingsStorageSettings,
    name: ROUTES_NAMES.adminSettingsStorageSettings,
    component: AdminSettingsStorageSettings,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.storageSettings,
      },
    },
  },
  {
    path: ROUTES.adminSettingsStateSettings,
    name: ROUTES_NAMES.adminSettingsStateSettings,
    component: AdminSettingsStateSettings,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.stateSetting,
      },
    },
  },
  {
    path: ROUTES.exploitationPbehaviors,
    name: ROUTES_NAMES.exploitationPbehaviors,
    component: ExploitationPbehaviors,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.pbehavior,
      },
    },
  },
  {
    path: ROUTES.exploitationEventFilters,
    name: ROUTES_NAMES.exploitationEventFilters,
    component: ExploitationEventFilters,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.eventFilter,
      },
    },
  },
  {
    path: ROUTES.exploitationSnmpRules,
    name: ROUTES_NAMES.exploitationSnmpRules,
    component: ExploitationSnmpRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.snmpRule,
      },
    },
  },
  {
    path: ROUTES.exploitationDynamicInfos,
    name: ROUTES_NAMES.exploitationDynamicInfos,
    component: ExploitationDynamicInfos,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.dynamicInfo,
      },
    },
  },
  {
    path: ROUTES.exploitationMetaAlarmRules,
    name: ROUTES_NAMES.exploitationMetaAlarmRules,
    component: ExploitationMetaAlarmRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.metaAlarmRule,
      },
    },
  },
  {
    path: ROUTES.exploitationScenarios,
    name: ROUTES_NAMES.exploitationScenarios,
    component: ExploitationScenarios,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.scenario,
      },
    },
  },
  {
    path: ROUTES.exploitationIdleRules,
    name: ROUTES_NAMES.exploitationIdleRules,
    component: ExploitationIdleRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.idleRules,
      },
    },
  },
  {
    path: ROUTES.exploitationFlappingRules,
    name: ROUTES_NAMES.exploitationFlappingRules,
    component: ExploitationFlappingRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.flappingRules,
      },
    },
  },
  {
    path: ROUTES.exploitationResolveRules,
    name: ROUTES_NAMES.exploitationResolveRules,
    component: ExploitationResolveRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.resolveRules,
      },
    },
  },
  {
    path: ROUTES.exploitationDeclareTicketRules,
    name: ROUTES_NAMES.exploitationDeclareTicketRules,
    component: ExploitationDeclareTicketRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.declareTicketRule,
      },
    },
  },
  {
    path: ROUTES.exploitationLinkRules,
    name: ROUTES_NAMES.exploitationLinkRules,
    component: ExploitationLinkRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.linkRule,
      },
    },
  },
  {
    path: ROUTES.exploitationRemediation,
    name: ROUTES_NAMES.exploitationRemediation,
    component: ExploitationRemediation,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: GROUPED_USER_PERMISSIONS.remediation,
      },
    },
  },
  {
    path: ROUTES.playlist,
    name: ROUTES_NAMES.playlist,
    component: Playlist,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: route => route.params.id,
        action: CRUD_ACTIONS.read,
      },
    },
    props: route => ({ id: route.params.id, autoplay: String(route.query.autoplay) === 'true' }),
  },
  {
    path: ROUTES.profilePatterns,
    name: ROUTES_NAMES.profilePatterns,
    component: ProfilePatterns,
    meta: {
      requiresLogin: true,
    },
  },
  {
    path: ROUTES.profileThemes,
    name: ROUTES_NAMES.profileThemes,
    component: ProfileThemes,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USER_PERMISSIONS.technical.profile.theme,
      },
    },
  },
  {
    path: ROUTES.notifications,
    name: ROUTES_NAMES.notifications,
    component: Notifications,
    meta: {
      requiresLogin: true,
    },
    props: route => ({ tabId: route.params.tabId, activeId: route.query.id }),
  },
  {
    path: ROUTES.error,
    name: ROUTES_NAMES.error,
    component: Error,
    meta: {
      hideHeader: true,
    },
    props: route => ({ message: route.query.message, redirect: route.query.redirect }),
  },

  /**
   * REDIRECTS FOR OLD ROUTES
   */
  {
    path: OLD_ROUTES.remediation,
    redirect: { name: ROUTES_NAMES.exploitationRemediation },
  },
  {
    path: OLD_ROUTES.externalDataTables,
    redirect: { name: ROUTES_NAMES.adminCustomObjectsExternalDataTables },
  },
  {
    path: OLD_ROUTES.entityInfosProperties,
    redirect: { name: ROUTES_NAMES.adminCustomObjectsEntityInfosProperties },
  },
  {
    path: OLD_ROUTES.parameters,
    redirect: { name: ROUTES_NAMES.adminSettingsUserInterface },
  },

  {
    path: '*',
    redirect: { name: ROUTES_NAMES.home },
  },
];

const router = new Router({
  mode: ROUTER_MODE,
  routes,
});

/**
 * If requiresLogin is undefined then we can visit this page with auth and without auth
 */
router.beforeEach(async (to, from, next) => {
  const isRequiresAuth = to.matched.some(v => v.meta.requiresLogin);
  const isDontRequiresAuth = to.matched.some(v => v.meta.requiresLogin === false);
  const isLoggedIn = store.getters['auth/isLoggedIn'];
  const { query: { [ROUTER_ACCESS_TOKEN_KEY]: accessToken, ...restQuery } = {} } = to;

  if (accessToken) {
    await store.dispatch('auth/applyAccessToken', accessToken);

    return router.replace({
      ...to,
      query: restQuery,
    });
  }

  if (!isLoggedIn && isRequiresAuth) {
    return next({
      name: ROUTES_NAMES.login,
      query: {
        redirect: to.fullPath,
        errorMessage: to.query.errorMessage,
      },
    });
  }

  if (isLoggedIn && isDontRequiresAuth) {
    return next({
      name: ROUTES_NAMES.home,
    });
  }

  return next();
});

router.beforeResolve(async (to, from, next) => {
  try {
    await checkAppInfoAccessForRoute(to);
    await checkUserAccessForRoute(to);

    next();
  } catch (err) {
    console.error(err);

    next({
      name: ROUTES_NAMES.home,
    });
  }
});

router.afterEach((to, from) => {
  if (to.path !== from.path) {
    store.dispatch('entities/sweep');
  }
});

router.onReady((route) => {
  const { errorMessage } = route.query;

  if (errorMessage) {
    store.dispatch('popups/error', { text: errorMessage, autoClose: false });
  }
});

/**
 * Promisified router replace method
 *
 * @param {Object} route
 * @returns {Promise<unknown>}
 */
router.replaceAsync = route => new Promise((resolve, reject) => {
  router.replace(route, resolve, reject);
});

export default router;
