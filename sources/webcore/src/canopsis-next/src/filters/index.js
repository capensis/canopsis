import get from './get';
import date from './date';
import json from './json';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', date);
    Vue.filter('json', json);
  },
};
