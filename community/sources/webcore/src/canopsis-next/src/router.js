import Vue from 'vue';
import Router from 'vue-router';

import { ROUTER_MODE, ROUTER_ACCESS_TOKEN_KEY } from '@/config';
import { CRUD_ACTIONS, ROUTES_NAMES, ROUTES, USERS_PERMISSIONS } from '@/constants';

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
const AdminParameters = () => import(/* webpackChunkName: "Parameters" */ '@/views/admin/parameters.vue');
const AdminBroadcastMessages = () => import(/* webpackChunkName: "BroadcastMessage" */ '@/views/admin/broadcast-messages.vue');
const AdminPlaylists = () => import(/* webpackChunkName: "Playlist" */ '@/views/admin/playlists.vue');
const AdminPlanning = () => import(/* webpackChunkName: "Planning" */ '@/views/admin/planning.vue');
const AdminRemediation = () => import(/* webpackChunkName: "Remediation" */ '@/views/admin/remediation.vue');
const AdminHealthcheck = () => import(/* webpackChunkName: "Healthcheck" */ '@/views/admin/healthcheck.vue');
const AdminKPI = () => import(/* webpackChunkName: "KPI" */ '@/views/admin/kpi.vue');
const AdminMaps = () => import(/* webpackChunkName: "Maps" */ '@/views/admin/maps.vue');
const AdminTags = () => import(/* webpackChunkName: "Tags" */ '@/views/admin/tags.vue');
const AdminStorageSettings = () => import(/* webpackChunkName: "Tags" */ '@/views/admin/storage-settings.vue');
const AdminStateSettings = () => import(/* webpackChunkName: "Tags" */ '@/views/admin/state-settings.vue');
const AdminEventsRecords = () => import(/* webpackChunkName: "EventsRecords" */ '@/views/admin/events-records.vue');
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
const ProfilePatterns = () => import(/* webpackChunkName: "Pattern" */ '@/views/profile/patterns.vue');
const ProfileThemes = () => import(/* webpackChunkName: "Theme" */ '@/views/profile/themes.vue');
const Playlist = () => import(/* webpackChunkName: "Playlist" */ '@/views/playlist.vue');
const NotificationInstructionStats = () => import(/* webpackChunkName: "InstructionStats" */ '@/views/notification/instruction-stats.vue');

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
        id: USERS_PERMISSIONS.technical.view,
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
        id: USERS_PERMISSIONS.technical.permission,
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
        id: USERS_PERMISSIONS.technical.user,
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
        id: USERS_PERMISSIONS.technical.role,
      },
    },
  },
  {
    path: ROUTES.adminParameters,
    name: ROUTES_NAMES.adminParameters,
    component: AdminParameters,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.parameters,
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
        id: USERS_PERMISSIONS.technical.broadcastMessage,
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
        id: USERS_PERMISSIONS.technical.playlist,
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
        id: USERS_PERMISSIONS.technical.planning,
      },
    },
  },
  {
    path: ROUTES.adminRemediation,
    name: ROUTES_NAMES.adminRemediation,
    component: AdminRemediation,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.remediation,
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
        id: USERS_PERMISSIONS.technical.healthcheck,
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
        id: USERS_PERMISSIONS.technical.kpi,
      },
    },
  },
  {
    path: ROUTES.adminMaps,
    name: ROUTES_NAMES.adminMaps,
    component: AdminMaps,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.map,
      },
    },
  },
  {
    path: ROUTES.adminTags,
    name: ROUTES_NAMES.adminTags,
    component: AdminTags,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.tag,
      },
    },
  },
  {
    path: ROUTES.adminStorageSettings,
    name: ROUTES_NAMES.adminStorageSettings,
    component: AdminStorageSettings,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.storageSettings,
      },
    },
  },
  {
    path: ROUTES.adminStateSettings,
    name: ROUTES_NAMES.adminStateSettings,
    component: AdminStateSettings,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.stateSetting,
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
        id: USERS_PERMISSIONS.technical.eventsRecord,
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
        id: USERS_PERMISSIONS.technical.exploitation.pbehavior,
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
        id: USERS_PERMISSIONS.technical.exploitation.eventFilter,
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
        id: USERS_PERMISSIONS.technical.exploitation.snmpRule,
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
        id: USERS_PERMISSIONS.technical.exploitation.dynamicInfo,
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
    path: ROUTES.exploitationMetaAlarmRules,
    name: ROUTES_NAMES.exploitationMetaAlarmRules,
    component: ExploitationMetaAlarmRules,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.exploitation.metaAlarmRule,
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
        id: USERS_PERMISSIONS.technical.exploitation.scenario,
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
        id: USERS_PERMISSIONS.technical.exploitation.idleRules,
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
        id: USERS_PERMISSIONS.technical.exploitation.flappingRules,
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
        id: USERS_PERMISSIONS.technical.exploitation.resolveRules,
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
        id: USERS_PERMISSIONS.technical.exploitation.declareTicketRule,
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
        id: USERS_PERMISSIONS.technical.exploitation.linkRule,
      },
    },
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
        id: USERS_PERMISSIONS.technical.profile.theme,
      },
    },
  },
  {
    path: ROUTES.notificationInstructionStats,
    name: ROUTES_NAMES.notificationInstructionStats,
    component: NotificationInstructionStats,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.notification.instructionStats,
      },
    },
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
  {
    path: '*',
    redirect: {
      name: ROUTES_NAMES.home,
    },
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
