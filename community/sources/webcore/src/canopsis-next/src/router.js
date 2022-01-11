import Vue from 'vue';
import Router from 'vue-router';
import Cookies from 'js-cookie';

import { ROUTER_MODE, COOKIE_SESSION_KEY } from '@/config';
import { CRUD_ACTIONS, USERS_PERMISSIONS } from '@/constants';
import store from '@/store';
import { checkAppInfoAccessForRoute, checkUserAccessForRoute, getKeepalivePathByRoute } from '@/helpers/router';

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
import AdminEngines from '@/views/admin/engines.vue';
import AdminRemediation from '@/views/admin/remediation.vue';
import ExploitationPbehaviors from '@/views/exploitation/pbehaviors.vue';
import ExploitationEventFilter from '@/views/exploitation/event-filter.vue';
import ExploitationSnmpRules from '@/views/exploitation/snmp-rules.vue';
import ExploitationHeartbeats from '@/views/exploitation/heartbeats.vue';
import ExploitationDynamicInfos from '@/views/exploitation/dynamic-infos.vue';
import Playlist from '@/views/playlist.vue';
import ExploitationMetaAlarmRules from '@/views/exploitation/meta-alarm-rules.vue';
import ExploitationScenarios from '@/views/exploitation/scenarios.vue';

Vue.use(Router);

const requiresLoginMeta = {
  requiresLogin: true,
};

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: {
      requiresLogin: false,
    },
  },
  {
    path: '/',
    name: 'home',
    component: Home,
    meta: requiresLoginMeta,
  },
  {
    path: '/view/:id',
    name: 'view',
    component: View,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: route => route.params.id,
      },
    },
    props: route => ({ id: route.params.id }),
  },
  {
    path: '/alarms/:id',
    name: 'alarms',
    component: Alarm,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.view,
      },
    },
    props: route => ({ id: route.params.id }),
  },
  {
    path: '/admin/rights',
    name: 'admin-rights',
    component: AdminPermissions,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.permission,
      },
    },
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: AdminUsers,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.user,
      },
    },
  },
  {
    path: '/admin/roles',
    name: 'admin-roles',
    component: AdminRoles,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.role,
      },
    },
  },
  {
    path: '/admin/parameters',
    name: 'admin-parameters',
    component: AdminParameters,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.parameters,
      },
    },
  },
  {
    path: '/admin/broadcast-messages',
    name: 'admin-broadcast-messages',
    component: AdminBroadcastMessages,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.broadcastMessage,
      },
    },
  },
  {
    path: '/admin/playlists',
    name: 'admin-playlists',
    component: AdminPlaylists,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.playlist,
      },
    },
  },
  {
    path: '/admin/planning',
    name: 'admin-planning-administration',
    component: AdminPlanning,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.planning,
      },
    },
  },
  {
    path: '/admin/engines',
    name: 'admin-engines',
    component: AdminEngines,
    meta: {
      requiresLogin: true,
      requiresRight: {
        action: CRUD_ACTIONS.can,
        id: USERS_PERMISSIONS.technical.engine,
      },
    },
  },
  {
    path: '/admin/remediation',
    name: 'admin-remediation-administration',
    component: AdminRemediation,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.remediation,
      },
    },
  },
  {
    path: '/exploitation/pbehaviors',
    name: 'exploitation-pbehaviors',
    component: ExploitationPbehaviors,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.pbehavior,
      },
    },
  },
  {
    path: '/exploitation/event-filter',
    name: 'exploitation-event-filter',
    component: ExploitationEventFilter,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.eventFilter,
      },
    },
  },
  {
    path: '/exploitation/snmp-rules',
    name: 'exploitation-snmp-rules',
    component: ExploitationSnmpRules,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.snmpRule,
      },
    },
  },
  {
    path: '/exploitation/heartbeats',
    name: 'exploitation-heartbeats',
    component: ExploitationHeartbeats,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.heartbeat,
      },
    },
  },
  {
    path: '/exploitation/dynamic-infos',
    name: 'exploitation-dynamic-infos',
    component: ExploitationDynamicInfos,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.dynamicInfo,
      },
    },
  },
  {
    path: '/playlist/:id',
    name: 'playlist',
    component: Playlist,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: route => route.params.id,
        action: CRUD_ACTIONS.read,
      },
    },
    props: route => ({ id: route.params.id, autoplay: String(route.query.autoplay) === 'true' }),
  },
  {
    path: '/exploitation/meta-alarm-rule',
    name: 'exploitation-meta-alarm-rules',
    component: ExploitationMetaAlarmRules,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.metaAlarmRule,
      },
    },
  },
  {
    path: '/exploitation/scenarios',
    name: 'exploitation-scenarios',
    component: ExploitationScenarios,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_PERMISSIONS.technical.exploitation.scenario,
      },
    },
  },
  {
    path: '*',
    redirect: {
      name: 'home',
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
router.beforeEach((to, from, next) => {
  const isRequiresAuth = to.matched.some(v => v.meta.requiresLogin);
  const isDontRequiresAuth = to.matched.some(v => v.meta.requiresLogin === false);
  const isLoggedIn = !!Cookies.get(COOKIE_SESSION_KEY);

  if (!isLoggedIn && isRequiresAuth) {
    return next({
      name: 'login',
      query: {
        redirect: to.fullPath,
        errorMessage: to.query.errorMessage,
      },
    });
  } else if (isLoggedIn && isDontRequiresAuth) {
    return next({
      name: 'home',
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
      name: 'home',
    });
  }
});

router.afterEach((to, from) => {
  const isLoggedIn = !!Cookies.get(COOKIE_SESSION_KEY);

  if (to.path !== from.path) {
    store.dispatch('entities/sweep');
  }

  if (isLoggedIn) {
    store.dispatch('keepalive/sessionTracePath', { path: getKeepalivePathByRoute(to) });
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
