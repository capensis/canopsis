import Vue from 'vue';
import Router from 'vue-router';
import Cookies from 'js-cookie';

import { ROUTER_MODE, COOKIE_SESSION_KEY } from '@/config';
import Login from '@/views/login.vue';
import Alarm from '@/views/alarm.vue';
import Home from '@/views/home.vue';
import About from '@/views/about.vue';
import Context from '@/views/context.vue';
import Filter from '@/views/filter.vue';
import View from '@/views/view.vue';
import Rrule from '@/components/forms/rrule.vue';
import Weather from '@/views/weather.vue';

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
    path: '/context',
    name: 'context',
    component: Context,
    meta: requiresLoginMeta,
  },
  {
    path: '/view/:id',
    name: 'view',
    component: View,
    props: route => ({ id: route.params.id }),
  },
];
if (process.env.NODE_ENV === 'development') {
  routes.push(
    {
      path: '/about',
      name: 'about',
      component: About,
      meta: requiresLoginMeta,
    },
    {
      path: '/alarms',
      name: 'alarms',
      component: Alarm,
      meta: requiresLoginMeta,
    },
    {
      path: '/filter',
      name: 'filter',
      component: Filter,
      meta: requiresLoginMeta,
    },
    {
      path: '/rrule',
      name: 'rrule',
      component: Rrule,
      meta: requiresLoginMeta,
    },
    {
      path: '/weather',
      name: 'weather',
      component: Weather,
    },
  );
}

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
