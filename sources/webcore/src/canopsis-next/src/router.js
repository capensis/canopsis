import Vue from 'vue';
import Router from 'vue-router';
import Cookies from 'js-cookie';

import { ROUTER_MODE, COOKIE_SESSION_KEY } from '@/config';
import { USERS_RIGHTS } from '@/constants';
import store from '@/store';
import { checkAppInfoAccessForRoute, checkUserAccessForRoute, getKeepalivePathByRoute } from '@/helpers/router';

import Login from '@/views/login.vue';
import Home from '@/views/home.vue';
import View from '@/views/view.vue';
import AdminRights from '@/views/admin/rights.vue';
import AdminUsers from '@/views/admin/users.vue';
import AdminRoles from '@/views/admin/roles.vue';
import AdminParameters from '@/views/admin/parameters.vue';
import ExploitationPbehaviors from '@/views/exploitation/pbehaviors.vue';
import ExploitationEventFilter from '@/views/exploitation/event-filter.vue';
import ExploitationWebhooks from '@/views/exploitation/webhooks.vue';
import ExploitationSnmpRules from '@/views/exploitation/snmp-rules.vue';
import ExploitationActions from '@/views/exploitation/actions.vue';
import ExploitationHeartbeats from '@/views/exploitation/heartbeats.vue';
import ExploitationDynamicInfos from '@/views/exploitation/dynamic-infos.vue';

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
    path: '/admin/rights',
    name: 'admin-rights',
    component: AdminRights,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_RIGHTS.technical.action,
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
        id: USERS_RIGHTS.technical.user,
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
        id: USERS_RIGHTS.technical.role,
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
        id: USERS_RIGHTS.technical.parameters,
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
        id: USERS_RIGHTS.technical.exploitation.pbehavior,
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
        id: USERS_RIGHTS.technical.exploitation.eventFilter,
      },
    },
  },
  {
    path: '/exploitation/webhooks',
    name: 'exploitation-webhooks',
    component: ExploitationWebhooks,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_RIGHTS.technical.exploitation.webhook,
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
        id: USERS_RIGHTS.technical.exploitation.snmpRule,
      },
    },
  },
  {
    path: '/exploitation/actions',
    name: 'exploitation-actions',
    component: ExploitationActions,
    meta: {
      requiresLogin: true,
      requiresRight: {
        id: USERS_RIGHTS.technical.exploitation.action,
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
        id: USERS_RIGHTS.technical.exploitation.heartbeat,
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
        id: USERS_RIGHTS.technical.exploitation.dynamicInfo,
      },
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
    store.dispatch('keepalive/sessionHide', { path: getKeepalivePathByRoute(to) });
  }
});

export default router;
