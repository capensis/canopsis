export default {
  install(Vue) {
    Vue.setSeveral = function setSeveral(source, key, payload) {
      if (source[key]) {
        Object.entries(payload).forEach(([itemKey, itemValue]) => Vue.set(source[key], itemKey, itemValue));
      } else {
        Vue.set(source, key, payload);
      }
    };
  },
};
