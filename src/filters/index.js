import get from './getProp';

export default {
  install(Vue) {
    Vue.filter('get', get);
  },
};
