import Vue from 'vue';
import Router from 'vue-router';

import Alarm from '@/views/alarm.vue';
import Home from '@/views/home.vue';
import About from '@/views/about.vue';

// EXAMPLES
import MFilterEditor from '@/components/mfilter-editor.vue';
import FilterSelector from '@/components/other/filter/selector.vue';
import RRuleForm from '@/components/other/rrule/rrule-form.vue';

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
      path: '/alarms',
      name: 'alarms',
      component: Alarm,
    },
    {
      path: '/filter-editor',
      name: 'filter-editor',
      component: MFilterEditor,
    },
    {
      path: '/filter-selector',
      name: 'filter-selector',
      component: FilterSelector,
    },
    {
      path: '/rrule-form',
      name: 'rruleForm',
      component: RRuleForm,
    },
  ],
});
