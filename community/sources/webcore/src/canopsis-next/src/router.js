import Vue from 'vue';
import Router from 'vue-router';

import { ROUTER_MODE } from '@/config';
import { CRUD_ACTIONS, ROUTES_NAMES, ROUTES, USERS_PERMISSIONS } from '@/constants';
import store from '@/store';
import {
  checkAppInfoAccessForRoute,
  checkUserAccessForRoute,
  getViewStatsPathByRoute,
} from '@/helpers/router';

import Login from '@/views/login.vue';
import Home from '@/views/home.vue';
import View from '@/views/view.vue';
import Alarm from '@/views/alarm.vue';
import AdminPermissions from '@/views/admin/permissions.vue';
import AdminUsers from '@/views/admin/users.vue';
import AdminRoles from '@/views/admin/roles.vue';
import AdminParameters from '@/views/admin/parameters.vue';
import AdminBroadcastMessages from '@/views/admin/broadcast-messages.vue';
import AdminPlaylists from '@/views/admin/playlists.vue';
import AdminPlanning from '@/views/admin/planning.vue';
import AdminRemediation from '@/views/admin/remediation.vue';
import AdminHealthcheck from '@/views/admin/healthcheck.vue';
import ExploitationPbehaviors from '@/views/exploitation/pbehaviors.vue';
import ExploitationEventFilter from '@/views/exploitation/event-filter.vue';
import ExploitationSnmpRules from '@/views/exploitation/snmp-rules.vue';
import ExploitationDynamicInfos from '@/views/exploitation/dynamic-infos.vue';
import ExploitationMetaAlarmRules from '@/views/exploitation/meta-alarm-rules.vue';
import ExploitationScenarios from '@/views/exploitation/scenarios.vue';
import ExploitationIdleRules from '@/views/exploitation/idle-rules.vue';
import ExploitationFlappingRules from '@/views/exploitation/flapping-rules.vue';
import ExploitationResolveRules from '@/views/exploitation/resolve-rules.vue';
import Playlist from '@/views/playlist.vue';
import NotificationInstructionStats from '@/views/notification/instruction-stats.vue';
import Error from '@/views/error.vue';

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
        id: USERS_PERMISSIONS.technical.action,
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
    path: ROUTES.exploitationEventFilter,
    name: ROUTES_NAMES.exploitationEventFilter,
    component: ExploitationEventFilter,
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

  const { query: { access_token: accessToken, ...restQuery } = {} } = to;

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
    next({
      name: ROUTES_NAMES.home,
    });
  }
});

router.afterEach((to, from) => {
  const isLoggedIn = store.getters['auth/isLoggedIn'];

  if (to.path !== from.path) {
    store.dispatch('entities/sweep');
  }

  if (isLoggedIn) {
    store.dispatch('viewStats/update', {
      data: {
        visible: !(document.visibilityState === 'hidden'),
        path: getViewStatsPathByRoute(to),
      },
    });
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
