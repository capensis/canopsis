import get from './get';
import formatContextSearch from './contextSearchFilter';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('formatContextSearch', formatContextSearch);
  },
};
