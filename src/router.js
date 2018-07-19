import Vue from 'vue';
import Router from 'vue-router';
import { ROUTER_MODE, ENVIRONNEMENT } from '@/config';
import Login from '@/views/login.vue';
import Alarm from '@/views/alarm.vue';
import Home from '@/views/home.vue';
import About from '@/views/about.vue';
import Context from '@/views/context.vue';
import Filter from '@/views/filter.vue';
import Rrule from '@/components/forms/rrule.vue';

Vue.use(Router);

let routes = [
  {
    path: '/',
    name: 'home',
    component: Home,
  },
  {
    path: '/context',
    name: 'context',
    component: Context,
  },
];
if (ENVIRONNEMENT === 'development') {
  routes = routes.concat([
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
    {
      path: '/about',
      name: 'about',
      component: About,
    },
    {
      path: '/alarms',
      name: 'alarms',
      component: Alarm,
    },
    {
      path: '/filter',
      name: 'filter',
      component: Filter,
    },
    {
      path: '/rrule',
      name: 'rrule',
      component: Rrule,
    },
  ]);
}

export default new Router({
  mode: ROUTER_MODE,
  routes,
});
