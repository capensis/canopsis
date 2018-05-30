import Vue from 'vue';
import Router from 'vue-router';

import Alarm from '@/views/alarm.vue';
import Home from '@/views/home.vue';
import About from '@/views/about.vue';
import Context from '@/views/context.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/about',
      name: 'about',
      component: About,
    },
    {
      path: '/alarm',
      name: 'alarm',
      component: Alarm,
    },
    {
      path: '/context',
      name: 'context',
      component: Context,
    },
  ],
});
