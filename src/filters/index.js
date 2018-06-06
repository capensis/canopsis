import get from '@/filters/getProp';

export default {
  install(Vue) {
    Vue.filter('get', get);
  },
};
