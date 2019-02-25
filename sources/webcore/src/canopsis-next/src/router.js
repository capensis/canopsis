import Vue from 'vue';
import Router from 'vue-router';
import Cookies from 'js-cookie';
import { isEmpty, isFunction } from 'lodash';

import { ROUTER_MODE, COOKIE_SESSION_KEY } from '@/config';
import { USERS_RIGHTS, USERS_RIGHTS_MASKS } from '@/constants';
import store from '@/store';
import i18n from '@/i18n';
import { checkUserAccess } from '@/helpers/right';

import Login from '@/views/login.vue';
import Home from '@/views/home.vue';
import View from '@/views/view.vue';
import AdminRights from '@/views/admin/rights.vue';
import AdminUsers from '@/views/admin/users.vue';
import AdminRoles from '@/views/admin/roles.vue';
import AdminParameters from '@/views/admin/parameters.vue';
import ExploitationPbehaviors from '@/views/exploitation/pbehaviors.vue';
import ExploitationEventFilter from '@/views/exploitation/event-filter.vue';

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
    meta: requiresLoginMeta,
  },
  {
    path: '/exploitation/pbehaviors',
    name: 'exploitation-pbehaviors',
    component: ExploitationPbehaviors,
    meta: requiresLoginMeta,
  },
  {
    path: '/exploitation/event-filter',
    name: 'exploitation-event-filter',
    component: ExploitationEventFilter,
    meta: requiresLoginMeta,
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

/**
 * if route has requiresRight we will wait currentUser object and check right
 */
router.beforeResolve((to, from, next) => {
  if (to.meta.requiresLogin && to.meta.requiresRight) {
    const { requiresRight } = to.meta;
    const rightId = isFunction(requiresRight.id) ? requiresRight.id(to) : requiresRight.id;
    const rightMask = requiresRight.mask ? requiresRight.mask : USERS_RIGHTS_MASKS.read;

    const checkProcess = (user) => {
      if (checkUserAccess(user, rightId, rightMask)) {
        next();
      } else {
        store.dispatch('popup/add', { text: i18n.t('common.forbidden') });

        next({
          name: 'home',
        });
      }
    };

    if (isEmpty(store.getters['auth/currentUser'])) {
      const unwatch = store.watch(
        state => state.auth.currentUser,
        (currentUser) => {
          if (!isEmpty(currentUser)) {
            unwatch();
            checkProcess(currentUser);
          }
        },
      );
    } else {
      checkProcess(store.getters['auth/currentUser']);
    }
  } else {
    next();
  }
});

router.afterEach(() => {
  store.dispatch('entities/sweep');
});

export default router;
