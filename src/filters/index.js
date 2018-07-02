import get from './getProp';
import formatContextSearch from './contextSearchFilter';

export default {
  install(Vue) {
    Vue.filter('get', get);
    Vue.filter('formatContextSearch', formatContextSearch);
  },
};
