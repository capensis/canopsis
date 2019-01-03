import get from './get';
import date from './date';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', date);
  },
};
