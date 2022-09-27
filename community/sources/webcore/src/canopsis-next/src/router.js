import Vue from 'vue';
import Router from 'vue-router';

import { ROUTER_MODE, ROUTER_ACCESS_TOKEN_KEY } from '@/config';
import {
  CRUD_ACTIONS,
  ROUTES_NAMES,
  ROUTES,
  USERS_PERMISSIONS,
} from '@/constants';
import store from '@/store';
import {
  checkAppInfoAccessForRoute,
  checkUserAccessForRoute,
} from '@/helpers/router';

import Login from '@/views/login.vue';
import Error from '@/views/error.vue';

const Home = () => import(/* webpackChunkName: "Home" */ '@/views/home.vue');
const View = () => import(/* webpackChunkName: "View" */ '@/views/view.vue');
const Alarm = () => import(/* webpackChunkName: "Alarm" */ '@/views/alarm.vue');
const AdminPermissions = () => import(/* webpackChunkName: "Permission" */ '@/views/admin/permissions.vue');
const AdminUsers = () => import(/* webpackChunkName: "User" */ '@/views/admin/users.vue');
const AdminRoles = () => import(/* webpackChunkName: "Role" */ '@/views/admin/roles.vue');
const AdminParameters = () => import(/* webpackChunkName: "Parameters" */ '@/views/admin/parameters.vue');
const AdminBroadcastMessages = () => import(/* webpackChunkName: "BroadcastMessage" */ '@/views/admin/broadcast-messages.vue');
const AdminPlaylists = () => import(/* webpackChunkName: "Playlist" */ '@/views/admin/playlists.vue');
const AdminPlanning = () => import(/* webpackChunkName: "Planning" */ '@/views/admin/planning.vue');
const AdminRemediation = () => import(/* webpackChunkName: "Remediation" */ '@/views/admin/remediation.vue');
const AdminEngines = () => import(/* webpackChunkName: "Engines" */ '@/views/admin/engines.vue');
const AdminHealthcheck = () => import(/* webpackChunkName: "Healthcheck" */ '@/views/admin/healthcheck.vue');
const AdminShareTokens = () => import(/* webpackChunkName: "ShareTokens" */ '@/views/admin/share-tokens.vue');
const AdminKPI = () => import(/* webpackChunkName: "KPI" */ '@/views/admin/kpi.vue');
const ExploitationPbehaviors = () => import(/* webpackChunkName: "Pbehavior" */ '@/views/exploitation/pbehaviors.vue');
const ExploitationEventFilters = () => import(/* webpackChunkName: "EventFilters" */ '@/views/exploitation/event-filters.vue');
const ExploitationSnmpRules = () => import(/* webpackChunkName: "SnmpRule" */ '@/views/exploitation/snmp-rules.vue');
const ExploitationDynamicInfos = () => import(/* webpackChunkName: "DynamicInfo" */ '@/views/exploitation/dynamic-infos.vue');
const ExploitationMetaAlarmRules = () => import(/* webpackChunkName: "MetaAlarmRule" */ '@/views/exploitation/meta-alarm-rules.vue');
const ExploitationScenarios = () => import(/* webpackChunkName: "Scenario" */ '@/views/exploitation/scenarios.vue');
const ExploitationIdleRules = () => import(/* webpackChunkName: "IdleRule" */ '@/views/exploitation/idle-rules.vue');
const ExploitationFlappingRules = () => import(/* webpackChunkName: "AlarmStatusRule" */ '@/views/exploitation/flapping-rules.vue');
const ExploitationResolveRules = () => import(/* webpackChunkName: "AlarmStatusRule" */ '@/views/exploitation/resolve-rules.vue');
const ProfilePatterns = () => import(/* webpackChunkName: "Pattern" */ '@/views/profile/patterns.vue');
const Playlist = () => import(/* webpackChunkName: "Playlist" */ '@/views/playlist.vue');
const NotificationInstructionStats = () => import(/* webpackChunkName: "InstructionStats" */ '@/views/notification/instruction-stats.vue');

Vue.use(Router);

const routes = [
  {
    path: ROUTES.login,
    name: ROUTES_NAMES.login,
    component: Login,
    meta: {
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
      requiresPermission: {
        id: route => route.params.id,
      },
    },
    props: route => ({ id: route.params.id }),
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
    path: ROUTES.adminEngines,
    name: ROUTES_NAMES.adminEngines,
    component: AdminEngines,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USERS_PERMISSIONS.technical.engine,
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
    path: ROUTES.adminShareTokens,
    name: ROUTES_NAMES.adminShareTokens,
    component: AdminShareTokens,
    meta: {
      requiresLogin: true,
      requiresPermission: {
        id: USERS_PERMISSIONS.technical.shareToken,
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
    path: ROUTES.profilePatterns,
    name: ROUTES_NAMES.profilePatterns,
    component: ProfilePatterns,
    meta: {
      requiresLogin: true,
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
