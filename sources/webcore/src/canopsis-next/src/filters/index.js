import get from './get';
import date from './date';
import duration from './duration';
import json from './json';
import percentage from './percentage';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', date);
    Vue.filter('duration', duration);
    Vue.filter('json', json);
    Vue.filter('percentage', percentage);
  },
};
