import Vue from 'vue';
import Router from 'vue-router';
import Cookies from 'js-cookie';

import { ROUTER_MODE, COOKIE_SESSION_KEY } from '@/config';
import Login from '@/views/login.vue';
import Home from '@/views/home.vue';
import View from '@/views/view.vue';
import AdminRights from '@/views/admin/rights.vue';

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
    meta: requiresLoginMeta,
    props: route => ({ id: route.params.id }),
  },
  {
    path: '/admin/rights',
    name: 'admin-rights',
    component: AdminRights,
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

export default router;
