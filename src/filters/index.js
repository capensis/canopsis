import get from './getProp';
import date from './date';
import formatContextSearch from './contextSearchFilter';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('date', date);
    Vue.filter('formatContextSearch', formatContextSearch);
  },
};
