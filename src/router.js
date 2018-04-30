import Vue from 'vue';
import Router from 'vue-router';

import Home from './views/home.vue';
import About from './views/about.vue';
import Context from './views/alarm.vue';

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
      path: '/context',
      name: 'context',
      component: Context,
    },
  ],
});
